package v1

import (
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/usermod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func LoginInputTest(c echo.Context) error {
	var ctx = apiutils.EchoWrapContext(c)

	templateFile, err := template.ParseFiles("resources/auth/test_login.html")
	if err != nil {
		return err
	}
	return templateFile.Execute(ctx.Response(), nil)
}

func LoginInputPrepare(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel LoginInputPrepareRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	nonce, err := authmod.PasswordPrepareNonce(reqModel.Username)
	if err != nil {
		return err
	}

	responseModel := LoginInputPrepareResponse{
		NonceID: nonce.ID,
		Nonce:   comutils.HexEncode(nonce.Value),
		Salt:    comutils.HexEncode(nonce.Salt),
	}
	return responses.Ok(ctx, responseModel)
}

func LoginInputExecute(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel LoginInputExecuteRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	user, err := authmod.PasswordValidate(
		reqModel.NonceID,
		reqModel.Username, reqModel.PasswordEncryptedHex)
	if err != nil {
		return err
	}
	commitToken, err := authmod.GenJwtKeyCommit(user)
	if err != nil {
		return err
	}

	authConfig := authmod.GetAuthConfig()
	c.SetCookie(&http.Cookie{
		Name:     middleware.JwtCookieAccessKey,
		Value:    commitToken,
		Path:     c.Echo().Reverse(RouteNameLoginInputCommit),
		MaxAge:   int(authConfig.AccessTimeout.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	responseModel := LoginInputExecuteResponse{
		CommitToken: commitToken,
	}
	return responses.Ok(ctx, responseModel)
}

func LoginInputCommit(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		reqModel LoginInputCommitRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	lock, err := loginGetInputTokenLock(uid)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	if !usermod.IsValidUserTOTP(uid, reqModel.AuthCode, false) {
		return utils.WrapError(constants.ErrorAuthInput)
	}
	jwtKeyPair, err := authmod.GenJwtKeyPair(ctx, uid, reqModel.DeviceUID)
	if err != nil {
		return err
	}

	loginInputSetRefreshTokenCookie(c, jwtKeyPair.RefreshToken)

	responseModel := LoginInputCommitResponse{jwtKeyPair}
	return responses.Ok(ctx, responseModel)
}

func LoginRefresh(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		uid      = apiutils.GetContextUidF(ctx)
		jwtToken = apiutils.GetContextJWT(ctx)
		reqModel LoginRefreshRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	lock, err := loginGetInputTokenLock(uid)
	if err != nil {
		return err
	}
	defer lock.Unlock()

	jwtKeyPair, err := authmod.RefreshJwtKeyPair(ctx, jwtToken, reqModel.DeviceUID, reqModel.Rotate)
	if err != nil {
		return err
	}
	if jwtKeyPair.RefreshToken != "" {
		loginInputSetRefreshTokenCookie(c, jwtKeyPair.RefreshToken)
	}

	responseModel := LoginRefreshResponse{jwtKeyPair}
	return responses.Ok(ctx, responseModel)
}

func LoginLogout(c echo.Context) error {
	ctx := apiutils.EchoWrapContext(c)

	// TODO: Mark blacklist or something...

	loginInputSetRefreshTokenCookie(c, "")

	return responses.OkEmpty(ctx)
}
