package blockchainmod

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var bitcoreClient = comtypes.NewSingleton(func() interface{} {
	client := utils.NewRestyClient(ClientRequestTimeout)
	client.SetBasicAuth(MyCoinGateAuthUsername, MyCoinGateAuthPassword)

	switch config.Env {
	case config.EnvStaging, config.EnvTest, config.EnvDev:
		client.SetHostURL(HostMyCoinGateBitcoreTest)
		break
	default:
		client.SetHostURL(HostMyCoinGateBitcoreProduction)
		break
	}

	return client
})

func getBitcoreClient() *resty.Client {
	return bitcoreClient.Get().(*resty.Client)
}

type (
	iBitcoreBalanceLikeClient interface {
		Client

		GetBlockTxnsByHeight(uint64) ([]Transaction, error)

		genBaseURI() string
		genRequest() *resty.Request
		getResponseBody(*resty.Response) (string, error)
	}

	tBitcoreClient struct {
		httpClient   *resty.Client
		chain        string
		currencyBase meta.Currency
		currencyUnit meta.Currency
	}
)

func newBitcoreClient(
	currencyBase meta.Currency, currencyUnit meta.Currency, chainCode string,
) tBitcoreClient {
	return tBitcoreClient{
		httpClient:   getBitcoreClient(),
		chain:        chainCode,
		currencyBase: currencyBase,
		currencyUnit: currencyUnit,
	}
}

func (this *tBitcoreClient) isMainChain() bool {
	return this.chain == BitcoreChainMainnet
}

func (this *tBitcoreClient) genRequest() *resty.Request {
	return this.httpClient.R()
}

func (this *tBitcoreClient) genBaseURI() string {
	return fmt.Sprintf("/api/%v/%v", this.currencyBase, this.chain)
}

func (this *tBitcoreClient) callRequestGet(request *resty.Request, uri string) (string, error) {
	response, err := request.Get(uri)
	if err != nil {
		return "", utils.WrapError(err)
	}

	return this.getResponseBody(response)
}

func (this *tBitcoreClient) getResponseBody(response *resty.Response) (string, error) {
	if !response.IsSuccess() {
		return "", utils.IssueErrorf(
			"api on %s has a fail response | uri=%v,status_code=%v,body=%v,response=%v",
			this.httpClient.HostURL, response.Request.RawRequest.URL.Path,
			response.StatusCode(), comutils.JsonEncodeF(response.Request.Body), response.String(),
		)
	}

	return response.String(), nil
}

func (*tBitcoreClient) baseGetBlock(that iBitcoreBalanceLikeClient, blockHash string) (_ Block, err error) {
	uri := that.genBaseURI() + fmt.Sprintf("/block/%v/", blockHash)
	response, err := that.genRequest().Get(uri)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	responseBody, err := that.getResponseBody(response)
	if err != nil {
		return
	}

	var block BitcoreBalanceLikeBlock
	err = comutils.JsonDecode(responseBody, &block)
	if err != nil {
		return
	}

	txns, err := that.GetBlockTxnsByHeight(block.GetHeight())
	if err != nil {
		return
	}
	block.setTxns(txns)

	return &block, nil
}

type bitcoreGetBlockTipResponse struct {
	Currency meta.Currency `json:"chain"`
	Network  string        `json:"network"`
	Hash     string        `json:"hash"`
	Height   uint64        `json:"height"`
}

func (this *tBitcoreClient) GetLatestBlockHeight() (_ uint64, err error) {
	uri := this.genBaseURI() + "/block/tip/"
	responseBody, err := this.callRequestGet(this.genRequest(), uri)
	if err != nil {
		return
	}

	var responseModel bitcoreGetBlockTipResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	return responseModel.Height, nil
}

func (this *tBitcoreClient) GetNextNonce(address string) (Nonce, error) {
	return nil, utils.IssueErrorf("%v doesn't have Nonce concept", this.currencyBase)
}

func (this *tBitcoreClient) GetUtxOutputs(address string, minAmount decimal.Decimal) (
	_ []UnspentTxnOutput, err error,
) {
	return nil, utils.IssueErrorf("%v doesn't have UTXO concept", this.currencyBase)
}

type bitcorePushTxnRawResponse struct {
	TxnHash string `json:"txid"`
}

func (this *tBitcoreClient) PushTxnRaw(data []byte) error {
	uri := this.genBaseURI() + "/tx/send/"
	body := meta.O{
		"rawTx": comutils.HexEncode(data),
	}

	request := this.genRequest().SetBody(body)
	response, err := request.Post(uri)
	if err != nil {
		return utils.WrapError(err)
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return err
	}
	var responseModel bitcorePushTxnRawResponse
	if err := comutils.JsonDecode(responseBody, &responseModel); err != nil {
		return utils.WrapError(err)
	}
	if responseModel.TxnHash == "" {
		return utils.IssueErrorf(
			"pushed a %v txn with a empty hash in response",
			this.currencyBase,
		)
	}

	return nil
}
