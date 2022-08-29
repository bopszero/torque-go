package meta

import "strings"

type BlockchainNetwork string

func NewBlockchainNetwork(code string) BlockchainNetwork {
	return BlockchainNetwork(strings.ToUpper(code))
}

func (this BlockchainNetwork) String() string {
	return string(this)
}

type BlockchainCurrencyIndex struct {
	Currency Currency
	Network  BlockchainNetwork
}

type BlockchainTxnStatus int8
