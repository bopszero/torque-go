package blockchainmod

import (
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreEthereumTokenTransaction struct {
	BitcoreEthereumTransaction

	Token   BitcoreEthereumTransactionTokenInfo    `json:"token"`
	AbiData BitcoreEthereumTokenTransactionAbiInfo `json:"abiType"`
}

type BitcoreEthereumTokenTransactionAbiInfo struct {
	Type   string                                    `json:"type"`
	Name   string                                    `json:"name"`
	Params []BitcoreEthereumTokenTransactionAbiParam `json:"params"`
}

type BitcoreEthereumTokenTransactionAbiParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type BitcoreEthereumTransactionTokenInfo struct {
	Name          string          `json:"name"`
	DecimalPlaces int32           `json:"decimals"`
	Currency      meta.Currency   `json:"symbol"`
	Address       string          `json:"address"`
	Amount        decimal.Decimal `json:"amount"`
}

func (this *BitcoreEthereumTokenTransaction) IsTransferTxn() bool {
	return this.AbiData.Name == "transfer" || this.AbiData.Name == "transferFrom"
}

func (this *BitcoreEthereumTokenTransaction) GetCurrency() meta.Currency {
	return this.Token.Currency
}

func (this *BitcoreEthereumTokenTransaction) GetToAddress() (string, error) {
	if !this.IsTransferTxn() {
		return this.ToAddress, nil
	}

	var toAddress string
	for _, param := range this.AbiData.Params {
		if param.Type == "address" && param.Name == "_to" {
			toAddress = param.Value
			break
		}
	}
	if toAddress == "" {
		return "", utils.IssueErrorf("ETH token ABI data doesn't contains to address | params=%v", this.AbiData.Params)
	}
	return toAddress, nil
}

func (this *BitcoreEthereumTokenTransaction) GetDirection() meta.Direction {
	toAddress, err := this.GetToAddress()
	if err != nil {
		return constants.DirectionTypeUnknown
	}

	switch strings.ToLower(this.ownerAddress) {

	case strings.ToLower(this.FromAddress):
		return constants.DirectionTypeSend

	case strings.ToLower(toAddress):
		return constants.DirectionTypeReceive

	default:
		if this.GetCurrency() != constants.CurrencyEthereum {
			return constants.DirectionTypeReceive
		}

		return constants.DirectionTypeUnknown
	}
}

func (this *BitcoreEthereumTokenTransaction) GetAmount() (_ decimal.Decimal, err error) {
	if !this.IsTransferTxn() {
		return decimal.Zero, nil
	}

	var amountValue string
	for _, param := range this.AbiData.Params {
		if param.Type == "uint256" && param.Name == "_value" {
			amountValue = param.Value
			break
		}
	}
	if amountValue == "" {
		err = utils.IssueErrorf("ETH token ABI data doesn't contains amount | params=%v", this.AbiData.Params)
		return
	}
	amount, err := decimal.NewFromString(amountValue)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	return amount.Shift(-this.Token.DecimalPlaces), nil
}
