package authmod

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type JwtExClaims struct {
	jwt.StandardClaims
	UID        meta.UID `json:"uid,omitempty"`
	Require2FA bool     `json:"2fa,omitempty"`
	NeedCommit bool     `json:"ncm,omitempty"`
}

func NewSystemJwtExClaims(uid meta.UID, timeout time.Duration) JwtExClaims {
	now := time.Now()
	return JwtExClaims{
		StandardClaims: jwt.StandardClaims{
			ID:        comutils.NewUuid4Code(),
			Issuer:    GetJwtSystemIssuer(),
			IssuedAt:  jwt.At(now),
			ExpiresAt: jwt.At(now.Add(timeout)),
		},
		UID: uid,
	}
}

func (this JwtExClaims) IsAuthorized() bool {
	return !this.NeedCommit && this.UID > 0
}

func (this JwtExClaims) CanRotate() bool {
	if config.Debug {
		return true
	}
	var (
		now         = time.Now()
		maxLifeTime = this.ExpiresAt.Time.Sub(this.IssuedAt.Time)
		lifeTime    = now.Sub(this.IssuedAt.Time)
	)
	return lifeTime.Seconds()/maxLifeTime.Seconds() > JwtRefreshLifeTimeMinRate
}

func (this JwtExClaims) RotateNew(timeout time.Duration) JwtExClaims {
	now := time.Now()
	this.ID = comutils.NewUuid4Code()
	this.IssuedAt = jwt.At(now)
	this.ExpiresAt = jwt.At(now.Add(timeout))
	return this
}
