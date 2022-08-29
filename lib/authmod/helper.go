package authmod

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetAccessSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS256
}

func GetRefreshSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS512
}

func ParseExClaim(claims jwt.Claims) (exClaims JwtExClaims, err error) {
	exClaimsPtr, ok := claims.(*JwtExClaims)
	if ok {
		exClaims = *exClaimsPtr
	} else {
		err = utils.DumpDataByJSON(claims, &exClaims)
	}
	return
}

func GetJwtSystemIssuer() string {
	return viper.GetString(config.KeyAppName)
}
