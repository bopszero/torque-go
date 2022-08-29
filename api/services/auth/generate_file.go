package auth

import (
	"bufio"
	"encoding/csv"
	"os"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/kycmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
)

func KycDOBNotMatchCsv(output string) {
	comutils.EchoWithTime("Generate CSV file started ... ")
	var kycRequestList []models.KycRequest
	err := database.GetDbF(database.AliasWalletSlave).
		Where(dbquery.Equal(models.KycRequestColStatus, constants.KycRequestStatusApproved)).
		Find(&kycRequestList).
		Error
	comutils.PanicOnError(err)

	csvFile, err := os.Create(output)
	comutils.PanicOnError(err)
	csvWriter := csv.NewWriter(bufio.NewWriter(csvFile))
	err = csvWriter.Write([]string{
		"kyc_request_id",
		"uid",
		"user_code",
		"username",
		"original_email",
		"full_name",
		"create_time",
		"dob_user",
		"dob_jumio",
	})
	comutils.PanicOnError(err)

	ctx := comcontext.NewContext()
	for _, request := range kycRequestList {
		scanDetails, err := kycmod.GetJumioScanDetails(ctx, request.Reference)
		comutils.PanicOnError(err)
		if request.DOB.Format(constants.DateFormatISO) != scanDetails.Document.DOB {
			user, err := usermod.GetUserFast(request.UID)
			comutils.PanicOnError(err)
			err = csvWriter.Write([]string{
				comutils.Stringify(request.ID),
				comutils.Stringify(request.UID),
				request.UserCode,
				user.Username,
				request.EmailOriginal,
				request.FullName,
				comutils.Stringify(request.CreateTime),
				request.DOB.Format(constants.DateFormatISO),
				scanDetails.Document.DOB,
			})
			comutils.PanicOnError(err)
		}
	}
	csvWriter.Flush()
	comutils.PanicOnError(csvFile.Close())

	comutils.EchoWithTime("Generate CSV file finished")
}
