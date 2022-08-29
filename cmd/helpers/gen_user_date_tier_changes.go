package helpers

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GenerateDateUserTierChanges(date string, clearDateBeforeRun bool) {
	comutils.EchoWithTime("gen_user_date_tier_changes|%s|started.", date)

	dateTime, err := time.ParseInLocation(constants.DateFormatISO, date, time.Local)
	comutils.PanicOnError(err)

	var tree *affiliate.Tree
	timeGap := time.Now().Sub(dateTime)
	if timeGap <= 0 {
		panic(constants.ErrorInvalidData)
	}

	comutils.EchoWithTime("gen_user_date_tier_changes|%s|generating_tree...", date)

	if timeGap < 24*time.Hour {
		tree = affiliate.GenerateTreeAtNow()
	} else {
		tree = affiliate.GenerateTreeAtDate(date)
	}

	var userTierChanges []models.UserTierChange
	for uid, node := range tree.GetNodeMap() {
		userCurrentTierMeta, ok := constants.TierMetaMap[node.User.TierType]
		if !ok {
			continue
		}

		realTier := node.GetTier()
		if realTier.Meta.Code == userCurrentTierMeta.Code {
			continue
		}

		comutils.EchoWithTime(
			"gen_user_date_tier_changes|%s|see_user_tier_change|uid=%d,to_tier=%s",
			date, uid, realTier.Meta.Code,
		)

		extraDataJSON := comutils.JsonEncodeF(genTierExtraData(realTier))
		userTierChange := models.UserTierChange{
			Date:         date,
			UID:          uid,
			Status:       constants.TierChangeStatusReady,
			FromTierType: userCurrentTierMeta.ID,
			ToTierType:   realTier.Meta.ID,

			ExtraData:  extraDataJSON,
			CreateTime: time.Now().Unix(),
		}
		userTierChanges = append(userTierChanges, userTierChange)
	}

	if len(userTierChanges) == 0 {
		comutils.EchoWithTime("gen_user_date_tier_changes|%s|no_data_to_save...", date)
		return
	}

	comutils.EchoWithTime("gen_user_date_tier_changes|%s|saving_data...", date)

	ctx := comcontext.NewContext()
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		if clearDateBeforeRun {
			comutils.PanicOnError(
				dbTxn.
					Delete(&models.UserTierChange{}, &models.UserTierChange{Date: date}).
					Error,
			)
		}
		if err := dbTxn.CreateInBatches(userTierChanges, 200).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	comutils.PanicOnError(err)

	comutils.EchoWithTime("gen_user_date_tier_changes|%s|finished.", date)
}

func genTierExtraData(nodeTier affiliate.TreeNodeTier) meta.O {
	extraData := meta.O{}

	baseUserInfoList := make([]meta.O, len(nodeTier.BaseNodes))
	for idx, baseNode := range nodeTier.BaseNodes {
		nodeInfo := meta.O{
			"uid":      baseNode.User.ID,
			"username": baseNode.User.Username,
			"tier":     baseNode.GetTier().Meta.Code,
		}

		baseUserInfoList[idx] = nodeInfo
	}
	if len(baseUserInfoList) == 0 {
		return extraData
	}

	extraData["base_on_users"] = baseUserInfoList
	extraData["is_fixed"] = nodeTier.IsFixed
	return extraData
}
