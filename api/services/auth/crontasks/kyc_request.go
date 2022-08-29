package crontasks

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func KycUpdateJumioScans() {
	defer apiutils.CronRunWithRecovery("KycUpdateJumioScans", nil)

	var (
		ctx    = comcontext.NewContext()
		logger = comlogging.GetLogger()

		err           error
		jumioScanList []models.JumioScan
	)
	err = database.GetDbF(database.AliasWalletSlave).
		Where(dbquery.Equal(models.JumioScanColStatus, kycmod.JumioScanStatusInit)).
		Find(&jumioScanList).
		Error
	comutils.PanicOnError(err)

	for _, jumioScan := range jumioScanList {
		if _, err := kycmod.UpdateJumioScan(ctx, jumioScan.Reference); err != nil {
			logger.
				WithError(err).
				WithField("scan_ref", jumioScan.Reference).
				Errorf("update jumio scan failed | err=%s", err.Error())
		}
	}
}

func KycExecuteRequests() {
	defer apiutils.CronRunWithRecovery("KycExecuteRequests", nil)

	var (
		ctx    = comcontext.NewContext()
		logger = comlogging.GetLogger()

		err            error
		kycRequestList []models.KycRequest
	)
	err = database.GetDbF(database.AliasWalletSlave).
		Where(dbquery.In(
			models.KycRequestColStatus,
			[]meta.KycRequestStatus{
				constants.KycRequestStatusInit,
				constants.KycRequestStatusPendingAnalysis,
			})).
		Find(&kycRequestList).
		Error
	comutils.PanicOnError(err)

	for _, request := range kycRequestList {
		if _, err := kycmod.ExecuteRequest(ctx, request.ID); err != nil {
			logger.
				WithError(err).
				WithField("request_id", request.ID).
				Errorf("execute kyc request failed | err=%s", err.Error())
		}
	}
	return
}
