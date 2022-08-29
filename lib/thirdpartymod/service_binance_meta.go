package thirdpartymod

import (
	"github.com/adshao/go-binance/v2/common"
)

const (
	BinanceErrorCodeOrderFound = -2013
)

func IsBinanceError(err error) bool {
	switch err.(type) {
	case common.APIError:
		return true
	default:
		return false
	}
}

func IsBinanceErrorCode(err error, code int64) bool {
	switch errT := err.(type) {
	case *common.APIError:
		return errT.Code == code
	default:
		return false
	}
}
