package emailmod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
)

func Send(ctx comcontext.Context, email *Email) error {
	return thirdpartymod.GetEmailServiceSystemClient().
		Send(ctx, &email.Message)
}
