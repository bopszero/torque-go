package blockchainmod

import (
	"fmt"
	"net/http"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreTronClient struct {
	BitcoreBalanceLikeClient
}

func NewTronBitcoreMainnetClient() *BitcoreTronClient {
	return &BitcoreTronClient{BitcoreBalanceLikeClient{
		tBitcoreClient: newBitcoreClient(
			constants.CurrencyTron, constants.CurrencySubTronSun,
			BitcoreChainMainnet,
		),
	}}
}

func NewTronBitcoreTestShastaClient() *BitcoreTronClient {
	return &BitcoreTronClient{BitcoreBalanceLikeClient{
		tBitcoreClient: newBitcoreClient(
			constants.CurrencyTron, constants.CurrencySubTronSun,
			BitcoreChainTestShasta,
		),
	}}
}

func (this *BitcoreTronClient) GetBlock(blockHash string) (Block, error) {
	return this.baseGetBlock(this, blockHash)
}

func (this *BitcoreTronClient) GetBlockByHeight(height uint64) (Block, error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BitcoreTronClient) GetBlockTxnsByHeight(height uint64) (
	txns []Transaction, err error,
) {
	uri := this.genBaseURI() + "/tx/"
	request := this.genRequest().SetQueryParams(map[string]string{
		"blockHeight": comutils.Stringify(height),
		"limit":       "5000",
	})
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}
	var tronTxns []BitcoreTronTransaction
	err = comutils.JsonDecode(responseBody, &tronTxns)
	if err != nil {
		return
	}

	txns = make([]Transaction, len(tronTxns))
	for i := range tronTxns {
		txn := &tronTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}
	return txns, nil
}

func (this *BitcoreTronClient) GetLatestBlock() (_ Block, err error) {
	latestBlockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}
	return this.GetBlockByHeight(latestBlockHeight)
}

func (this *BitcoreTronClient) GetTxn(hash string) (_ Transaction, err error) {
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

	var txn BitcoreTronTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	return &txn, nil
}

func (this *BitcoreTronClient) GetTxns(address string, paging meta.Paging) (
	[]Transaction, error,
) {
	return this.GetNormalTxns(address, paging)
}

func (this *BitcoreTronClient) GetNormalTxns(address string, paging meta.Paging) (
	_ []Transaction, err error,
) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/txs/", address)

	request := this.genRequest()
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
		request.SetQueryParam("index", comutils.Stringify(paging.Offset))
	}
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseTxns []BitcoreTronTransaction
	err = comutils.JsonDecode(responseBody, &responseTxns)
	if err != nil {
		return
	}

	txns := make([]Transaction, len(responseTxns))
	for i := range responseTxns {
		txn := &responseTxns[i]
		txn.SetOwnerAddress(address)
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}

	return txns, nil
}

func (this *BitcoreTronClient) GetRC20TokenTxns(address string, tokenMeta TokenMeta, paging meta.Paging) (
	_ []Transaction, err error,
) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/txs/", address)

	request := this.genRequest().
		SetQueryParam("tokenAddress", tokenMeta.Address)
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
		request.SetQueryParam("index", comutils.Stringify(paging.Offset))
	}
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseTxns []BitcoreTronTokenTransaction
	err = comutils.JsonDecode(responseBody, &responseTxns)
	if err != nil {
		return
	}

	txns := make([]Transaction, 0, len(responseTxns))
	for i := range responseTxns {
		txn := &responseTxns[i]
		txn.SetOwnerAddress(address)
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		if err = txn.initTokenMeta(tokenMeta); err != nil {
			if !utils.IsOurError(err, constants.ErrorCodeDataNotFound) {
				return
			}
			continue
		}
		txns = append(txns, txn)
	}

	return txns, nil
}

func (this *BitcoreTronClient) GetNextNonce(address string) (_ Nonce, err error) {
	latestBlock, err := this.GetLatestBlock()
	if err != nil {
		return
	}

	return NewTronFrozenBlockNonce(latestBlock), nil
}
