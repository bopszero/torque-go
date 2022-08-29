package meta

import (
	"fmt"
	"gitlab.com/snap-clickstaff/go-common/comlogging"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

type ErrorCode string

func (this ErrorCode) String() string {
	return string(this)
}

type GeneralError struct {
	code ErrorCode

	msgKey  string
	msgData O
}

func NewGeneralError(code ErrorCode) *GeneralError {
	return &GeneralError{
		code: code,
	}
}

func (this GeneralError) Error() string {
	return this.code.String()
}

func (this *GeneralError) Code() ErrorCode {
	return this.code
}

func (this *GeneralError) WithData(data O) *GeneralError {
	clone := NewGeneralError(this.code)
	clone.msgData = data

	return clone
}

func (this *GeneralError) WithKey(key string) *GeneralError {
	return this.WithMessage(key, nil)
}

func (this *GeneralError) WithMessage(key string, data O) *GeneralError {
	clone := NewGeneralError(this.code)
	clone.msgKey = key
	clone.msgData = data

	return clone
}

func (this *GeneralError) Message(ctx comcontext.Context) string {
	var errKey string
	if this.msgKey == "" {
		errKey = this.code.String()
	} else {
		errKey = this.msgKey
	}

	message, err := comlocale.TranslateKeyData(ctx, errKey, this.msgData)
	if err == nil {
		return message
	}

	comlogging.GetLogger().
		WithContext(ctx).
		WithField("key", errKey).
		WithError(err).
		Warn("translation failed")
	if config.Test {
		message = fmt.Sprintf(
			"translate key `%s` failed | data=%s,err=%v",
			errKey, comutils.JsonEncodeF(this.msgData), err,
		)
		return message
	}
	switch err.(type) {
	case *i18n.MessageNotFoundErr:
		var (
			defaultCtx               = comcontext.NewContext()
			fallbackMsg, fallbackErr = comlocale.TranslateKeyData(defaultCtx, errKey, this.msgData)
		)
		if fallbackErr == nil {
			message = fallbackMsg
		}
		break
	}
	return message
}

type MessageError struct {
	msg string
}

func NewMessageError(msg string, params ...interface{}) *MessageError {
	return &MessageError{fmt.Sprintf(msg, params...)}
}

func (this MessageError) Error() string {
	return this.msg
}
