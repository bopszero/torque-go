package usermod

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func GetUserByModel(queryModel *models.User) (user models.User, err error) {
	err = database.GetDbSlave().
		First(&user, queryModel).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = constants.ErrorUserNotFound
		}
		err = utils.WrapError(err)
	}
	return
}

func GetUser(uid meta.UID) (user models.User, err error) {
	return GetUserByModel(&models.User{ID: uid})
}

func GetUserFast(uid meta.UID) (user models.User, err error) {
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		fmt.Sprintf("user:%v", uid),
		30*time.Second,
		&user,
		func() (interface{}, error) {
			return GetUser(uid)
		},
	)
	return
}

func GetUserFastF(uid meta.UID) models.User {
	user, err := GetUserFast(uid)
	comutils.PanicOnError(err)
	return user
}

func GetUserByUsername(username string) (user models.User, err error) {
	return GetUserByModel(&models.User{Username: username})
}

func GetUserByReferralCode(referalCode string) (models.User, error) {
	return GetUserByModel(&models.User{ReferralCode: referalCode})
}

func IsValidUserTOTP(uid meta.UID, input string, require2FA bool) bool {
	if config.Test && input == constants.TestTwoFaCode {
		return true
	}

	user, err := GetUserFast(uid)
	if err != nil {
		return false
	}
	if user.TwoFaKey == "" {
		return !require2FA
	}

	return utils.IsValidateTOTP(input, user.TwoFaKey)
}

// GetUsersByTier do as its name (limit<=0 means unlimited).
func GetUsersByTier(tierMeta meta.TierMeta, limit int) (users []models.User, err error) {
	var (
		db        = database.GetDbSlave()
		queryCond = models.User{
			TierType:  tierMeta.ID,
			IsDeleted: models.NewBool(false),
		}
		pageLimit = limit
	)
	if pageLimit <= 0 || pageLimit > 1000 {
		pageLimit = 1000
	}
	usersItr := utils.NewChunkIterator(pageLimit, func(lastItem interface{}) (items interface{}, err error) {
		var (
			lastUID meta.UID
			users   []models.User
		)
		if lastItem != nil {
			lastUID = lastItem.(models.User).ID
		}
		err = db.
			Where(&queryCond).
			Where(dbquery.Gt(models.UserColID, lastUID)).
			Order(dbquery.OrderAsc(models.UserColID)).
			Limit(pageLimit).
			Find(&users).
			Error
		if err != nil {
			err = utils.WrapError(err)
			return
		}
		return users, nil
	})
	var (
		chunkUsers          []models.User
		iterNext, iterState = usersItr.GetNext(&chunkUsers)
	)
	for ; iterState.OK(); iterNext, iterState = iterNext(&chunkUsers) {
		users = append(users, chunkUsers...)
		if 0 < limit && limit <= len(users) {
			users = users[:limit]
			break
		}
	}
	err = iterState.Error()
	return
}
