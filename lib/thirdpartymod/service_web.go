package thirdpartymod

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var webServiceSystemClient = comtypes.NewSingleton(func() interface{} {
	return NewWebServiceClient(
		viper.GetString(config.KeyServiceWebMacSecret),
	)
})

type WebServiceClient struct {
	httpClient  *resty.Client
	secretBytes []byte
}

func NewWebServiceClient(secret string) *WebServiceClient {
	httpClient := utils.NewRestyClient(5 * time.Second)
	httpClient.SetHostURL(viper.GetString(config.KeyServiceWebBaseURL))

	return &WebServiceClient{
		httpClient:  httpClient,
		secretBytes: []byte(secret),
	}
}

func GetWebServiceSystemClient() *WebServiceClient {
	return webServiceSystemClient.Get().(*WebServiceClient)
}

func (this *WebServiceClient) callRequest(
	ctx comcontext.Context,
	request *resty.Request, uri string,
) (_ *resty.Response, err error) {
	var (
		logger    = comlogging.GetLogger()
		startTime = time.Now()
	)
	if err = DumpJsonRequestHMAC(this.secretBytes, request); err != nil {
		return
	}
	response, err := request.Post(uri)

	defer func() {
		latency := time.Now().Sub(startTime)
		logEntry := logger.
			GenHttpOutboundEnty(
				request.Method, latency, request.URL,
				request.Body, nil, response.StatusCode(), "", response.String(),
			).
			WithContext(ctx)

		if err != nil || !response.IsSuccess() {
			logEntry.Error("web service request failed")
		} else {
			var repsonseData meta.O
			if err = comutils.JsonDecode(response.String(), &repsonseData); err != nil {
				logEntry.WithError(err).Error("web service cannot parse json response")
			} else if code, ok := repsonseData["code"]; code != "success" || !ok {
				err = utils.IssueErrorf("web service fail response code `%v`", code)
				logEntry.Error("web service response failed")
			} else {
				logEntry.Info("web service request")
			}
		}
	}()

	return response, err
}

func (this *WebServiceClient) SendEmailProfitWithdraw(ctx comcontext.Context, torqueTxnID uint64) error {
	request := this.httpClient.R().SetBody(meta.O{
		"withdraw_id": torqueTxnID,
	})
	_, err := this.callRequest(ctx, request, "/api/s2s/email-torque-confirm/")
	if err != nil {
		return err
	}

	return nil
}
