package authmod

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func newJwtKeyRefreshToken(
	ctx comcontext.Context,
	claims JwtExClaims, deviceUID string,
) (token string, err error) {
	var (
		authConfig = GetAuthConfig()
		tokenInfo  = jwt.NewWithClaims(GetRefreshSigningMethod(), claims)
	)
	token, err = tokenInfo.SignedString(authConfig.RefreshSecret)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		jwtRefresh := models.UserJwtRefresh{
			ID:         claims.ID,
			UID:        claims.UID,
			DeviceUID:  deviceUID,
			ExpireTime: claims.ExpiresAt.Unix(),
		}
		if err := dbTxn.Create(&jwtRefresh).Error; err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	return
}

func GenJwtKeyCommit(user models.User) (token string, err error) {
	var (
		authConfig = GetAuthConfig()
		jwtClaims  = NewSystemJwtExClaims(user.ID, authConfig.AccessTimeout)
	)
	jwtClaims.NeedCommit = true
	jwtClaims.Require2FA = user.TwoFaKey != ""

	tokenInfo := jwt.NewWithClaims(GetAccessSigningMethod(), jwtClaims)
	token, err = tokenInfo.SignedString(authConfig.AccessSecret)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	return
}

func GenJwtKeyAccessToken(uid meta.UID) (token string, err error) {
	var (
		authConfig = GetAuthConfig()
		claims     = NewSystemJwtExClaims(uid, authConfig.AccessTimeout)
		tokenInfo  = jwt.NewWithClaims(GetAccessSigningMethod(), claims)
	)
	token, err = tokenInfo.SignedString(authConfig.AccessSecret)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	return
}

func GenJwtKeyRefreshToken(
	ctx comcontext.Context,
	uid meta.UID, deviceUID string,
) (token string, err error) {
	var (
		authConfig = GetAuthConfig()
		claims     = NewSystemJwtExClaims(uid, authConfig.RefreshTimeout)
	)
	return newJwtKeyRefreshToken(ctx, claims, deviceUID)
}

func GenJwtKeyPair(ctx comcontext.Context, uid meta.UID, deviceUID string) (pair JwtKeyPair, err error) {
	accessToken, err := GenJwtKeyAccessToken(uid)
	if err != nil {
		return
	}
	refreshToken, err := GenJwtKeyRefreshToken(ctx, uid, deviceUID)
	if err != nil {
		return
	}
	pair = JwtKeyPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

func RefreshJwtKeyPair(
	ctx comcontext.Context,
	tokenInfo *jwt.Token, deviceUID string, doRotate bool,
) (pair JwtKeyPair, err error) {
	exClaims, err := ParseExClaim(tokenInfo.Claims)
	if err != nil {
		return
	}

	var jwtRefresh models.UserJwtRefresh
	err = database.GetDbF(database.AliasWalletSlave).
		First(
			&jwtRefresh,
			&models.UserJwtRefresh{
				ID:        exClaims.ID,
				DeviceUID: deviceUID,
			},
		).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = constants.ErrorAuth
		}
		err = utils.WrapError(err)
		return
	}

	if pair.AccessToken, err = GenJwtKeyAccessToken(exClaims.UID); err != nil {
		return
	}
	if !doRotate {
		return
	}

	if !exClaims.CanRotate() {
		err = utils.WrapError(constants.ErrorInvalidParams)
		return
	}
	var (
		authConfig = GetAuthConfig()
		newClaims  = exClaims.RotateNew(authConfig.RefreshTimeout)
	)
	if pair.RefreshToken, err = newJwtKeyRefreshToken(ctx, newClaims, deviceUID); err != nil {
		return
	}

	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) error {
		err := dbTxn.
			Model(&jwtRefresh).
			Updates(&models.UserJwtRefresh{RotateTime: time.Now().Unix()}).
			Error
		if err != nil {
			return utils.WrapError(err)
		}
		return nil
	})
	if err != nil {
		return
	}

	return pair, nil
}
