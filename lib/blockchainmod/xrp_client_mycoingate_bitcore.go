package blockchainmod

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreRippleClient struct {
	BitcoreBalanceLikeClient
}

func NewRippleBitcoreMainnetClient() *BitcoreRippleClient {
	return &BitcoreRippleClient{BitcoreBalanceLikeClient{
		tBitcoreClient: newBitcoreClient(
			constants.CurrencyRipple, constants.CurrencySubRippleDrop,
			BitcoreChainMainnet,
		),
	}}
}

func NewRippleBitcoreTestnetClient() *BitcoreRippleClient {
	return &BitcoreRippleClient{BitcoreBalanceLikeClient{
		tBitcoreClient: newBitcoreClient(
			constants.CurrencyRipple, constants.CurrencySubRippleDrop,
			BitcoreChainTestnet,
		),
	}}
}

func (this *BitcoreRippleClient) ParseRootAddress(address string) (string, error) {
	xAddress, err := RippleParseXAddress(address, this.isMainChain())
	if err != nil {
		return "", err
	}
	return xAddress.GetRootAddress(), nil
}

func (this *BitcoreRippleClient) GetBlock(blockHash string) (Block, error) {
	return this.baseGetBlock(this, blockHash)
}

func (this *BitcoreRippleClient) GetBlockByHeight(height uint64) (Block, error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BitcoreRippleClient) GetBlockTxnsByHeight(height uint64) (
	txns []Transaction, err error,
) {
	var (
		uri          = this.genBaseURI() + "/tx/"
		allTxns      []BitcoreRippleTransaction
		pagingOffset = 0
		pagingLimit  = 200
	)
	for {
		var pageTxns []BitcoreRippleTransaction
		request := this.genRequest().SetQueryParams(map[string]string{
			"blockHeight": comutils.Stringify(height),
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
	for i := range txns {
		txn := &allTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}
	return txns, nil
}

func (this *BitcoreRippleClient) GetLatestBlock() (_ Block, err error) {
	latestBlockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}
	return this.GetBlockByHeight(latestBlockHeight)
}

func (this *BitcoreRippleClient) GetTxn(hash string) (_ Transaction, err error) {
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

	var txn BitcoreRippleTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	return &txn, nil
}

func (this *BitcoreRippleClient) GetTxns(address string, paging meta.Paging) (
	_ []Transaction, err error,
) {
	address, err = this.ParseRootAddress(address)
	if err != nil {
		return
	}
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

	var responseTxns []BitcoreRippleTransaction
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

func (this *BitcoreRippleClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	rootAddress, err := this.ParseRootAddress(address)
	if err != nil {
		return
	}
	balance, err := this.BitcoreBalanceLikeClient.GetBalance(rootAddress)
	if err != nil {
		return
	}

	var availableBalance decimal.Decimal
	if balance.GreaterThan(RippleReserveBalance) {
		availableBalance = balance.Sub(RippleReserveBalance)
	} else {
		availableBalance = decimal.Zero
	}
	return availableBalance, nil
}

func (this *BitcoreRippleClient) GetNextNonce(address string) (_ Nonce, err error) {
	rootAddress, err := this.ParseRootAddress(address)
	if err != nil {
		return
	}
	return this.BitcoreBalanceLikeClient.GetNextNonce(rootAddress)
}
