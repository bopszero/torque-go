package blockchainmod

import (
	"context"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"golang.org/x/time/rate"
)

var (
	etherscanClient      *resty.Client
	etherscanRateLimiter = rate.NewLimiter(rate.Every(220*time.Millisecond), 1)
)

type EtherscanEthereumClient struct {
	httpClient  *resty.Client
	rateLimiter *rate.Limiter
	apiKey      string
}

func NewEthereumEtherscanClientWithSystemKey() (*EtherscanEthereumClient, error) {
	if config.BlockchainUseTestnet {
		return NewEthereumEtherscanTestRopstenSystemClient()
	} else {
		return NewEthereumEtherscanMainnetSystemClient()
	}
}

func NewEthereumEtherscanMainnetSystemClient() (*EtherscanEthereumClient, error) {
	apiKey := viper.GetString(config.KeyApiEtherscanKey)
	if apiKey == "" {
		return nil, utils.IssueErrorf("missing API key for %s", HostEtherscanProduction)
	}

	return NewEthereumEtherscanMainnetClient(apiKey), nil
}

func NewEthereumEtherscanTestRopstenSystemClient() (*EtherscanEthereumClient, error) {
	apiKey := viper.GetString(config.KeyApiEtherscanKey)
	if apiKey == "" {
		return nil, utils.IssueErrorf("missing API key for %s", HostEtherscanProduction)
	}

	return NewEthereumEtherscanTestRopstenClient(apiKey), nil
}

func NewEthereumEtherscanMainnetClient(apiKey string) *EtherscanEthereumClient {
	if etherscanClient == nil {
		client := utils.NewRestyClient(ClientRequestTimeout)
		client.SetHostURL(HostEtherscanProduction)

		etherscanClient = client
	}

	return &EtherscanEthereumClient{
		apiKey:     apiKey,
		httpClient: etherscanClient,
	}
}

func NewEthereumEtherscanTestRopstenClient(apiKey string) *EtherscanEthereumClient {
	if etherscanClient == nil {
		client := utils.NewRestyClient(ClientRequestTimeout)
		client.SetHostURL(HostEtherscanSandboxRopsten)

		etherscanClient = client
		etherscanRateLimiter.SetLimit(rate.Every(1000 * time.Millisecond))
	}

	return &EtherscanEthereumClient{
		apiKey:     apiKey,
		httpClient: etherscanClient,
	}
}

func (this *EtherscanEthereumClient) genRequest() *resty.Request {
	comutils.PanicOnError(
		etherscanRateLimiter.Wait(context.Background()),
	)
	return this.httpClient.R().SetQueryParam("apikey", this.apiKey)
}

func (this *EtherscanEthereumClient) callRequestGet(request *resty.Request) (string, error) {
	response, err := request.Get("/api")
	if err != nil {
		return "", utils.WrapError(err)
	}
	return this.getResponseBody(response)
}

func (this *EtherscanEthereumClient) getResponseBody(response *resty.Response) (string, error) {
	if !response.IsSuccess() {
		return "", utils.IssueErrorf(
			"api on %s has a fail response | uri=%v,query=%v,status_code=%v,body=%v,response=%v",
			this.httpClient.HostURL, response.Request.RawRequest.URL.Path, response.Request.QueryParam,
			response.StatusCode(), comutils.JsonEncodeF(response.Request.Body), response.String(),
		)
	}

	return response.String(), nil
}

func (this *EtherscanEthereumClient) GetBlock(blockHash string) (_ Block, err error) {
	return nil, utils.IssueErrorf("%s doesn't support get block by hash", this.httpClient.HostURL)
}

type etherscanGetBlockByHeightResponse struct {
	JsonRPC string                 `json:"jsonrpc"`
	Error   meta.O                 `json:"error"`
	Block   EtherscanEthereumBlock `json:"result"`
}

func (this *EtherscanEthereumClient) GetBlockByHeight(height uint64) (_ Block, err error) {
	request := this.genRequest().SetQueryParams(map[string]string{
		"module": EtherscanModuleProxy,
		"action": EtherscanActionGetBlockByHeight,

		"boolean": "true", // Get detail of transactions
		"tag":     uint64ToHex0x(height),
	})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetBlockByHeightResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}
	if responseModel.Block.Hash == "" {
		return nil, utils.WrapError(constants.ErrorDataNotFound)
	}

	for i := range responseModel.Block.Transactions {
		responseModel.Block.Transactions[i].SetClient(this)
	}

	return &responseModel.Block, nil
}

type getLatestBlockResponse struct {
	JsonRPC        string `json:"jsonrpc"`
	Error          meta.O `json:"error"`
	BlockHeightHex string `json:"result"`
}

func (this *EtherscanEthereumClient) GetLatestBlockHeight() (_ uint64, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleProxy,
			"action": EtherscanActionGetBlockHeight,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel getLatestBlockResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}

	blockHeight, err := hex0xToUint64(responseModel.BlockHeightHex)
	if err != nil {
		return
	}

	return blockHeight, nil
}

func (this *EtherscanEthereumClient) GetLatestBlock() (_ Block, err error) {
	blockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}

	return this.GetBlockByHeight(blockHeight)
}

type etherscanGetBalanceResponse struct {
	Status     string          `json:"status"`
	BalanceWei decimal.Decimal `json:"result"`
}

func (this *EtherscanEthereumClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleAccount,
			"action": EtherscanActionGetBalance,

			"address": address,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetBalanceResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != EtherscanStatusSuccess {
		err = utils.IssueErrorf(
			"api to %s failed | data=%v",
			this.httpClient.HostURL, request.QueryParam)
		return
	}

	balanceEth := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    responseModel.BalanceWei,
		},
		constants.CurrencyEthereum,
	)
	return balanceEth.Value, nil
}

type etherscanGetTransactionResponse struct {
	JsonRPC string                          `json:"jsonrpc"`
	Error   meta.O                          `json:"error"`
	Txn     EtherscanEthereumHexTransaction `json:"result"`
}

func (this *EtherscanEthereumClient) GetTxn(hash string) (_ Transaction, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleProxy,
			"action": EtherscanActionGetTxn,

			"txhash": hash,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetTransactionResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}

	responseModel.Txn.SetClient(this)
	return &responseModel.Txn, nil
}

func (this *EtherscanEthereumClient) GetTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
	return this.GetNormalTxns(address, paging)
}

type etherscanListTransactionsResponse struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Txns    []EtherscanEthereumTransaction `json:"result"`
}

func (this *EtherscanEthereumClient) GetNormalTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleAccount,
			"action": EtherscanActionListTxns,

			"address": address,
			"sort":    "desc",
		})
	if paging.Limit > 0 {
		page := (paging.Offset / paging.Limit) + 1
		request.SetQueryParam("offset", comutils.Stringify(paging.Limit))
		request.SetQueryParam("page", comutils.Stringify(page))
	}
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanListTransactionsResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != EtherscanStatusSuccess {
		// We accept the empty transcation list
		if responseModel.Message != "No transactions found" {
			err = utils.IssueErrorf(
				"api to %s failed | data=%v",
				this.httpClient.HostURL, request.QueryParam)
			return
		}
	}

	txns := make([]Transaction, len(responseModel.Txns))
	for i := range responseModel.Txns {
		txns[i] = &responseModel.Txns[i]
		txns[i].SetOwnerAddress(address)
	}

	return txns, nil
}

type etherscanListIntenralTransactionsResponse struct {
	Status  string                                 `json:"status"`
	Message string                                 `json:"message"`
	Txns    []EtherscanEthereumInternalTransaction `json:"result"`
}

type EtherscanEthereumInternalTransaction struct {
	baseTransaction

	Hash        string `json:"hash"`
	BlockHeight uint64 `json:"blockNumber,string"`

	FromAddress string          `json:"from"`
	ToAddress   string          `json:"to"`
	AmountWei   decimal.Decimal `json:"value"`
	InputHex    string          `json:"input"`
	Time        int64           `json:"timeStamp,string"`
	HasError    int8            `json:"isError,string"`

	GasLimit uint32 `json:"gas,string"`
	GasUsed  uint32 `json:"gasUsed,string"`
}

func (this *EtherscanEthereumClient) GetBlockInternalTxns(
	blockHeight uint64, paging meta.Paging,
) (_ []EtherscanEthereumInternalTransaction, err error) {
	request := this.genRequest().SetQueryParams(map[string]string{
		"module": EtherscanModuleAccount,
		"action": EtherscanActionListInternalTxns,

		"startblock": comutils.Stringify(blockHeight),
		"endblock":   comutils.Stringify(blockHeight),
		"sort":       "asc",
	})
	if paging.Limit > 0 {
		page := (paging.Offset / paging.Limit) + 1
		request.SetQueryParam("offset", comutils.Stringify(paging.Limit))
		request.SetQueryParam("page", comutils.Stringify(page))
	}
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanListIntenralTransactionsResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != EtherscanStatusSuccess {
		// We accept the empty transcation list
		if responseModel.Message != "No transactions found" {
			err = utils.IssueErrorf(
				"api to %s failed | data=%v",
				this.httpClient.HostURL, request.QueryParam)
			return
		}
	}

	return responseModel.Txns, nil
}

type etherscanListTokenTransactionsResponse struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Txns    []EtherscanEthereumTokenTransaction `json:"result"`
}

func (this *EtherscanEthereumClient) GetRC20TokenTxns(
	address string, contractAddress string, paging meta.Paging,
) (_ []Transaction, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleAccount,
			"action": EtherscanActionListErc20TokenTxns,

			"contractaddress": contractAddress,
			"address":         address,
			"sort":            "desc",
		})
	if paging.Limit > 0 {
		page := (paging.Offset / paging.Limit) + 1
		request.SetQueryParam("offset", comutils.Stringify(paging.Limit))
		request.SetQueryParam("page", comutils.Stringify(page))
	}
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanListTokenTransactionsResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Status != EtherscanStatusSuccess {
		// We accept the empty transcation list
		if responseModel.Message != "No transactions found" {
			err = utils.IssueErrorf(
				"api to %s failed | data=%v",
				this.httpClient.HostURL, request.QueryParam)
			return
		}
	}

	txns := make([]Transaction, len(responseModel.Txns))
	for i := range responseModel.Txns {
		txns[i] = &responseModel.Txns[i]
		txns[i].SetOwnerAddress(address)
	}

	return txns, nil
}

type etherscanGetTxnReceiptResponse struct {
	JsonRPC    string                      `json:"jsonrpc"`
	Error      meta.O                      `json:"error"`
	TxnReceipt etherscanEthereumTxnReceipt `json:"result"`
}

func (this *EtherscanEthereumClient) GetTxnReceipt(hash string) (
	_ *etherscanEthereumTxnReceipt, err error,
) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleProxy,
			"action": EtherscanActionGetTxnReceipt,

			"txhash": hash,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetTxnReceiptResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}

	return &responseModel.TxnReceipt, nil
}

type etherscanGetTxnsCountResponse struct {
	JsonRPC  string `json:"jsonrpc"`
	Error    meta.O `json:"error"`
	CountHex string `json:"result"`
}

func (this *EtherscanEthereumClient) GetNextNonce(address string) (_ Nonce, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleProxy,
			"action": EtherscanActionGetTxnsCount,

			"address": address,
			"tag":     "latest", // TODO: Check if need to exlude pending
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetTxnsCountResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}

	nonce, err := hex0xToUint64(responseModel.CountHex)
	if err != nil {
		return
	}
	return NewNumberNonce(nonce), nil
}

func (this *EtherscanEthereumClient) GetUtxOutputs(address string, minAmount decimal.Decimal) (
	_ []UnspentTxnOutput, err error,
) {
	return nil, utils.IssueErrorf("Ethereum doesn't have UTXO concept")
}

type etherscanGetCodeResponse struct {
	JsonRPC string `json:"jsonrpc"`
	Error   meta.O `json:"error"`
	CodeHex string `json:"result"`
}

func (this *EtherscanEthereumClient) GetCode(address string) (hexVal string, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleProxy,
			"action": EtherscanActionGetCode,

			"address": address,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetCodeResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	if responseModel.Error != nil {
		err = utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.QueryParam, responseModel.Error)
		return
	}

	return responseModel.CodeHex, nil
}

type etherscanPushTxnRawResponse struct {
	JsonRPC string `json:"jsonrpc"`
	Error   meta.O `json:"error"`
}

func (this *EtherscanEthereumClient) PushTxnRaw(data []byte) error {
	body := url.Values{}
	body.Add("module", EtherscanModuleProxy)
	body.Add("action", EtherscanActionPushTxn)
	body.Add("hex", comutils.HexEncode(data))

	request := this.genRequest().
		SetBody(body.Encode()).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationForm)
	response, err := request.Post("/api")
	if err != nil {
		return utils.WrapError(err)
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return err
	}
	var responseModel etherscanPushTxnRawResponse
	if err := comutils.JsonDecode(responseBody, &responseModel); err != nil {
		return utils.WrapError(err)
	}
	if responseModel.Error != nil {
		return utils.IssueErrorf(
			"rpc api to %s failed | data=%v,error=%v",
			this.httpClient.HostURL,
			request.Body, responseModel.Error)
	}

	return nil
}
