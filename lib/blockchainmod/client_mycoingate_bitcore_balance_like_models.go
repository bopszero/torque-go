package blockchainmod

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type BitcoreBalanceLikeBlock struct {
	tBitcoreBaseBlock
}

type BitcoreBalanceLikeTransaction struct {
	tBitcoreBaseTransaction

	Status      string `json:"status"`
	DataHex     string `json:"data"`
	Nonce       uint64 `json:"nonce"`
	ToAddress   string `json:"to"`
	FromAddress string `json:"from"`
}

func (this *BitcoreBalanceLikeTransaction) isEmptyStatus() bool {
	return this.Status == ""
}

func (this *BitcoreBalanceLikeTransaction) GetDirection() meta.Direction {
	switch this.ownerAddress {
	case this.FromAddress:
		return constants.DirectionTypeSend
	case this.ToAddress:
		return constants.DirectionTypeReceive
	default:
		return constants.DirectionTypeUnknown
	}
}

func (this *BitcoreBalanceLikeTransaction) GetFromAddress() (string, error) {
	return this.FromAddress, nil
}

func (this *BitcoreBalanceLikeTransaction) GetToAddress() (string, error) {
	return this.ToAddress, nil
}

func (this *BitcoreBalanceLikeTransaction) GetInputDataHex() (string, error) {
	return this.DataHex, nil
}
