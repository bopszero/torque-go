package rewardmod

import (
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetTotalCleanUIDs() (uids []meta.UID, err error) {
	kycUIDs, err := GetTotalKycUIDs()
	if err != nil {
		return
	}
	whitelistUIDS, err := GetWhitelistUIDs()
	if err != nil {
		return
	}
	blacklistUIDS, err := GetBlacklistUIDs()
	if err != nil {
		return
	}
	kycSet, err := comtypes.NewHashSetFromList(kycUIDs)
	comutils.PanicOnError(err)
	whitelistSet, err := comtypes.NewHashSetFromList(whitelistUIDS)
	comutils.PanicOnError(err)
	blacklistSet, err := comtypes.NewHashSetFromList(blacklistUIDS)
	comutils.PanicOnError(err)

	err = kycSet.Union(whitelistSet).Diff(blacklistSet).DumpItems(&uids)
	return
}

func GetTotalKycUIDs() (_ []meta.UID, err error) {
	var (
		dbMain   = database.GetDbSlave()
		dbWallet = database.GetDbF(database.AliasWalletSlave)
	)
	var kycRequests []models.KycRequest
	err = dbWallet.
		Model(&models.KycRequest{}).
		Where(dbquery.In(
			models.KycRequestColStatus,
			[]meta.KycRequestStatus{
				constants.KycRequestStatusPendingApproval,
				constants.KycRequestStatusApproved,
			},
		)).
		Group(models.KycRequestColEmail).
		Select(models.KycRequestColEmail).
		Find(&kycRequests).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	kycEmails := make([]string, len(kycRequests))
	for i, req := range kycRequests {
		kycEmails[i] = req.EmailOriginal
	}

	var (
		kycUidSet    = make(comtypes.HashSet, len(kycEmails))
		kycEmailsItr = utils.NewChunkIterator(1000, func(lastItem interface{}) (items interface{}, err error) {
			if lastItem != nil {
				return nil, nil
			}
			return kycEmails, nil
		})
		chunkEmails         []string
		iterNext, iterState = kycEmailsItr.GetNext(&chunkEmails)
	)
	for ; iterState.OK(); iterNext, iterState = iterNext(&chunkEmails) {
		var chunkUsers []models.User
		err = dbMain.
			Select(models.UserColID).
			Where(dbquery.In(models.UserColEmailOriginal, chunkEmails)).
			Find(&chunkUsers).
			Error
		if err != nil {
			err = utils.WrapError(err)
			return
		}
		for _, user := range chunkUsers {
			kycUidSet.Add(user.ID)
		}
	}
	if iterState.Error() != nil {
		err = iterState.Error()
		return
	}
	var kycPendingUserInfoList []models.KycSpecialUserInfo
	err = dbWallet.
		Select(models.CommonColUID).
		Where(&models.KycSpecialUserInfo{IsPending: models.NewBool(true)}).
		Find(&kycPendingUserInfoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	for _, userInfo := range kycPendingUserInfoList {
		kycUidSet.Add(userInfo.UID)
	}

	uids := make([]meta.UID, 0, len(kycUidSet))
	for uidObj := range kycUidSet {
		uids = append(uids, uidObj.(meta.UID))
	}
	return uids, nil
}

func FilterKycUIDs(uids []meta.UID) (_ []meta.UID, err error) {
	var (
		dbMain   = database.GetDbSlave()
		dbWallet = database.GetDbF(database.AliasWalletSlave)
	)
	var users []models.User
	err = dbMain.
		Where(dbquery.In(models.UserColID, uids)).
		Find(&users).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	var (
		emailMap    = make(map[meta.UID]string, len(users))
		inputEmails = make([]string, len(users))
	)
	for i, user := range users {
		emailMap[user.ID] = user.OriginalEmail
		inputEmails[i] = user.OriginalEmail
	}
	var (
		kycRequests            []models.KycRequest
		kycPendingUserInfoList []models.KycSpecialUserInfo
	)
	err = dbWallet.
		Model(&models.KycRequest{}).
		Where(dbquery.In(
			models.KycRequestColStatus,
			[]meta.KycRequestStatus{
				constants.KycRequestStatusPendingApproval,
				constants.KycRequestStatusApproved,
			},
		)).
		Where(dbquery.In(models.KycRequestColEmail, inputEmails)).
		Group(models.KycRequestColEmail).
		Select(models.KycRequestColEmail).
		Find(&kycRequests).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	err = dbWallet.
		Select(models.CommonColUID).
		Where(dbquery.In(models.CommonColUID, uids)).
		Where(&models.KycSpecialUserInfo{IsPending: models.NewBool(true)}).
		Find(&kycPendingUserInfoList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	kycEmailSet := make(comtypes.HashSet, len(kycRequests))
	for _, req := range kycRequests {
		kycEmailSet.Add(req.EmailOriginal)
	}
	pendingUidSet := make(comtypes.HashSet, len(kycPendingUserInfoList))
	for _, userInfo := range kycPendingUserInfoList {
		pendingUidSet.Add(userInfo.UID)
	}
	validUIDs := make([]meta.UID, 0, len(uids))
	for _, uid := range uids {
		email, ok := emailMap[uid]
		if (ok && kycEmailSet.Contains(email)) || pendingUidSet.Contains(uid) {
			validUIDs = append(validUIDs, uid)
		}
	}
	return validUIDs, nil
}

func getRewardUIDs(rewardType meta.RewardType) (uids []meta.UID, err error) {
	var configList []models.UserConfig
	err = database.GetDbSlave().
		Where(&models.UserConfig{RewardType: rewardType}).
		Find(&configList).
		Error
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	uids = make([]meta.UID, len(configList))
	for i, conf := range configList {
		uids[i] = conf.UID
	}
	return
}

func GetWhitelistUIDs() (uids []meta.UID, err error) {
	return getRewardUIDs(constants.RewardTypeWhitelist)
}

func GetBlacklistUIDs() (uids []meta.UID, err error) {
	return getRewardUIDs(constants.RewardTypeBlacklist)
}

// FilterCleanUIDs provide clean identity users
func FilterCleanUIDs(uids []meta.UID) (_ []meta.UID, err error) {
	if len(uids) == 0 {
		return uids, nil
	}
	uidSet, err := comtypes.NewHashSetFromList(uids)
	comutils.PanicOnError(err)

	blacklistUIDS, err := GetBlacklistUIDs()
	if err != nil {
		return
	}
	blacklistSet, err := comtypes.NewHashSetFromList(blacklistUIDS)
	comutils.PanicOnError(err)

	uidSet = uidSet.Diff(blacklistSet)
	var cleanUIDs []meta.UID
	comutils.PanicOnError(
		uidSet.DumpItems(&cleanUIDs),
	)
	return FilterKycUIDs(cleanUIDs)
}

func FilterEnoughBalanceUIDs(uids []meta.UID) (_ []meta.UID, err error) {
	if len(uids) == 0 {
		return uids, nil
	}

	userBalanceMapMap, err := tradingbalance.GenUserCoinBalanceMapMap(uids, false)
	if err != nil {
		return
	}
	var (
		minThresholdMap = GetCurrencyDepositMinThresholdMap()
		validUIDs       = make([]meta.UID, 0, len(uids))
	)
	for _, uid := range uids {
		balanceMap, ok := userBalanceMapMap[uid]
		if !ok {
			continue
		}
		for currency, threshold := range minThresholdMap {
			if balanceMap[currency].GreaterThanOrEqual(threshold) {
				validUIDs = append(validUIDs, uid)
				break
			}
		}
	}
	return validUIDs, nil
}

func FilterValidUIDs(uids []meta.UID) (_ []meta.UID, err error) {
	if len(uids) == 0 {
		return uids, nil
	}
	cleanUIDs, err := FilterCleanUIDs(uids)
	if err != nil {
		return
	}
	return FilterEnoughBalanceUIDs(cleanUIDs)
}
