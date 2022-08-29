package blockchainmod

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type BitcoreBalanceLikeClient struct {
	tBitcoreClient
}

type bitcoreBalanceLikeGetBalanceResponse struct {
	BalanceUnit decimal.Decimal `json:"balance"`
}

func (this *BitcoreBalanceLikeClient) GetBalance(address string) (_ decimal.Decimal, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/balance/", address)
	responseBody, err := this.callRequestGet(this.genRequest(), uri)
	if err != nil {
		return
	}

	var responseModel bitcoreBalanceLikeGetBalanceResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}
	balanceBase := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.currencyUnit,
			Value:    responseModel.BalanceUnit,
		},
		this.currencyBase,
	)
	return balanceBase.Value, nil
}

type bitcoreGetTxnsCountResponse struct {
	Nonce uint64 `json:"nonce"`
}

func (this *BitcoreBalanceLikeClient) GetNextNonce(address string) (_ Nonce, err error) {
	uri := this.genBaseURI() + fmt.Sprintf("/address/%s/txs/count/", address)
	responseBody, err := this.callRequestGet(this.genRequest(), uri)
	if err != nil {
		return
	}

	var responseModel bitcoreGetTxnsCountResponse
	err = comutils.JsonDecode(responseBody, &responseModel)
	if err != nil {
		return
	}

	return NewNumberNonce(responseModel.Nonce), nil
}
