package blockchainmod

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"golang.org/x/time/rate"
)

var (
	blockchainInfoRateLimiter = rate.NewLimiter(rate.Every(220*time.Millisecond), 1)
	blockchainInfoClientProxy = comtypes.NewSingleton(func() interface{} {
		client := utils.NewRestyClient(30 * time.Second)
		client.SetHostURL(HostBlockchainInfoProduction)
		return client
	})
)

type BlockchainInfoBitcoinClient struct {
	httpClient *resty.Client
	apiKey     string
}

var _ Client = (*BlockchainInfoBitcoinClient)(nil)

func NewBitcoinBlockchainInfoMainnetClientWithSystemKey() (*BlockchainInfoBitcoinClient, error) {
	apiKey := viper.GetString(config.KeyApiBlockchainInfoKey)
	if apiKey == "" {
		return nil, utils.IssueErrorf("missing API key for %s", HostBlockchainInfoProduction)
	}

	return NewBitcoinBlockchainInfoMainnetClient(apiKey), nil
}

func NewBitcoinBlockchainInfoMainnetClient(apiKey string) *BlockchainInfoBitcoinClient {
	return &BlockchainInfoBitcoinClient{
		apiKey:     apiKey,
		httpClient: blockchainInfoClientProxy.Get().(*resty.Client),
	}
}

func (this *BlockchainInfoBitcoinClient) genRequest() *resty.Request {
	comutils.PanicOnError(
		blockchainInfoRateLimiter.Wait(context.Background()),
	)
	return this.httpClient.R().SetQueryParam("api_code", this.apiKey)
}

func (this *BlockchainInfoBitcoinClient) getResponseBody(response *resty.Response) (string, error) {
	if !response.IsSuccess() {
		return "", utils.IssueErrorf(
			"api on %s has a fail response | uri=%v,status_code=%v,body=%v",
			this.httpClient.HostURL, response.Request.RawRequest.URL.Path,
			response.StatusCode(), response.String(),
		)
	}

	return response.String(), nil
}

func (this *BlockchainInfoBitcoinClient) GetBlock(blockHash string) (_ Block, err error) {
	uri := fmt.Sprintf("/rawblock/%s", blockHash)

	response, err := this.genRequest().Get(uri)
	if err != nil {
		return
	}
	if response.StatusCode() == http.StatusInternalServerError &&
		response.String() == "Index: 0, Size: 0" {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return nil, err
	}

	var block BlockchainInfoBitcoinBlock
	err = comutils.JsonDecode(responseBody, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (this *BlockchainInfoBitcoinClient) GetBlockByHeight(height uint64) (_ Block, err error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BlockchainInfoBitcoinClient) GetLatestBlock() (_ Block, err error) {
	request := this.genRequest()
	responseBody, err := func() (string, error) {
		response, err := request.Get("/latestblock")
		if err != nil {
			return "", err
		}
		return this.getResponseBody(response)
	}()

	var responseData meta.O
	err = comutils.JsonDecode(responseBody, &responseData)
	if err != nil {
		return nil, err
	}

	blockHash, ok := responseData["hash"]
	if !ok {
		return nil, utils.IssueErrorf("blockchain.info latest block missing `hash` field")
	}

	return this.GetBlock(blockHash.(string))
}

type blockchainInfoAddressInfo struct {
	Address       string                             `json:"address"`
	TxnCount      uint64                             `json:"n_tx"`
	TotalReceived decimal.Decimal                    `json:"total_received"`
	TotalSent     decimal.Decimal                    `json:"total_sent"`
	Balance       decimal.Decimal                    `json:"final_balance"`
	Txns          []BlockchainInfoBitcoinTransaction `json:"txs"`
}

func (this *BlockchainInfoBitcoinClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	uri := fmt.Sprintf("/rawaddr/%s", address)

	request := this.genRequest().SetQueryParam("limit", "0")
	responseBody, err := func() (string, error) {
		response, err := request.Get(uri)
		if err != nil {
			return "", err
		}
		return this.getResponseBody(response)
	}()
	if err != nil {
		return
	}

	var addressInfo blockchainInfoAddressInfo
	err = comutils.JsonDecode(responseBody, &addressInfo)
	if err != nil {
		return
	}

	return addressInfo.Balance, nil
}

func (this *BlockchainInfoBitcoinClient) GetTxn(hash string) (_ Transaction, err error) {
	uri := fmt.Sprintf("/rawtx/%s", hash)

	request := this.genRequest()
	responseBody, err := func() (string, error) {
		response, err := request.Get(uri)
		if err != nil {
			return "", err
		}
		return this.getResponseBody(response)
	}()
	if err != nil {
		return
	}

	var txn BlockchainInfoBitcoinTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	return &txn, nil
}

func (this *BlockchainInfoBitcoinClient) GetTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
	uri := fmt.Sprintf("/rawaddr/%s", address)

	request := this.genRequest()
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
	}
	if paging.Offset > 0 {
		request.SetQueryParam("offset", comutils.Stringify(paging.Offset))
	}
	responseBody, err := func() (string, error) {
		response, err := request.Get(uri)
		if err != nil {
			return "", err
		}
		return this.getResponseBody(response)
	}()
	if err != nil {
		return
	}

	var addressInfo blockchainInfoAddressInfo
	err = comutils.JsonDecode(responseBody, &addressInfo)
	if err != nil {
		return
	}

	txns := make([]Transaction, len(addressInfo.Txns))
	for i := range addressInfo.Txns {
		txns[i] = &addressInfo.Txns[i]
		txns[i].SetOwnerAddress(address)
	}

	return txns, nil
}

func (this *BlockchainInfoBitcoinClient) GetNextNonce(address string) (Nonce, error) {
	return nil, utils.IssueErrorf("BTC doesn't have Nonce concept")
}

func (this *BlockchainInfoBitcoinClient) GetUtxOutputs(address string, minAmount decimal.Decimal) (
	_ []UnspentTxnOutput, err error,
) {
	utxOutputs := make([]UnspentTxnOutput, 0)

	paging := meta.Paging{Limit: 10}
	filledAmount := minAmount
	for i := 0; filledAmount.GreaterThan(decimal.Zero); i++ {
		paging.SetPage(uint(i))
		pageUtxOutputs, err := this.getUtxOutputs(address, paging)
		if err != nil {
			return nil, err
		}
		if len(pageUtxOutputs) == 0 {
			break
		}

		for i := range pageUtxOutputs {
			utxo := pageUtxOutputs[i]
			utxOutputs = append(utxOutputs, &utxo)
			filledAmount = filledAmount.Sub(utxo.GetAmount())
		}
	}

	return utxOutputs, nil
}

type blockchainInfoGetUtxOutputsResponse struct {
	UnspentOutputs []BlockchainInfoBitcoinUtxOutput `json:"unspent_outputs"`
}

func (this *BlockchainInfoBitcoinClient) getUtxOutputs(
	address string, paging meta.Paging,
) (_ []BlockchainInfoBitcoinUtxOutput, err error) {
	request := this.genRequest().SetQueryParam("active", address)
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
	}
	if paging.Offset > 0 {
		request.SetQueryParam("offset", comutils.Stringify(paging.Offset))
	}
	responseBody, err := func() (string, error) {
		response, err := request.Get("/unspent")
		if err != nil {
			return "", err
		}
		return this.getResponseBody(response)
	}()
	if err != nil {
		return
	}

	var addressInfo blockchainInfoGetUtxOutputsResponse
	err = comutils.JsonDecode(responseBody, &addressInfo)
	if err != nil {
		return
	}

	for i := range addressInfo.UnspentOutputs {
		addressInfo.UnspentOutputs[i].address = address
	}

	return addressInfo.UnspentOutputs, nil
}

func (this *BlockchainInfoBitcoinClient) PushTxnRaw(data []byte) error {
	body := meta.O{
		"hex": comutils.HexEncode(data),
	}
	request := this.genRequest().
		SetBody(body).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response, err := request.Post("/pushtx")
	if err != nil {
		return utils.WrapError(err)
	}
	if _, err = this.getResponseBody(response); err != nil {
		return err
	}

	// TODO: Validate response

	return nil
}
