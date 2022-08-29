package utils

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.com/snap-clickstaff/torque-go/buildmeta"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

func NewRestyClient(timeout time.Duration) *resty.Client {
	client := resty.New()
	client.SetTimeout(timeout)
	client.SetHeader("User-Agent", fmt.Sprintf("torque-wallet/%v", buildmeta.Version))
	client.Debug = config.Debug

	return client
}
