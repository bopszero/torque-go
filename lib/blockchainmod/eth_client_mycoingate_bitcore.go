package blockchainmod

import (
	"fmt"
	"net/http"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreEthereumClient struct {
	BitcoreBalanceLikeClient
	etherscanClient *EtherscanEthereumClient
}

func NewEthereumBitcoreClient() (*BitcoreEthereumClient, error) {
	if config.BlockchainUseTestnet {
		return nil, utils.IssueErrorf("bitcore testnet hasn't been supported yet")
	} else {
		return NewEthereumBitcoreMainnetClient()
	}
}

func NewEthereumBitcoreMainnetClient() (*BitcoreEthereumClient, error) {
	etherscanClient, err := NewEthereumEtherscanClientWithSystemKey()
	if err != nil {
		return nil, err
	}
	client := BitcoreEthereumClient{
		BitcoreBalanceLikeClient: BitcoreBalanceLikeClient{
			tBitcoreClient: newBitcoreClient(
				constants.CurrencyEthereum, constants.CurrencySubEthereumWei,
				BitcoreChainMainnet,
			),
		},
		etherscanClient: etherscanClient,
	}
	return &client, nil
}

func (this *BitcoreEthereumClient) GetBlock(blockHash string) (Block, error) {
	return this.baseGetBlock(this, blockHash)
}

func (this *BitcoreEthereumClient) GetBlockByHeight(height uint64) (Block, error) {
	return this.GetBlock(comutils.Stringify(height))
}

func (this *BitcoreEthereumClient) GetBlockTxnsByHeight(height uint64) (
	txns []Transaction, err error,
) {
	var (
		uri          = this.genBaseURI() + "/tx/"
		allTxns      []BitcoreEthereumTransaction
		pagingOffset = 0
		pagingLimit  = 200
	)
	for {
		var pageTxns []BitcoreEthereumTransaction
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
	for i := range allTxns {
		txn := &allTxns[i]
		txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
		txns[i] = txn
	}
	return txns, nil
}

func (this *BitcoreEthereumClient) GetLatestBlock() (_ Block, err error) {
	latestBlockHeight, err := this.GetLatestBlockHeight()
	if err != nil {
		return
	}
	return this.GetBlockByHeight(latestBlockHeight)
}

func (this *BitcoreEthereumClient) GetTxn(hash string) (_ Transaction, err error) {
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

	var txn BitcoreEthereumTransaction
	err = comutils.JsonDecode(responseBody, &txn)
	if err != nil {
		return
	}

	txn.SetCurrencyPair(this.currencyBase, this.currencyUnit)
	return &txn, nil
}

func (this *BitcoreEthereumClient) GetTxns(address string, paging meta.Paging) (
	_ []Transaction, err error,
) {
	return this.GetNormalTxns(address, paging)
}

func (this *BitcoreEthereumClient) GetNormalTxns(
	address string, paging meta.Paging,
) (_ []Transaction, err error) {
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

	var responseTxns []BitcoreEthereumTransaction
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

func (this *BitcoreEthereumClient) GetRC20TokenTxns(
	address string, contractAddress string, paging meta.Paging,
) (_ []Transaction, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/txs/", address)

	request := this.genRequest().
		SetQueryParam("tokenAddress", contractAddress)
	if paging.Limit > 0 {
		request.SetQueryParam("limit", comutils.Stringify(paging.Limit))
		request.SetQueryParam("index", comutils.Stringify(paging.Offset))
	}
	responseBody, err := this.callRequestGet(request, uri)
	if err != nil {
		return
	}

	var responseTxns []BitcoreEthereumTokenTransaction
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

func (this *BitcoreEthereumClient) PushTxnRaw(data []byte) (err error) {
	if err = this.BitcoreBalanceLikeClient.PushTxnRaw(data); err == nil {
		return nil
	}
	if etherscanErr := this.etherscanClient.PushTxnRaw(data); etherscanErr == nil {
		return nil
	}
	return err
}
