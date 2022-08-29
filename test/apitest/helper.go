package test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func newEcho() *echo.Echo {
	return api.CreateEchoObject()
}

func patchContext(c echo.Context) echo.Context {
	c.Set(config.ContextKeyDebugUID, ApiUID)

	return c
}

func assertSuccessResponse(t *testing.T, rec *httptest.ResponseRecorder) {
	assert.Truef(t, 200 <= rec.Code && rec.Code < 300, "status code `%s` means fail", rec.Code)

	responseCode := rec.Header().Get(config.HttpHeaderApiResponseCode)
	assert.Equal(
		t,
		ErrorCodeSuccess, responseCode,
		"response code `%s` means fail", responseCode,
	)
}

func readResponseBodyString(rec *httptest.ResponseRecorder) string {
	bodyBytes, err := ioutil.ReadAll(rec.Body)
	comutils.PanicOnError(err)

	return string(bodyBytes)
}
