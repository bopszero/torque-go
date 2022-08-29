package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api/middleware"
	"gitlab.com/snap-clickstaff/torque-go/lib/authmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func loginInputSetRefreshTokenCookie(c echo.Context, token string) {
	var (
		authConfig = authmod.GetAuthConfig()
		maxAge     = int(authConfig.RefreshTimeout.Seconds())
	)
	if token == "" {
		maxAge *= -1
	}
	c.SetCookie(&http.Cookie{
		Name:     middleware.JwtCookieRefreshKey,
		Value:    token,
		Path:     c.Echo().Reverse(RouteNameLoginRefresh),
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

func loginGetInputTokenLock(uid meta.UID) (*redsync.Mutex, error) {
	key := fmt.Sprintf("login:input:%v", uid)
	return lockmod.LockNoRetry(key, 5*time.Second)
}
