package blockchainmod

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"golang.org/x/time/rate"
)

var (
	soChainRateLimiter = rate.NewLimiter(rate.Every(220*time.Millisecond), 1)
	soChainClient      = comtypes.NewSingleton(func() interface{} {
		client := utils.NewRestyClient(ClientRequestTimeout)
		client.SetHostURL(HostSoChainProduction)

		return client
	})
)

const (
	soChainNetworkBitcoin  = "btc"
	soChainNetworkLitecoin = "ltc"

	soChainNetworkTestSuffix = "test"
)

func getSoChainClient() *resty.Client {
	return soChainClient.Get().(*resty.Client)
}

type SoChainUtxoLikeClient struct {
	httpClient *resty.Client

	currency    meta.Currency
	networkCode string
}

func NewBitcoinSoChainClient() *SoChainUtxoLikeClient {
	if config.BlockchainUseTestnet {
		return NewBitcoinSoChainTestnetClient()
	} else {
		return NewBitcoinSoChainMainnetClient()
	}
}

func NewBitcoinSoChainMainnetClient() *SoChainUtxoLikeClient {
	return &SoChainUtxoLikeClient{
		currency:    constants.CurrencyBitcoin,
		networkCode: soChainNetworkBitcoin,

		httpClient: getSoChainClient(),
	}
}

func NewBitcoinSoChainTestnetClient() *SoChainUtxoLikeClient {
	return &SoChainUtxoLikeClient{
		currency:    constants.CurrencyBitcoin,
		networkCode: soChainNetworkBitcoin + soChainNetworkTestSuffix,

		httpClient: getSoChainClient(),
	}
}

func (this *SoChainUtxoLikeClient) genBaseURI(action string) string {
	return fmt.Sprintf("/api/v2/%s/%s/", action, this.networkCode)
}

func (this *SoChainUtxoLikeClient) genRequest() *resty.Request {
	comutils.PanicOnError(
		soChainRateLimiter.Wait(context.Background()),
	)
	return this.httpClient.R()
}

func (this *SoChainUtxoLikeClient) getResponseBody(response *resty.Response) (string, error) {
	if !response.IsSuccess() {
		return "", utils.IssueErrorf(
			"api on %s has a fail response | uri=%v,status_code=%v,body=%v,response=%v",
			this.httpClient.HostURL, response.Request.RawRequest.URL.Path,
			response.StatusCode(), comutils.JsonEncodeF(response.Request.Body), response.String(),
		)
	}

	return response.String(), nil
}

type bitcoinSoChainGetBlockResponse struct {
	Status string              `json:"status"`
	Block  SoChainBitcoinBlock `json:"data"`
}

func (this *SoChainUtxoLikeClient) GetBlock(blockHash string) (_ Block, err error) {
	uri := this.genBaseURI("get_block") + fmt.Sprintf("%s/", blockHash)

	response, err := this.genRequest().Get(uri)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return
	}

	var responseModel bitcoinSoChainGetBlockResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	block := responseModel.Block
	block.setClient(this)

	return &block, nil
}

func (this *SoChainUtxoLikeClient) GetBlockByHeight(height uint64) (Block, error) {
	return this.GetBlock(comutils.Stringify(height))
}

type bitcoinSoChainGetNetworkInfoResponse struct {
	Status string                       `json:"status"`
	Data   bitcoinSoChainGetNetworkInfo `json:"data"`
}

type bitcoinSoChainGetNetworkInfo struct {
	Height        uint64          `json:"blocks"`
	Price         decimal.Decimal `json:"price"`
	PriceCurrency meta.Currency   `json:"price_base"`
}

func (this *SoChainUtxoLikeClient) GetLatestBlock() (_ Block, err error) {
	uri := this.genBaseURI("get_info")

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

	var responseModel bitcoinSoChainGetNetworkInfoResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	return this.GetBlockByHeight(responseModel.Data.Height)
}

type bitcoinSoChainBalance struct {
	Network          string          `json:"network"`
	Address          string          `json:"address"`
	ValueConfirmed   decimal.Decimal `json:"confirmed_balance"`
	ValueUnconfirmed decimal.Decimal `json:"unconfirmed_balance"`
}

type bitcoinSoChainGetBalanceResponse struct {
	Status string                `json:"status"`
	Data   bitcoinSoChainBalance `json:"data"`
}

func (this *SoChainUtxoLikeClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	uri := this.genBaseURI("get_address_balance") + fmt.Sprintf("%s/", address)

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

	var responseModel bitcoinSoChainGetBalanceResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	return responseModel.Data.ValueConfirmed, nil
}

type bitcoinSoChainGetTxnResponse struct {
	Status string                    `json:"status"`
	Data   SoChainBitcoinTransaction `json:"data"`
}

func (this *SoChainUtxoLikeClient) GetTxn(hash string) (_ Transaction, err error) {
	uri := this.genBaseURI("tx") + fmt.Sprintf("%s/", hash)

	request := this.genRequest()
	responseBody, err := func() (string, error) {
		response, err := request.Get(uri)
		if err != nil {
			return "", err
		}
		if response.StatusCode() == http.StatusInternalServerError {
			return "", constants.ErrorDataNotFound
		}
		return this.getResponseBody(response)
	}()
	if err != nil {
		return
	}

	var responseModel bitcoinSoChainGetTxnResponse
	if err = comutils.JsonDecode(responseBody, &responseModel); err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	return &responseModel.Data, nil
}

type soChainAddressInfo struct {
	TotalTxnCount uint64                             `json:"total_txs"`
	Txns          []SoChainBitcoinDisplayTransaction `json:"txs"`
}

type soChainGetAdressInfoResponse struct {
	Status string
	Data   soChainAddressInfo `json:"data"`
}

func (this *SoChainUtxoLikeClient) GetTxns(
	address string, paging meta.Paging,
) (txns []Transaction, err error) {
	uri := this.genBaseURI("address") + fmt.Sprintf("%s/", address)

	response, err := this.genRequest().Get(uri)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return
	}

	var responseModel soChainGetAdressInfoResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	var (
		txnLen       = uint(len(responseModel.Data.Txns))
		fromIdx      = comutils.MinUint(paging.Offset, txnLen)
		toIdx        = comutils.MinUint(paging.Offset+paging.Limit, txnLen)
		responseTxns = responseModel.Data.Txns[fromIdx:toIdx]
	)

	txns = make([]Transaction, 0, len(responseTxns))
	for i := range responseTxns {
		txn := &responseTxns[i]
		txn.setCurrency(this.currency)
		txn.SetOwnerAddress(address)
		txns = append(txns, txn)
	}

	return txns, nil
}

func (this *SoChainUtxoLikeClient) GetNextNonce(address string) (Nonce, error) {
	return nil, utils.IssueErrorf("%v doesn't have Nonce concept", this.currency)
}

func (this *SoChainUtxoLikeClient) GetUtxOutputs(address string, minAmount decimal.Decimal) (
	_ []UnspentTxnOutput, err error,
) {
	utxOutputs := make([]UnspentTxnOutput, 0)

	var (
		lastHash     = ""
		filledAmount = minAmount
	)
	for filledAmount.GreaterThan(decimal.Zero) {
		pageUtxOutputs, err := this.getUtxOutputs(address, lastHash)
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
		lastHash = pageUtxOutputs[len(pageUtxOutputs)-1].GetTxnHash()
	}

	return utxOutputs, nil
}

type bitcoinSoChainListUtxoInfo struct {
	Network    string                    `json:"network"`
	Address    string                    `json:"address"`
	UtxOutputs []SoChainBitcoinUtxOutput `json:"txs"`
}

type bitcoinSoChainListUtxoResponse struct {
	Status string                     `json:"status"`
	Data   bitcoinSoChainListUtxoInfo `json:"data"`
}

func (this *SoChainUtxoLikeClient) getUtxOutputs(
	address string, afterHash string,
) (_ []SoChainBitcoinUtxOutput, err error) {
	uri := this.genBaseURI("get_tx_unspent") + fmt.Sprintf("%s/", address)
	if afterHash != "" {
		uri += fmt.Sprintf("%s/", afterHash)
	}
	response, err := this.genRequest().Get(uri)
	if err != nil {
		return
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return
	}

	var responseModel bitcoinSoChainListUtxoResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != SoChainStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | uri=%v",
			this.httpClient.HostURL, uri)
		return
	}

	confirmedOtxOutputs := make([]SoChainBitcoinUtxOutput, 0, len(responseModel.Data.UtxOutputs))
	for _, utxOutput := range responseModel.Data.UtxOutputs {
		if utxOutput.Confirmations < 1 {
			continue
		}

		utxOutput.SetAddress(address)
		confirmedOtxOutputs = append(confirmedOtxOutputs, utxOutput)
	}

	return confirmedOtxOutputs, nil
}

type soChainPushTxnInfo struct {
	TxnHash string `json:"txid"`
}

type soChainPushTxnResponse struct {
	Status string             `json:"status"`
	Data   soChainPushTxnInfo `json:"data"`
}

func (this *SoChainUtxoLikeClient) PushTxnRaw(data []byte) error {
	uri := this.genBaseURI("send_tx")
	body := meta.O{
		"tx_hex": comutils.HexEncode(data),
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
	var responseModel soChainPushTxnResponse

	if err := comutils.JsonDecode(responseBody, &responseModel); err != nil {
		return utils.WrapError(err)
	}
	if responseModel.Data.TxnHash == "" {
		return utils.IssueErrorf(
			"SoChain push Bitcoin transaction failed | hex=%v,response=%v",
			comutils.HexEncode(data), responseBody,
		)
	}

	return nil
}
