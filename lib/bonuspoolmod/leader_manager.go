package bonuspoolmod

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func LeaderCalcPoolAmount(fromDate, toDate string) (amount decimal.Decimal, err error) {
	fromTime, toTime, err := leaderValidateDateRange(fromDate, toDate)
	if err != nil {
		return
	}
	var (
		expectDateCount = int(toTime.Sub(fromTime).Hours()/24) + 1
		payoutStatsList []models.PayoutStats
	)
	err = database.GetDbSlave().
		Where(dbquery.Between(models.CommonColDate, fromDate, toDate)).
		Order(dbquery.OrderAsc(models.CommonColDate)).
		Find(&payoutStatsList).
		Error
	if err != nil {
		return
	}
	if len(payoutStatsList) != expectDateCount {
		err = meta.NewMessageError(
			"Leader Bonus Pool cannot collect enough Payout Stats %v/%v",
			len(payoutStatsList), expectDateCount)
		return
	}

	for _, stats := range payoutStatsList {
		amount = amount.
			Add(stats.RemainingAmountAffiliate).
			Add(stats.RemainingAmountLeaderReward)
	}
	return
}

func LeaderFetchTierInfoList(ctx comcontext.Context) (infoList []LeaderTierInfo, err error) {
	infoList = make([]LeaderTierInfo, 0, len(LeaderTierMetaList))
	for _, tierMeta := range LeaderTierMetaList {
		tierInfo, err := LeaderFetchTierInfo(ctx, tierMeta)
		if err != nil {
			return nil, err
		}
		if len(tierInfo.UIDs) == 0 {
			continue
		}
		infoList = append(infoList, tierInfo)
	}
	return
}

func LeaderFetchTierInfo(ctx comcontext.Context, tierMeta meta.TierMeta) (info LeaderTierInfo, err error) {
	tierRate, err := leaderGetTierDefaultRate(tierMeta)
	if err != nil {
		return
	}

	users, err := usermod.GetUsersByTier(tierMeta, 0)
	if err != nil {
		return
	}
	if config.Debug && len(users) > 10 {
		comlogging.GetLogger().
			WithType(constants.LogTypeBonusPoolLeader).
			WithContext(ctx).
			Debugf("shrink user list of tier `%v` to 10", tierMeta.Code)
		users = users[:10]
	}
	uids := make([]meta.UID, len(users))
	for i, user := range users {
		uids[i] = user.ID
	}
	validUIDs, err := filterValidTierUIDs(ctx, tierMeta, uids)
	if err != nil {
		return
	}
	sort.Slice(validUIDs, func(i, j int) bool { return validUIDs[i] < validUIDs[j] })
	info = LeaderTierInfo{
		TierType: tierMeta.ID,
		Rate:     tierRate,
		UIDs:     validUIDs,
	}
	return info, nil
}

func LeaderCalcExecutionHash(
	fromDate, toDate string,
	amount decimal.Decimal, tierInfoList []LeaderTierInfo,
) string {
	hashMeta := LeaderExecutionHashMeta{
		Secret:       config.SecretKey,
		FromDate:     fromDate,
		ToDate:       toDate,
		TotalAmount:  amount,
		TierInfoList: tierInfoList,
	}
	hashBytes, err := msgpack.Marshal(hashMeta)
	comutils.PanicOnError(err)
	return comutils.HashSha256Hex(hashBytes)
}

func LeaderCreateExecution(
	ctx comcontext.Context,
	hash string, fromDate, toDate string,
	amount decimal.Decimal, tierInfoList []LeaderTierInfo,
) (execution models.LeaderBonusPoolExecution, err error) {
	fromTime, toTime, err := leaderValidateDateRange(fromDate, toDate)
	if err != nil {
		return
	}

	logger := comlogging.GetLogger().
		WithType(constants.LogTypeBonusPoolLeader).
		WithContext(ctx)

	execHash := LeaderCalcExecutionHash(fromDate, toDate, amount, tierInfoList)
	if hash != execHash {
		logger.
			WithFields(logrus.Fields{
				"expected_hash": execHash,
				"actual_hash":   hash,
			}).
			Debug("create execution hash mismatch")
		if !config.Debug {
			err = utils.WrapError(constants.ErrorInvalidParams)
			return
		}
		logger.
			WithFields(logrus.Fields{
				"expected_hash": execHash,
				"actual_hash":   hash,
			}).
			Debug("bypass a mismatch hash")
	}

	lock, err := lockmod.LockSimple("bonus_pool:leader:create")
	if err != nil {
		return
	}
	defer lock.Unlock()

	execution = models.LeaderBonusPoolExecution{
		Hash:        hash,
		FromDate:    dbfields.NewDateFieldFromTime(fromTime),
		ToDate:      dbfields.NewDateFieldFromTime(toTime),
		TotalAmount: amount,
		Status:      constants.BonusPoolLeaderExecutionStatusInit,
	}
	totalReceiverCount := 0
	for _, tierInfo := range tierInfoList {
		totalReceiverCount += len(tierInfo.UIDs) + len(tierInfo.AdditionalUserInfoList)
		switch tierInfo.TierType {
		case constants.TierMetaSeniorPartner.ID:
			execution.TierSeniorRate = tierInfo.Rate
			execution.TierSeniorReceiverCount = uint16(len(tierInfo.UIDs))
			break
		case constants.TierMetaRegionalPartner.ID:
			execution.TierRegionalRate = tierInfo.Rate
			execution.TierRegionalReceiverCount = uint16(len(tierInfo.UIDs))
			break
		case constants.TierMetaGlobalPartner.ID:
			execution.TierGlobalRate = tierInfo.Rate
			execution.TierGlobalReceiverCount = uint16(len(tierInfo.UIDs))
			break
		default:
			err = utils.IssueErrorf(
				"leader bonus pool cannot process an unknown tier type %v",
				tierInfo.TierType)
			return
		}
	}
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		var overlapExec models.LeaderBonusPoolExecution
		err = dbTxn.
			Where(dbquery.Lte(models.LeaderBonusPoolExecutionColFromDate, execution.ToDate)).
			Where(dbquery.Gte(models.LeaderBonusPoolExecutionColToDate, execution.FromDate)).
			First(&overlapExec).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if overlapExec.ID > 0 {
			return meta.NewMessageError(
				"Leader bonus Pool already has an existing range %v --> %v with id=%v",
				overlapExec.FromDate, overlapExec.ToDate, overlapExec.ID,
			)
		}
		if err = dbTxn.Create(&execution).Error; err != nil {
			return utils.WrapError(err)
		}
		var (
			details = make([]models.LeaderBonusPoolDetail, 0, totalReceiverCount)
			now     = time.Now()
		)
		for _, tierInfo := range tierInfoList {
			for _, uid := range tierInfo.UIDs {
				detail := models.LeaderBonusPoolDetail{
					ExecutionID: execution.ID,
					TierType:    tierInfo.TierType,
					UID:         uid,
					CreateTime:  now.Unix(),
					UpdateTime:  now.Unix(),
				}
				details = append(details, detail)
			}
			for _, info := range tierInfo.AdditionalUserInfoList {
				detail := models.LeaderBonusPoolDetail{
					ExecutionID: execution.ID,
					TierType:    tierInfo.TierType,
					UID:         info.UID,
					Note:        info.Note,
					CreateTime:  now.Unix(),
					UpdateTime:  now.Unix(),
				}
				details = append(details, detail)
			}
		}
		if err = dbTxn.CreateInBatches(details, 200).Error; err != nil {
			return utils.WrapError(err)
		}
		return
	})
	if err != nil {
		return
	}
	comlogging.GetLogger().
		WithType(constants.LogTypeBonusPoolLeader).
		WithContext(ctx).
		WithFields(logrus.Fields{
			"id":   execution.ID,
			"hash": execution.Hash,
		}).
		Info("created an execution")
	return
}

func LeaderRunExecution(
	ctx comcontext.Context, executionID uint32,
) (execution models.LeaderBonusPoolExecution, err error) {
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		return database.Atomic(ctx, database.AliasWalletMaster, func(dbWalletTxn *gorm.DB) (err error) {
			leaderDB := LeaderGetDB(dbTxn)
			execution, err = leaderDB.GetAndLockExecution(
				executionID,
				constants.BonusPoolLeaderExecutionStatusInit)
			if err != nil {
				return
			}
			details, err := leaderDB.GetExecutionDetails(execution.ID)
			if err != nil {
				return utils.WrapError(err)
			}

			orders := make([]models.Order, 0, len(details))
			for _, detail := range details {
				order, genErr := leaderGenOrder(execution, detail)
				if genErr != nil {
					err = genErr
					return
				}
				orders = append(orders, order)
			}
			if err = dbWalletTxn.CreateInBatches(orders, 200).Error; err != nil {
				return utils.WrapError(err)
			}
			execution.Status = constants.BonusPoolLeaderExecutionStatusExecuting
			if err = dbTxn.Save(&execution).Error; err != nil {
				return utils.WrapError(err)
			}
			return
		})
	})
	if err != nil {
		return
	}
	comlogging.GetLogger().
		WithType(constants.LogTypeBonusPoolLeader).
		WithContext(ctx).
		WithFields(logrus.Fields{
			"id":   execution.ID,
			"hash": execution.Hash,
		}).
		Info("run an execution")
	return
}

func LeaderTryCompleteExecution(ctx comcontext.Context, executionID uint32) (
	execution models.LeaderBonusPoolExecution, err error,
) {
	logger := comlogging.GetLogger().
		WithType(constants.LogTypeBonusPoolLeader).
		WithContext(ctx)
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		leaderDB := LeaderGetDB(dbTxn)
		execution, err = leaderDB.GetAndLockExecution(
			executionID,
			constants.BonusPoolLeaderExecutionStatusExecuting)
		if err != nil {
			return
		}
		details, err := leaderDB.GetExecutionDetails(execution.ID)
		if err != nil {
			return utils.WrapError(err)
		}
		for _, detail := range details {
			if detail.OrderID == 0 {
				return nil
			}
		}

		execution.Status = constants.BonusPoolLeaderExecutionStatusCompleted
		if err = dbTxn.Save(&execution).Error; err != nil {
			return utils.WrapError(err)
		}
		logger.
			WithFields(logrus.Fields{
				"id":   execution.ID,
				"hash": execution.Hash,
			}).
			Info("complete an execution")
		return nil
	})
	return
}
