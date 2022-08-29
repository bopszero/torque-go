package blockchainmod

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreUtxoLikeClient struct {
	tBitcoreClient
}

func NewBitcoreUtxoLikeClient(caseCurrency meta.Currency, chain string) BitcoreUtxoLikeClient {
	return BitcoreUtxoLikeClient{
		tBitcoreClient: newBitcoreClient(caseCurrency, constants.CurrencySubBitcoinSatoshi, chain),
	}
}

func (this *BitcoreUtxoLikeClient) GetBlock(blockHash string) (Block, error) {
	return this.baseGetBlock(this, blockHash)
}

func (this *BitcoreUtxoLikeClient) GetBlockByHeight(height uint64) (_ Block, err error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BitcoreUtxoLikeClient) GetBlockTxnsByHeight(height uint64) (
	txns []Transaction, err error,
) {
	var (
		uri          = this.genBaseURI() + "/tx/"
		allTxns      []BitcoreUtxoLikeTransaction
		pagingOffset = 0
		pagingLimit  = 200
	)
	for {
		var pageTxns []BitcoreUtxoLikeTransaction
		request := this.genRequest().SetQueryParams(map[string]string{
			"blockHeight": comutils.Stringify(height),
			"full":        "true",
			"index":       comutils.Stringify(pagingOffset),
			"limit":       comutils.Stringify(pagingLimit),
		})
		responseBody, reqErr := this.callRequestGet(request, uri)
		if reqErr != nil {
			err = reqErr
			return
		}
		err = comutils.JsonDecode(responseBody, &pageTxns)
		if err != nil {
			return
		}
		allTxns = append(allTxns, pageTxns...)

		pagingOffset += pagingLimit
		if len(pageTxns) < pagingLimit {
			break
		}
	}

	txns = make([]Transaction, len(allTxns))
	for i := range allTxns {
		txn := &allTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}
	return txns, nil
}

func (this *BitcoreUtxoLikeClient) GetLatestBlock() (_ Block, err error) {
	latestBlockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}
	return this.GetBlockByHeight(latestBlockHeight)
}

type bitcoreBitcoinGetBalanceResponse struct {
	Confirmed   decimal.Decimal `json:"confirmed"`
	Unconfirmed decimal.Decimal `json:"unconfirmed"`
	Balance     decimal.Decimal `json:"balance"`
}

func (this *BitcoreUtxoLikeClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/balance/", address)
	responseBody, err := this.callRequestGet(this.genRequest(), uri)
	if err != nil {
		return
	}

	var responseModel bitcoreBitcoinGetBalanceResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}

	balance := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.currencyUnit,
			Value:    responseModel.Confirmed,
		},
		this.currencyBase,
	)
	return balance.Value, nil
}

func (this *BitcoreUtxoLikeClient) GetTxn(hash string) (_ Transaction, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/tx/%s/", hash)

	request := this.genRequest().
		SetQueryParam("full", "true")
	response, err := request.Get(uri)
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

	var txn BitcoreUtxoLikeTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	return &txn, nil
}

func (this *BitcoreUtxoLikeClient) GetTxns(address string, paging meta.Paging) (
	_ []Transaction, err error,
) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/txs/", address)

	request := this.genRequest().
		SetQueryParam("full", "true")
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
		request.SetQueryParam("index", comutils.Stringify(paging.Offset))
	}
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseTxns []BitcoreUtxoLikeTransaction
	err = comutils.JsonDecode(responseBody, &responseTxns)
	if err != nil {
		return
	}

	txns := make([]Transaction, len(responseTxns))
	for i := range responseTxns {
		txn := &responseTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txn.SetOwnerAddress(address)
		txns[i] = txn
	}

	return txns, nil
}

func (this *BitcoreUtxoLikeClient) GetUtxOutputs(address string, minAmount decimal.Decimal) (
	_ []UnspentTxnOutput, err error,
) {
	var (
		utxOutputs   = make([]UnspentTxnOutput, 0)
		paging       = meta.Paging{Limit: 50}
		filledAmount = minAmount
	)
	for i := 0; filledAmount.GreaterThan(decimal.Zero); i++ {
		paging.SetPage(uint(i))
		pageUtxOutputs, err := this.getUtxOutputs(address, paging)
		if err != nil {
			return nil, err
		}

		for i := range pageUtxOutputs {
			utxo := pageUtxOutputs[i]
			utxOutputs = append(utxOutputs, &utxo)
			filledAmount = filledAmount.Sub(utxo.GetAmount())
		}
		if len(pageUtxOutputs) < int(paging.Limit) {
			break
		}
	}

	return utxOutputs, nil
}

func (this *BitcoreUtxoLikeClient) getUtxOutputs(address string, paging meta.Paging) (
	_ []BitcoreUtxoLikeTransactionUtxOutput, err error,
) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/", address)

	var (
		limit  uint = ApiDefaultListPaging
		offset uint = 0
	)
	if paging.Limit > 0 {
		limit = paging.Limit
	}
	if paging.Offset > 0 {
		offset = paging.Offset
	}

	request := this.genRequest().
		SetQueryParam("unspent", "true").
		SetQueryParam("limit", comutils.Stringify(limit)).
		SetQueryParam("index", comutils.Stringify(offset))
	response, err := request.Get(uri)
	if err != nil {
		return
	}
	responseBody, err := this.getResponseBody(response)
	if err != nil {
		return
	}

	var utxOutputs []BitcoreUtxoLikeTransactionUtxOutput
	err = comutils.JsonDecode(responseBody, &utxOutputs)
	if err != nil {
		return
	}

	confirmedOtxOutputs := make([]BitcoreUtxoLikeTransactionUtxOutput, 0, len(utxOutputs))
	for _, utxOutput := range utxOutputs {
		if utxOutput.Confirmations > 0 {
			confirmedOtxOutputs = append(confirmedOtxOutputs, utxOutput)
		}
	}

	return confirmedOtxOutputs, nil
}
