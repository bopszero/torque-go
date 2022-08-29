package utils

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func IssueErrorf(format string, args ...interface{}) error {
	return comlogging.NewPreloadSentryError(fmt.Errorf(format, args...))
}

func WrapError(err error) error {
	if err == nil {
		panic("cannot wrap nil error")
	}

	return comlogging.NewPreloadSentryError(err)
}

func IsOurError(err error, ourError meta.ErrorCode) bool {
	if err == nil {
		return false
	}

	return err.Error() == ourError.String()
}

func ErrorCatchWithLog(ctx comcontext.Context, keyword string, err error) {
	if err != nil {
		comlogging.GetLogger().
			WithContext(ctx).
			WithError(WrapError(err)).
			WithField("keyword", keyword).
			Errorf("system catch an error `%v`", err)
	}

	panicErrLike := recover()
	if panicErrLike == nil {
		return
	}

	panicErr := WrapError(comutils.ToError(panicErrLike))
	comlogging.GetLogger().
		WithContext(ctx).
		WithError(panicErr).
		WithField("keyword", keyword).
		Errorf("system panic with error `%v`", panicErr)
}
