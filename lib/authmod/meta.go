package authmod

import (
	"github.com/dgrijalva/jwt-go/v4"
)

const (
	JwtRefreshLifeTimeMinRate = 0.5
)

type JwtKeyPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	JwtParser = jwt.NewParser(
		jwt.WithValidMethods(
			[]string{
				jwt.SigningMethodHS256.Alg(),
				jwt.SigningMethodHS512.Alg(),
			},
		),
		jwt.WithJSONNumber(),
	)
)
