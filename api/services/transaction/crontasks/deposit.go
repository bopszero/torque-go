package crontasks

import (
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/depositmod/depositcrawler"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

const (
	DepositCrawBlockAcceptAge = time.Minute
)

func DepositCrawBlocks(currency meta.Currency) {
	defer apiutils.CronRunWithRecovery(
		"DepositCrawBlocks",
		meta.O{"currency": currency},
	)

	var (
		ctx  = comcontext.NewContext()
		coin = blockchainmod.GetCoinNativeF(currency)
		err  error
	)
	client, err := coin.NewClientDefault()
	comutils.PanicOnError(err)

	crawler, err := depositcrawler.NewCrawlerDefaultOptions(ctx, coin)
	comutils.PanicOnError(err)

	scanHeights, err := crawler.GetScanBlockHeights()
	comutils.PanicOnError(err)

	var (
		wg          sync.WaitGroup
		blockErrMap = make(map[uint64]error)
	)
	poolConsumeFunc := func(blockObj interface{}) {
		defer wg.Done()
		var (
			block       = blockObj.(blockchainmod.Block)
			blockHeight = block.GetHeight()
		)

		err := database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
			if err := crawler.ConsumeBlock(block); err != nil {
				return err
			}
			comlogging.GetLogger().
				WithFields(logrus.Fields{
					"currency": coin.GetCurrency(),
					"height":   block.GetHeight(),
				}).
				Infof("deposit crawled `%s` block", coin.GetIndexCode())
			return nil
		})
		blockErrMap[blockHeight] = err
	}
	pool, err := ants.NewPoolWithFunc(1, poolConsumeFunc)
	comutils.PanicOnError(err)
	defer pool.Release()

	var scanErr error
	for _, height := range scanHeights {
		block, err := client.GetBlockByHeight(height)
		if err != nil {
			if !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
				scanErr = err
			}
			break
		}
		acceptTimeUnix := time.Now().Add(-DepositCrawBlockAcceptAge).Unix()
		if block.GetTimeUnix() > acceptTimeUnix {
			break
		}

		wg.Add(1)
		if scanErr = pool.Invoke(block); scanErr != nil {
			err = utils.WrapError(err)
			break
		}
	}
	wg.Wait()

	var (
		acceptedBlockHeight uint64 = 0
		consumeErr          error
	)
	for _, height := range scanHeights {
		err, ok := blockErrMap[height]
		if !ok || err != nil {
			consumeErr = err
			break
		}
		acceptedBlockHeight = height
	}
	if acceptedBlockHeight > 0 {
		err := database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
			return depositcrawler.UpdateCrawledBlock(ctx, coin, acceptedBlockHeight)
		})
		comutils.PanicOnError(err)
	}
	comutils.PanicOnError(consumeErr)
	comutils.PanicOnError(scanErr)
}

func DepositCollectAndApprove() {
	defer apiutils.CronRunWithRecovery("DepositCollectAndApprove", nil)

	var (
		logger      = comlogging.GetLogger()
		pendingTxns []models.DepositCryptoTxn
		err         error
	)
	err = database.GetDbSlave().
		Find(
			&pendingTxns,
			&models.DepositCryptoTxn{
				IsAccepted: models.NewBool(false),
			},
		).
		Error
	comutils.PanicOnError(err)

	for _, txn := range pendingTxns {
		if err := depositCollectAndApproveTxn(txn); err != nil {
			logger.
				WithType(constants.LogTypeDeposit).
				WithError(err).
				WithFields(logrus.Fields{
					"to_address": txn.ToAddress,
					"to_index":   txn.ToIndex,
					"hash":       txn.Hash,
				}).
				Errorf("collect and approve crypto txn failed | id=%v,hash=%v", txn.ID, txn.Hash)
		}
	}
}

func depositCollectAndApproveTxn(txn models.DepositCryptoTxn) (err error) {
	var (
		ctx     = comcontext.NewContext()
		dbSlave = database.GetDbSlave()

		coin             = blockchainmod.GetCoinF(txn.Currency, txn.Network)
		networkInfo      = coin.GetModelNetwork()
		minConfirmations = networkInfo.DepositMinConfirmations
	)
	if minConfirmations == 0 {
		comlogging.GetLogger().
			WithField("network", networkInfo.Network).
			Warnf("blockchain network `%v`'s min deposit confirmations hasn't been set", networkInfo.Network)
		minConfirmations = 20 // Fail safe
	}

	var existsDeposit models.Deposit
	err = dbSlave.
		First(
			&existsDeposit,
			&models.Deposit{
				Currency: txn.Currency,
				Network:  txn.Network,
				TxnHash:  txn.Hash,
				TxnIndex: txn.ToIndex,
			},
		).
		Error
	if database.IsDbError(err) {
		return utils.WrapError(err)
	}
	if existsDeposit.ID > 0 && txn.Confirmations < minConfirmations {
		return nil
	}

	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		if existsDeposit.ID == 0 {
			userDepositAddress, err := depositmod.GetDepositAddressByAddress(coin, txn.ToAddress)
			if err != nil {
				return err
			}
			existsDeposit, err = depositmod.SubmitDeposit(
				ctx, userDepositAddress.UID,
				txn.Currency, txn.Network,
				txn.Hash, txn.ToIndex, txn.ToAddress, txn.Amount,
			)
			if err != nil {
				return err
			}
		}

		if txn.Confirmations < minConfirmations {
			return nil
		}

		if _, err = depositmod.ApproveDeposit(ctx, existsDeposit.ID, ""); err != nil {
			return
		}

		txn.IsAccepted = models.NewBool(true)
		txn.UpdateTime = time.Now().Unix()
		if err = dbTxn.Save(&txn).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
}
