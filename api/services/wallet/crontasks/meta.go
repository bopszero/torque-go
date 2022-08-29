package crontasks

import (
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var httpClientMap map[time.Duration]*resty.Client

func init() {
	httpClientMap = make(map[time.Duration]*resty.Client)
}

func getHttpClient(timeout time.Duration) *resty.Client {
	httpClient, ok := httpClientMap[timeout]
	if !ok {
		httpClient = utils.NewRestyClient(timeout)
		httpClientMap[timeout] = httpClient
	}

	return httpClient
}
