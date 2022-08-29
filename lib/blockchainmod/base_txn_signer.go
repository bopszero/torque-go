package blockchainmod

type baseTxnSigner struct {
	client          Client
	feeInfo         FeeInfo
	nonce           Nonce
	isPreferOffline bool
}

func (this *baseTxnSigner) GetFeeInfo() FeeInfo {
	return this.feeInfo
}

func (this *baseTxnSigner) SetClient(client Client) {
	this.client = client
}

func (this *baseTxnSigner) SetNonce(nonce Nonce) error {
	this.nonce = nonce
	return nil
}

func (this *baseTxnSigner) MarkAsPreferOffline() {
	this.isPreferOffline = true
}
