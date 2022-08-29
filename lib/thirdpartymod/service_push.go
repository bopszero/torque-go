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

const (
	ServicePushActionActivity = "activity"

	ServicePushActionDestinationTradingTxn = "TradingWalletTransaction"
	ServicePushActionDestinationWalletTxn  = "PersonalWalletTransaction"
)

var (
	pushServiceSystemClient = comtypes.NewSingleton(func() interface{} {
		return NewPushServiceClient(
			viper.GetString(config.KeyServicePushMacSecret),
		)
	})
)

type PushServiceClient struct {
	httpClient  *resty.Client
	secretBytes []byte
}

type PushServiceMessageData struct {
	Title   string `json:"title"`
	Message string `json:"message"`

	Action            string      `json:"action"`
	ActionDestination string      `json:"action_destination"`
	ActionData        interface{} `json:"action_id"`
}

func NewPushServiceClient(secret string) *PushServiceClient {
	httpClient := utils.NewRestyClient(5 * time.Second)
	httpClient.SetHostURL(viper.GetString(config.KeyServicePushBaseURL))

	return &PushServiceClient{
		httpClient:  httpClient,
		secretBytes: []byte(secret),
	}
}

func GetPushServiceSystemClient() *PushServiceClient {
	return pushServiceSystemClient.Get().(*PushServiceClient)
}

func (this *PushServiceClient) callRequest(
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
			logEntry.Error("push service request failed")
		} else {
			var repsonseData meta.O
			if err = comutils.JsonDecode(response.String(), &repsonseData); err != nil {
				logEntry.WithError(err).Error("push service cannot parse json response")
			} else if code, ok := repsonseData["code"]; code != "success" || !ok {
				err = utils.IssueErrorf("push service fail response code `%v`", code)
				logEntry.Error("push service response failed")
			} else {
				logEntry.Info("push service request")
			}
		}
	}()

	return response, err
}

type PushServicePushRequest struct {
	UID  meta.UID               `json:"user_id"`
	Data PushServiceMessageData `json:"data"`
}

func (this *PushServiceClient) Push(ctx comcontext.Context, uid meta.UID, data PushServiceMessageData) error {
	request := this.httpClient.R().SetBody(
		PushServicePushRequest{
			UID:  uid,
			Data: data,
		},
	)
	_, err := this.callRequest(ctx, request, "/v1/push/")
	if err != nil {
		return err
	}

	return nil
}
