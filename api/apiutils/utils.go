package apiutils

import (
	"github.com/dgrijalva/jwt-go/v4"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func GetContextJWT(ctx comcontext.Context) *jwt.Token {
	jwTokenObj := ctx.Get(config.ContextKeyJWT)
	if jwTokenObj == nil {
		return nil
	}
	return jwTokenObj.(*jwt.Token)
}

func getContextUID(ctx comcontext.Context) (meta.UID, error) {
	jwToken := GetContextJWT(ctx)
	if jwToken != nil {
		exClaims, err := authmod.ParseExClaim(jwToken.Claims)
		if err != nil {
			return 0, err
		}
		return exClaims.UID, nil
	}

	if config.Debug {
		debugUID := ctx.Get(config.ContextKeyDebugUID)
		if debugUID != nil {
			return meta.UID(debugUID.(uint64)), nil
		}
	}

	return 0, constants.ErrorDataNotFound
}

func GetContextUID(ctx comcontext.Context) (meta.UID, error) {
	cachedUID := ctx.Get(config.ContextKeyUID)
	if cachedUID != nil {
		return cachedUID.(meta.UID), nil
	}

	uid, err := getContextUID(ctx)
	if err != nil {
		return 0, err
	}

	ctx.Set(config.ContextKeyUID, uid)
	return uid, nil
}

func GetContextUidF(ctx comcontext.Context) meta.UID {
	uid, err := GetContextUID(ctx)
	comutils.PanicOnError(err)

	return uid
}
