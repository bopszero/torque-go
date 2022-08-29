package thirdpartymod

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func DumpJsonRequestHMAC(secret []byte, request *resty.Request) (err error) {
	var bodyJSON string
	switch bodyWithType := request.Body.(type) {
	case string:
		bodyJSON = bodyWithType
	default:
		bodyJSON, err = comutils.JsonEncode(bodyWithType)
		if err != nil {
			return
		}
	}

	nonce := comutils.RandomStringSimple(8)
	mac := comutils.HmacSha256(secret, []byte(bodyJSON+nonce))
	macHeader := fmt.Sprintf("HMAC 1.0:%s:%s", nonce, comutils.HexEncode(mac))

	request.
		SetHeader(echo.HeaderAuthorization, macHeader).
		SetHeader(echo.HeaderContentType, constants.ContentTypeJSON).
		SetBody(bodyJSON)
	return nil
}
