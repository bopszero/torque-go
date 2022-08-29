package blockchainmod

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreTronTokenClient struct {
	*BitcoreTronClient
	tokenMeta TokenMeta
}

func NewTronTokenBitcoreClient(tokenMeta TokenMeta) (*BitcoreTronTokenClient, error) {
	if config.BlockchainUseTestnet {
		return NewTronTokenBitcoreTestShastaClient(tokenMeta), nil
	} else {
		return NewTronTokenBitcoreMainnetClient(tokenMeta), nil
	}
}

func NewTronTokenBitcoreTestShastaClient(tokenMeta TokenMeta) *BitcoreTronTokenClient {
	return &BitcoreTronTokenClient{
		BitcoreTronClient: NewTronBitcoreTestShastaClient(),
		tokenMeta:         tokenMeta,
	}
}

func NewTronTokenBitcoreMainnetClient(tokenMeta TokenMeta) *BitcoreTronTokenClient {
	return &BitcoreTronTokenClient{
		BitcoreTronClient: NewTronBitcoreMainnetClient(),
		tokenMeta:         tokenMeta,
	}
}

func (this *BitcoreTronTokenClient) GetBlockTxnsByHeight(height uint64) (
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
	var tronTxns []BitcoreTronTokenTransaction
	err = comutils.JsonDecode(responseBody, &tronTxns)
	if err != nil {
		return
	}

	txns = make([]Transaction, 0, len(tronTxns))
	for i := range tronTxns {
		txn := &tronTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		if err = txn.initTokenMeta(this.tokenMeta); err != nil {
			if !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
				return
			}
			continue
		}
		txns = append(txns, txn)
	}
	return txns, nil
}

type bitcoreTronTokenGetBalanceResponse struct {
	Token bitcoreTronTokenGetBalance `json:"token"`
}

type bitcoreTronTokenGetBalance struct {
	Balance       decimal.Decimal `json:"balance"`
	Name          string          `json:"name"`
	Currency      meta.Currency   `json:"symbol"`
	DecimalPlaces int32           `json:"decimals"`
	Address       string          `json:"address"`
}

func (this *BitcoreTronTokenClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	var (
		uri     = this.genBaseURI() + fmt.Sprintf("/address/%s/balance/", address)
		request = this.genRequest().
			SetQueryParam("tokenAddress", this.tokenMeta.Address)
	)
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseModel bitcoreTronTokenGetBalanceResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	balanceBase := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.currencyUnit,
			Value:    responseModel.Token.Balance,
		},
		this.currencyBase,
	)
	return balanceBase.Value, nil
}

func (this *BitcoreTronTokenClient) GetTxn(hash string) (_ Transaction, err error) {
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

	var txn BitcoreTronTokenTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	if err = txn.initTokenMeta(this.tokenMeta); err != nil {
		return
	}

	return &txn, nil
}

func (this *BitcoreTronTokenClient) GetTxns(address string, paging meta.Paging) (
	[]Transaction, error,
) {
	return this.GetRC20TokenTxns(address, this.tokenMeta, paging)
}
