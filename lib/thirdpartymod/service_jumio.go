package thirdpartymod

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	jumioSystemClient *JumioServiceClient
)

type JumioServiceClient struct {
	httpClient *resty.Client
}

type JumioCredentials struct {
	ApiToken  string
	ApiSecret string
	Scheme    string
	Host      string
}

func NewJumioServiceClient() (*JumioServiceClient, error) {
	jumioCredentials, err := GetJumioCredentials()
	if err != nil {
		return nil, err
	}
	jumioHost := fmt.Sprintf("%v://%v", jumioCredentials.Scheme, jumioCredentials.Host)
	httpClient := utils.NewRestyClient(5 * time.Second)
	httpClient.SetHeader(echo.HeaderAccept, constants.ContentTypeJSON)
	httpClient.SetHeader(echo.HeaderContentType, constants.ContentTypeJSON)
	httpClient.SetHostURL(jumioHost)
	httpClient.SetBasicAuth(jumioCredentials.ApiToken, jumioCredentials.ApiSecret)
	jumioServiceClient := &JumioServiceClient{
		httpClient,
	}
	return jumioServiceClient, nil
}

func GetJumioCredentials() (*JumioCredentials, error) {
	jumioURL, err := url.Parse(viper.GetString(config.KeyJumioDSN))
	if err != nil {
		return nil, utils.WrapError(err)
	}
	apiSecret, passwordSet := jumioURL.User.Password()
	if !passwordSet {
		return nil, utils.IssueErrorf("jumio API secret has not setted")
	}
	jumioCredentials := &JumioCredentials{}
	jumioCredentials.ApiToken = jumioURL.User.Username()
	jumioCredentials.ApiSecret = apiSecret
	jumioCredentials.Scheme = jumioURL.Scheme
	jumioCredentials.Host = jumioURL.Host
	return jumioCredentials, nil
}

func GetJumioServiceSystemClient() (*JumioServiceClient, error) {
	if jumioSystemClient == nil {
		newJumioServiceClient, err := NewJumioServiceClient()
		if err != nil {
			return nil, err
		}
		jumioSystemClient = newJumioServiceClient
	}
	return jumioSystemClient, nil
}

func (this *JumioServiceClient) callPostRequest(ctx comcontext.Context, request *resty.Request, uri string) (
	response *resty.Response, err error,
) {
	startTime := time.Now()
	response, err = request.Post(uri)

	defer this.logRequest(ctx, startTime, request, response, err)
	return response, err
}

func (this *JumioServiceClient) callGetRequest(ctx comcontext.Context, request *resty.Request, uri string) (
	response *resty.Response, err error,
) {
	startTime := time.Now()
	response, err = request.Get(uri)

	defer this.logRequest(ctx, startTime, request, response, err)
	return response, err
}

func (this *JumioServiceClient) logRequest(
	ctx comcontext.Context,
	startTime time.Time, request *resty.Request, response *resty.Response, errResponse error,
) {
	logger := comlogging.GetLogger()
	latency := time.Now().Sub(startTime)
	logEntry := logger.
		GenHttpOutboundEnty(
			request.Method, latency, request.URL,
			request.Body, nil, response.StatusCode(), "", response.String(),
		).
		WithContext(ctx)
	if errResponse != nil || !response.IsSuccess() {
		logEntry.Error("request to Jumio failed")
	} else {
		var responseData meta.O
		if err := comutils.JsonDecode(response.String(), &responseData); err != nil {
			logEntry.WithError(err).Error("jumio service cannot parse json response")
		} else {
			logEntry.Info("jumio service request")
		}
	}
}

type InitTransactionRequest struct {
	CustomerInternalReference string `json:"customerInternalReference"`
	UserReference             string `json:"userReference"`
	Locale                    string `json:"locale"`
}

func (this *JumioServiceClient) InitTransaction(
	ctx comcontext.Context,
	customerInternalReference string, userReference string, locale string,
) (*resty.Response, error) {
	requestBody := this.httpClient.R().SetBody(
		InitTransactionRequest{
			customerInternalReference,
			userReference,
			locale,
		},
	)
	response, err := this.callPostRequest(ctx, requestBody, "/api/v4/initiate/")
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (this *JumioServiceClient) GetScanDetailsResp(
	ctx comcontext.Context, scanReference string,
) (*resty.Response, error) {
	uri := fmt.Sprintf("/api/netverify/v2/scans/%v/data", scanReference)
	response, err := this.callGetRequest(ctx, this.httpClient.R(), uri)
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get scan details from Jumio failed | status_code=%v",
			response.StatusCode())
		return nil, err
	}
	return response, nil
}

func (this *JumioServiceClient) GetVerificationDataResp(
	ctx comcontext.Context, scanReference string,
) (*resty.Response, error) {
	uri := fmt.Sprintf("/api/netverify/v2/scans/%v/data/verification", scanReference)
	response, err := this.callGetRequest(ctx, this.httpClient.R(), uri)
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		err = utils.IssueErrorf(
			"get verification data from Jumio failed | status_code=%v",
			response.StatusCode())
		return nil, err
	}
	return response, nil
}
