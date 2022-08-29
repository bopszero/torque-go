package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/ipmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func MetaStatus(c echo.Context) (err error) {
	var (
		ctx                     = apiutils.EchoWrapContext(c)
		isEnableBanIp           = ipmod.IsEnableBanIp(ctx)
		metaStatusResponseAllow = MetaStatusResponse{
			Registration{
				IsBanned: false,
			},
		}
	)
	if !isEnableBanIp {
		return responses.Ok(ctx, metaStatusResponseAllow)
	}

	userIP := utils.GetRequestUserIP(c.Request())
	whiteListIp, err := ipmod.GetWhiteListIp()
	if err != nil {
		return err
	}
	if whiteListIp.Contains(userIP) {
		return responses.Ok(ctx, metaStatusResponseAllow)
	}

	ipCountry, err := ipmod.GetIpCountry(userIP)
	if err != nil {
		return err
	}
	var userCountry models.Country
	err = database.GetDbF(database.AliasMainSlave).
		First(&userCountry, &models.Country{CodeIso2: ipCountry}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			comlogging.GetLogger().
				WithContext(ctx).
				WithError(err).
				WithFields(
					logrus.Fields{
						"iso_code":   ipCountry,
						"request_ip": userIP,
					},
				).
				Errorf("cannot find code iso2 of IP country | err=%s", err.Error())
			return responses.Ok(ctx, metaStatusResponseAllow)
		}
		return utils.WrapError(err)
	}

	metaStatusResponseBanned := MetaStatusResponse{
		Registration{
			BannedInfo: BannedInfo{
				IP:      userIP,
				Country: userCountry.Name,
				Message: constants.ErrorBannedCountry.Message(ctx),
			},
			IsBanned: true,
		},
	}
	blackListIp, err := ipmod.GetBlackListIp()
	if err != nil {
		return utils.WrapError(err)
	}
	if blackListIp.Contains(userIP) {
		return responses.Ok(ctx, metaStatusResponseBanned)
	}
	if userCountry.IsBanned.Bool {
		return responses.Ok(ctx, metaStatusResponseBanned)
	}
	return responses.Ok(ctx, metaStatusResponseAllow)
}
