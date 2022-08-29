package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type BitcoreTronTransaction struct {
	BitcoreBalanceLikeTransaction

	Type    string                             `json:"type"`
	AbiInfo BitcoreTronTokenTransactionAbiInfo `json:"tronAbi"`
}

func (this *BitcoreTronTransaction) isConfirmed() bool {
	return this.GetConfirmations() >= TronMinConfirmations
}

func (this *BitcoreTronTransaction) IsSuccess() bool {
	return this.Status == TronTxnStatusSuccess
}

func (this *BitcoreTronTransaction) IsTransferContract() bool {
	return this.Type == TronTxnTypeTransferContract
}

func (this *BitcoreTronTransaction) IsTransferAssetContract() bool {
	return this.Type == TronTxnTypeTransferAssetContract
}

func (this *BitcoreTronTransaction) IsTriggerSmartContract() bool {
	return this.Type == TronTxnTypeTriggerSmartContract
}

func (this *BitcoreTronTransaction) GetTypeCode() string {
	return this.Type
}

func (this *BitcoreTronTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if !this.isConfirmed() || this.isEmptyStatus() {
		return constants.BlockchainTxnStatusPending
	}
	if this.IsSuccess() {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusFailed
	}
}

func (this *BitcoreTronTransaction) GetRC20Transfers(tokenMeta TokenMeta) ([]RC20Transfer, error) {
	if this.ToAddress != tokenMeta.Address {
		return nil, nil
	}
	if this.AbiInfo.Type != TronTRC20TransferTokenMethodType ||
		this.AbiInfo.Name != TronTRC20TransferTokenMethodName {
		return nil, nil
	}

	var (
		toAddress string
		amount    decimal.Decimal
	)
	for _, param := range this.AbiInfo.Params {
		switch {
		case param.Name == "_to" && param.Type == "address":
			toAddress = param.Value
			break
		case param.Name == "_value" && param.Type == "uint256":
			rawAmount, err := decimal.NewFromString(param.Value)
			comutils.PanicOnError(err)
			amount = rawAmount.Shift(-int32(tokenMeta.DecimalPlaces))
			break
		default:
			break
		}
	}

	transfer := BitcoreTronTransactionERC20Transfer{
		TokenMeta:   tokenMeta,
		FromAddress: this.FromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}
	return []RC20Transfer{&transfer}, nil
}

type BitcoreTronTokenTransactionAbiInfo struct {
	Type   string                                 `json:"type"`
	Name   string                                 `json:"name"`
	Params []BitcoreTronTokenTransactionAbiParams `json:"params"`
}

type BitcoreTronTokenTransactionAbiParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type BitcoreTronTransactionERC20Transfer struct {
	TokenMeta   TokenMeta
	FromAddress string
	ToAddress   string
	Amount      decimal.Decimal
}

func (this *BitcoreTronTransactionERC20Transfer) GetTokenMeta() TokenMeta {
	return this.TokenMeta
}

func (this *BitcoreTronTransactionERC20Transfer) GetFromAddress() string {
	return this.FromAddress
}

func (this *BitcoreTronTransactionERC20Transfer) GetToAddress() string {
	return this.ToAddress
}

func (this *BitcoreTronTransactionERC20Transfer) GetAmount() decimal.Decimal {
	return this.Amount
}
