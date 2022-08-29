package blockchainmod

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreEthereumTokenClient struct {
	*BitcoreEthereumClient
	tokenMeta TokenMeta
}

func NewEthereumTokenBitcoreClient(tokenMeta TokenMeta) (*BitcoreEthereumTokenClient, error) {
	if config.BlockchainUseTestnet {
		return nil, utils.IssueErrorf("bitcore testnet hasn't been supported yet")
	} else {
		return NewEthereumTokenBitcoreMainnetClient(tokenMeta)
	}
}

func NewEthereumTokenBitcoreMainnetClient(tokenMeta TokenMeta) (*BitcoreEthereumTokenClient, error) {
	ethClient, err := NewEthereumBitcoreMainnetClient()
	if err != nil {
		return nil, err
	}
	client := BitcoreEthereumTokenClient{
		BitcoreEthereumClient: ethClient,
		tokenMeta:             tokenMeta,
	}
	return &client, nil
}

func (this *BitcoreEthereumTokenClient) GetBlock(blockHash string) (Block, error) {
	return this.baseGetBlock(this, blockHash)
}

func (this *BitcoreEthereumTokenClient) GetBlockByHeight(height uint64) (_ Block, err error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BitcoreEthereumTokenClient) GetBlockTxnsByHeight(height uint64) (
	txns []Transaction, err error,
) {
	uri := this.genBaseURI() + "/tx/"
	request := this.genRequest().SetQueryParams(map[string]string{
		"tokenAddress": this.tokenMeta.Address,
		"blockHeight":  comutils.Stringify(height),
		"limit":        "5000",
	})
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}
	var tokenTxns []BitcoreEthereumTokenTransaction
	err = comutils.JsonDecode(responseBody, &tokenTxns)
	if err != nil {
		return
	}

	txns = make([]Transaction, len(tokenTxns))
	for i := range tokenTxns {
		txn := &tokenTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}
	return txns, nil
}

type bitcoreEthereumTokenGetBlockTipResponse struct {
	Currency meta.Currency `json:"chain"`
	Network  string        `json:"network"`
	Hash     string        `json:"hash"`
	Height   uint64        `json:"height"`
}

func (this *BitcoreEthereumTokenClient) GetLatestBlock() (_ Block, err error) {
	uri := this.genBaseURI() + "/block/tip/"
	responseBody, err := this.callRequestGet(this.genRequest(), uri)
	if err != nil {
		return
	}

	var responseModel bitcoreEthereumTokenGetBlockTipResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}

	return this.GetBlockByHeight(responseModel.Height)
}

type bitcoreEthereumGetTokenBalanceResponse struct {
	Token bitcoreEthereumGetTokenBalanceTokenInfo `json:"token"`
}

type bitcoreEthereumGetTokenBalanceTokenInfo struct {
	Balance       decimal.Decimal `json:"balance"`
	DecimalPlaces int32           `json:"decimals"`
	Symbol        meta.Currency   `json:"symbol"`
	Address       string          `json:"address"`
}

func (this *BitcoreEthereumTokenClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/balance/", address)
	request := this.genRequest().
		SetQueryParam("tokenAddress", this.tokenMeta.Address)
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseModel bitcoreEthereumGetTokenBalanceResponse
	if err = comutils.JsonDecode(responseBody, &responseModel); err != nil {
		return
	}

	tokenBalance := responseModel.Token.Balance.Shift(-responseModel.Token.DecimalPlaces)
	return tokenBalance, nil
}

func (this *BitcoreEthereumTokenClient) GetTxn(hash string) (_ Transaction, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/tx/%s/", hash)
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

	var txn BitcoreEthereumTokenTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	return &txn, nil
}

func (this *BitcoreEthereumTokenClient) GetTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
	return this.GetRC20TokenTxns(address, this.tokenMeta.Address, paging)
}
