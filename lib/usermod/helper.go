package usermod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func IsUserContext(ctx comcontext.Context) bool {
	return ctx.Get(config.ContextKeyUID) != nil
}
