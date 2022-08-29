package middleware

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type JwtSimpleMapClaims struct {
	jwt.MapClaims
}

func JwtNewSimpleMapClaims() JwtSimpleMapClaims {
	return JwtSimpleMapClaims{jwt.MapClaims{}}
}

func (this JwtSimpleMapClaims) Valid() error {
	return this.MapClaims.Valid(nil)
}
