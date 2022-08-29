package blockchainmod

import (
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type EtherscanEthereumTokenClient struct {
	*EtherscanEthereumClient
	tokenMeta TokenMeta
}

func NewEthereumTokenEtherscanClientWithSystemKey(tokenMeta TokenMeta) (
	*EtherscanEthereumTokenClient, error,
) {
	if config.BlockchainUseTestnet {
		return NewEthereumTokenEtherscanTestRopstenSystemClient(tokenMeta)
	} else {
		return NewEthereumTokenEtherscanMainnetSystemClient(tokenMeta)
	}
}

func NewEthereumTokenEtherscanMainnetSystemClient(tokenMeta TokenMeta) (
	*EtherscanEthereumTokenClient, error,
) {
	apiKey := viper.GetString(config.KeyApiEtherscanKey)
	if apiKey == "" {
		return nil, utils.IssueErrorf("missing API key for %s", HostEtherscanProduction)
	}

	return NewEthereumTokenEtherscanMainnetClient(tokenMeta, apiKey), nil
}

func NewEthereumTokenEtherscanTestRopstenSystemClient(tokenMeta TokenMeta) (
	*EtherscanEthereumTokenClient, error,
) {
	apiKey := viper.GetString(config.KeyApiEtherscanKey)
	if apiKey == "" {
		return nil, utils.IssueErrorf("missing API key for %s", HostEtherscanProduction)
	}

	return NewEthereumTokenEtherscanTestnetClient(tokenMeta, apiKey), nil
}

func NewEthereumTokenEtherscanMainnetClient(tokenMeta TokenMeta, apiKey string) *EtherscanEthereumTokenClient {
	client := NewEthereumEtherscanMainnetClient(apiKey)

	return &EtherscanEthereumTokenClient{
		EtherscanEthereumClient: client,
		tokenMeta:               tokenMeta,
	}
}

func NewEthereumTokenEtherscanTestnetClient(tokenMeta TokenMeta, apiKey string) *EtherscanEthereumTokenClient {
	client := NewEthereumEtherscanTestRopstenClient(apiKey)

	return &EtherscanEthereumTokenClient{
		EtherscanEthereumClient: client,
		tokenMeta:               tokenMeta,
	}
}

func (this *EtherscanEthereumTokenClient) GetBlockByHeight(height uint64) (_ Block, err error) {
	request := this.genRequest().SetQueryParams(map[string]string{
		"module": EtherscanModuleProxy,
		"action": EtherscanActionGetBlockByHeight,

		"boolean": "true",
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
		err = utils.WrapError(constants.ErrorDataNotFound)
		return
	}

	var tokenBlock EtherscanEthereumTokenBlock
	if err = copier.Copy(&tokenBlock, &responseModel.Block); err != nil {
		err = utils.WrapError(err)
		return
	}
	blockTxns, err := this.GetBlockTxnsByHeight(tokenBlock.GetHeight())
	if err != nil {
		return
	}
	tokenBlock.setTxns(blockTxns)

	return &tokenBlock, nil
}

func (this *EtherscanEthereumTokenClient) GetBlockTxnsByHeight(height uint64) (
	txns []EtherscanEthereumTokenTransaction, err error,
) {
	heightStr := comutils.Stringify(height)
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleAccount,
			"action": EtherscanActionListErc20TokenTxns,

			"contractaddress": this.tokenMeta.Address,
			"startblock":      heightStr,
			"endblock":        heightStr,

			"offset": "10000", // All txns
		})
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

	return responseModel.Txns, nil
}

func (this *EtherscanEthereumTokenClient) GetLatestBlock() (_ Block, err error) {
	blockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}

	return this.GetBlockByHeight(blockHeight)
}

type etherscanGetTokenBalanceResponse struct {
	Status string          `json:"status"`
	Amount decimal.Decimal `json:"result"`
}

func (this *EtherscanEthereumTokenClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	request := this.genRequest().
		SetQueryParams(map[string]string{
			"module": EtherscanModuleAccount,
			"action": EtherscanActionGetTokenBalance,

			"contractaddress": this.tokenMeta.Address,
			"address":         address,
		})
	responseBody, err := this.callRequestGet(request)
	if err != nil {
		return
	}

	var responseModel etherscanGetTokenBalanceResponse
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

	tokenAmount := responseModel.Amount.Shift(-int32(this.tokenMeta.DecimalPlaces))
	return tokenAmount, nil
}

func (this *EtherscanEthereumTokenClient) GetTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
	return this.GetRC20TokenTxns(address, this.tokenMeta.Address, paging)
}
