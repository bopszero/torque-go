package sysforwardingmod

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comutils"

	"github.com/jszwec/csvutil"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func NewForwardingHandler(currency meta.Currency, network meta.BlockchainNetwork, date string) (ForwardingHandler, error) {
	if constants.BlockchainChannelUtxoCurrencySet.Contains(currency) {
		return NewUtxoLikeForwardingHandler(currency, network, date)
	} else {
		return NewBalanceLikeForwardingHandler(currency, network, date)
	}
}

func updateTradingDeposits(
	ctx comcontext.Context,
	ids []uint64, status meta.SystemForwardingOrderTxnStatus,
) error {
	return database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.Model(&models.Deposit{}).
			Where(dbquery.In(models.DepositColID, ids)).
			Updates(&models.Deposit{
				ForwardStatus: status,
				UpdateTime:    time.Now().Unix(),
			}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
}

func ImportAddresses(ctx comcontext.Context, csvReader *csv.Reader) (count uint32, err error) {
	headers, err := csvReader.Read()
	if err != nil {
		return 0, utils.WrapError(err)
	}
	if len(headers) < 3 {
		return 0, utils.WrapError(constants.ErrorInvalidParams)
	}

	headerNetwork := headers[ImportIndexNetwork]
	headerAddress := headers[ImportIndexAddress]
	headerKey := headers[ImportIndexKey]
	isValidHeaders := headerNetwork == ImportColNetwork &&
		headerAddress == ImportColAddress &&
		headerKey == ImportColKey
	if !isValidHeaders {
		return 0, utils.WrapError(constants.ErrorInvalidParams)
	}

	genModel := func(values []string) (model models.SystemForwardingAddress, err error) {
		var (
			addressNetwork = meta.NewBlockchainNetwork(values[ImportIndexNetwork])
			addressValue   = values[ImportIndexAddress]
			addressKey     = values[ImportIndexKey]
		)

		coin, err := blockchainmod.GetNetworkMainCoin(addressNetwork)
		if err != nil {
			return
		}
		normalizedAddress, err := coin.NormalizeAddress(addressValue)
		if err != nil {
			err = utils.IssueErrorf(
				"system forwarding import invalid %v address `%v`",
				addressNetwork, addressValue,
			)
			return
		}

		keyHolder, err := coin.LoadKey(addressKey, normalizedAddress)
		if err != nil {
			return
		}
		if keyHolder.GetAddress() != normalizedAddress {
			err = utils.IssueErrorf(
				"system forwarding import address and key don't match | network=%v,address=%v",
				addressNetwork, addressValue,
			)
			return
		}

		model.Network = addressNetwork
		model.Address = normalizedAddress
		model.Key = models.NewUserAddressKeyEncryptedField(addressKey)
		model.CreateTime = time.Now().Unix()

		return model, nil
	}
	err = database.Atomic(ctx, database.AliasInternalMaster, func(dbTxn *gorm.DB) error {
		addressModels := make([]models.SystemForwardingAddress, 0)
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return utils.WrapError(err)
			}

			addressModel, err := genModel(record)
			if err != nil {
				return err
			}

			addressModels = append(addressModels, addressModel)
			count++

			if len(addressModels) >= 200 {
				if err := dbTxn.CreateInBatches(addressModels, 1000).Error; err != nil {
					return utils.WrapError(err)
				}
				addressModels = make([]models.SystemForwardingAddress, 0)
			}
		}
		if err := dbTxn.CreateInBatches(addressModels, 1000).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func GenerateReport(date string, filePath string) (err error) {
	fileIO, err := os.Create(filePath)
	if err != nil {
		return utils.WrapError(err)
	}
	defer comutils.PanicOnError(fileIO.Close())

	var (
		dbMain     = database.GetDbSlave()
		dbInternal = database.GetDbF(database.AliasInternalMaster)
	)
	var orders []models.SystemForwardingOrder
	err = dbInternal.
		Find(
			&orders,
			&models.SystemForwardingOrder{
				Date: dbfields.NewDateFieldFromStringF(date),
			},
		).
		Error
	if err != nil {
		return utils.WrapError(err)
	}
	orderIDs := make([]uint64, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	var txns []models.SystemForwardingOrderTxn
	err = dbInternal.
		Where(dbquery.In(models.SystemForwardingOrderTxnColOrderID, orderIDs)).
		Find(&txns).
		Error
	if err != nil {
		return utils.WrapError(err)
	}

	var (
		hashToTxnsMap = make(map[string][]models.SystemForwardingOrderTxn)
		depositIDs    = make([]uint64, len(txns))
	)
	for i, txn := range txns {
		depositIDs[i] = txn.DepositID
		hashToTxnsMap[txn.Hash.String] = append(hashToTxnsMap[txn.Hash.String], txn)
	}

	var deposits []models.Deposit
	err = dbMain.
		Where(dbquery.In(models.DepositColID, depositIDs)).
		Find(&deposits).
		Error
	if err != nil {
		return utils.WrapError(err)
	}
	depositMap := make(map[uint64]*models.Deposit)
	for i := range deposits {
		d := deposits[i]
		depositMap[d.ID] = &d
	}

	reportItems := make([]ReportItem, 0, len(txns))
	for _, txns := range hashToTxnsMap {
		firstTxn := txns[0]
		reportItems = append(
			reportItems,
			dumpTxnToReportItem(firstTxn, depositMap[firstTxn.DepositID]),
		)
		for _, txn := range txns[1:] {
			txn.Hash.String = "" // For readable

			reportItems = append(
				reportItems,
				dumpTxnToReportItem(txn, depositMap[txn.DepositID]),
			)
		}
	}

	csvBytes, err := csvutil.Marshal(reportItems)
	if err != nil {
		return utils.WrapError(err)
	}

	_, err = fileIO.Write(csvBytes)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}
