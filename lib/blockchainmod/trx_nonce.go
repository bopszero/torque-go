package blockchainmod

type TronFrozenBlockNonce struct {
	block Block
}

func NewTronFrozenBlockNonce(block Block) Nonce {
	return &TronFrozenBlockNonce{block}
}

func (this *TronFrozenBlockNonce) GetNumber() (uint64, error) {
	return this.block.GetHeight(), nil
}

func (this *TronFrozenBlockNonce) GetValue() interface{} {
	return this.block
}

func (this *TronFrozenBlockNonce) Next() (Nonce, error) {
	return this, nil
}
