package notifmod

import (
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/httpmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func RegisterFirebaseToken(ctx comcontext.Context, uid meta.UID, token string, userAgent string) (
	_ *models.UserFirebaseToken, err error,
) {
	lock, err := lockmod.LockSimple("notif:token:register:%v", uid)
	if err != nil {
		return
	}
	defer lock.Unlock()

	var userFirebaseToken models.UserFirebaseToken
	err = database.GetDbSlave().
		First(
			&userFirebaseToken,
			&models.UserFirebaseToken{
				Token: token,
			}).
		Error
	if database.IsDbError(err) {
		return nil, utils.WrapError(err)
	}
	if userFirebaseToken.ID != 0 && userFirebaseToken.UID == uid {
		return &userFirebaseToken, nil
	}

	uaInfo := httpmod.UserAgentParse(userAgent)
	err = database.Atomic(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.
			First(
				&userFirebaseToken,
				&models.UserFirebaseToken{
					Token: token,
				}).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if userFirebaseToken.ID != 0 && userFirebaseToken.UID == uid {
			return nil
		}

		userFirebaseToken.UID = uid
		userFirebaseToken.Token = token
		userFirebaseToken.CreateTime = time.Now().Unix()
		userFirebaseToken.DeviceOS = uaInfo.Os.ToVersionString()
		userFirebaseToken.DeviceDesc = uaInfo.Device.ToString()
		if userFirebaseToken.DeviceOS == "" {
			userFirebaseToken.DeviceOS = "Unknown"
			if userAgent == "" {
				userFirebaseToken.DeviceDesc = "Unknown"
			} else {
				userFirebaseToken.DeviceDesc = userAgent
			}
		}

		if err = dbTxn.Save(&userFirebaseToken).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userFirebaseToken, nil
}
