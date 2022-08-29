package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	walletv1 "gitlab.com/snap-clickstaff/torque-go/api/services/wallet/v1"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func TestGetOverview(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioGetOverview(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestGetCurrencyTorque(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	body := comutils.JsonEncodeF(walletv1.PortfolioGetCurrencyRequest{
		Currency: constants.CurrencyTorque,
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioGetCurrency(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestGetCurrencyBitcoin(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	body := comutils.JsonEncodeF(walletv1.PortfolioGetCurrencyRequest{
		Currency: constants.CurrencyBitcoin,
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioGetCurrency(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestListCurrencyOrderTorque(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	body := comutils.JsonEncodeF(walletv1.PortfolioListOrdersRequest{
		Currency: constants.CurrencyTorque,
		Paging:   meta.Paging{Limit: 10},
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioListCurrencyOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
}

func TestListCurrencyOrderBitcoin(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	body := comutils.JsonEncodeF(walletv1.PortfolioListOrdersRequest{
		Currency: constants.CurrencyBitcoin,
		Paging:   meta.Paging{Limit: 10},
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioListCurrencyOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
	t.Logf("Response:\n%s", readResponseBodyString(rec))
}

func TestListCurrencyOrderRipple(t *testing.T) {
	e := newEcho()
	rec := httptest.NewRecorder()

	body := comutils.JsonEncodeF(walletv1.PortfolioListOrdersRequest{
		Currency: constants.CurrencyRipple,
		Paging:   meta.Paging{Limit: 10},
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := patchContext(e.NewContext(req, rec))

	if !assert.NoError(t, walletv1.PortfolioListCurrencyOrder(c)) {
		return
	}

	assertSuccessResponse(t, rec)
	t.Logf("Response:\n%s", readResponseBodyString(rec))
}
