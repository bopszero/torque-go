package kycmod

import (
	"strings"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/settingmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetJumioScanDetails(ctx comcontext.Context, jumioIdScanReference string) (*JumioScanDetailsResponse, error) {
	jumioClient, err := thirdpartymod.GetJumioServiceSystemClient()
	if err != nil {
		return nil, err
	}
	resp, err := jumioClient.GetScanDetailsResp(ctx, jumioIdScanReference)
	if err != nil {
		return nil, err
	}
	var scanDetailsResp JumioScanDetailsResponse
	if err = comutils.JsonDecode(resp.String(), &scanDetailsResp); err != nil {
		return nil, err
	}
	return &scanDetailsResp, nil
}

func GetJumioVerificationData(ctx comcontext.Context, jumioIdScanReference string) (*JumioDataVerificationResponse, error) {
	jumioClient, err := thirdpartymod.GetJumioServiceSystemClient()
	if err != nil {
		return nil, err
	}
	resp, err := jumioClient.GetVerificationDataResp(ctx, jumioIdScanReference)
	if err != nil {
		return nil, err
	}
	var verificationResp JumioDataVerificationResponse
	if err = comutils.JsonDecode(resp.String(), &verificationResp); err != nil {
		return nil, err
	}
	return &verificationResp, nil
}

func GetJumioIpWhiteList() []string {
	jumioWhiteListIPStr, err := settingmod.GetSettingValueFast(constants.SettingKeyKycJumioWhiteListIP)
	if err != nil {
		return nil
	}
	jumioWhiteListIP := strings.Split(jumioWhiteListIPStr, ",")
	return jumioWhiteListIP
}

func initJumioScan(ctx comcontext.Context, scanRef string) (scan models.JumioScan, err error) {
	now := time.Now()
	scan = models.JumioScan{
		Reference: scanRef,
		Status:    JumioScanStatusInit,

		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		if err := dbTxn.Save(&scan).Error; err != nil && !database.IsDuplicateEntryError(err) {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func UpdateJumioScan(ctx comcontext.Context, scanRef string) (scan models.JumioScan, err error) {
	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(&scan, &models.JumioScan{Reference: scanRef}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		if scan.Status != JumioScanStatusInit {
			return
		}

		now := time.Now()
		if scan.CreateTime > now.Add(-JumioScanInitLifeTime).Unix() {
			scanDetailsResp, err := GetJumioScanDetails(ctx, scan.Reference)
			if err != nil {
				return err
			}
			scan.UserCode = scanDetailsResp.Transaction.UserCode
			scan.RequestCode = scanDetailsResp.Transaction.KycCode
			scan.Status = JumioScanStatusUpdatedData
		} else {
			scan.Status = JumioScanStatusExpired
		}
		scan.UpdateTime = now.Unix()

		if err := dbTxn.Save(&scan).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}
