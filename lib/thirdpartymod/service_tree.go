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
	"gitlab.com/snap-clickstaff/torque-go/lib/affiliate"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

const (
	TreeApiGetNodeDown = "/v1/basic/get_node_down/"
)

var treeServiceSystemClient = comtypes.NewSingleton(func() interface{} {
	return NewTreeServiceClient()
})

type TreeServiceClient struct {
	httpClient *resty.Client
}

func NewTreeServiceClient() *TreeServiceClient {
	httpClient := utils.NewRestyClient(5 * time.Second)
	httpClient.SetHostURL(viper.GetString(config.KeyServiceTreeBaseURL))

	return &TreeServiceClient{
		httpClient: httpClient,
	}
}

func GetTreeServiceSystemClient() *TreeServiceClient {
	return treeServiceSystemClient.Get().(*TreeServiceClient)
}

func (this *TreeServiceClient) callRequest(
	ctx comcontext.Context,
	request *resty.Request, uri string,
) (_ *resty.Response, err error) {
	var (
		logger    = comlogging.GetLogger()
		startTime = time.Now()
	)
	response, err := request.Post(uri)

	defer func() {
		latency := time.Now().Sub(startTime)
		logEntry := logger.
			GenHttpOutboundEnty(
				request.Method, latency, request.URL,
				request.Body, nil, response.StatusCode(), "", response.String(),
			).
			WithContext(ctx).
			WithField("service", "tree")

		if err != nil || !response.IsSuccess() {
			logEntry.Error("tree service request failed")
		} else {
			var repsonseData meta.O
			if err = comutils.JsonDecode(response.String(), &repsonseData); err != nil {
				logEntry.WithError(err).Error("tree service cannot parse json repsonse")
			} else if code, ok := repsonseData["code"]; code != "success" || !ok {
				err = utils.IssueErrorf("tree service fail response code `%v`", code)
				logEntry.Error("tree service response failed")
			} else {
				logEntry.Info("tree service request")
			}
		}
	}()

	return response, err
}

func (this *TreeServiceClient) GetNodeDown(
	ctx comcontext.Context,
	uid meta.UID, options affiliate.ScanOptions,
) (info affiliate.ScanNodeInfo, err error) {
	request := this.httpClient.R().SetBody(meta.O{
		"user_id": uid,
		"options": options,
	})
	response, err := this.callRequest(ctx, request, TreeApiGetNodeDown)
	if err != nil {
		return
	}

	var responseModel struct {
		Data struct {
			Node affiliate.ScanNodeInfo `json:"node"`
		} `json:"data"`
	}
	if err = comutils.JsonDecode(response.String(), &responseModel); err != nil {
		return
	}
	return responseModel.Data.Node, nil
}
