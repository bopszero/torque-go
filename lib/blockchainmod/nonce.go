package blockchainmod

type NumberNonce uint64

func NewNumberNonce(value uint64) Nonce {
	nonce := NumberNonce(value)
	return &nonce
}

func (this *NumberNonce) GetNumber() (uint64, error) {
	return uint64(*this), nil
}

func (this *NumberNonce) GetValue() interface{} {
	return uint64(*this)
}

func (this *NumberNonce) Next() (Nonce, error) {
	return NewNumberNonce(uint64(*this) + 1), nil
}
