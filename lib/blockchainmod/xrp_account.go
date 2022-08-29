package blockchainmod

type RippleAccount struct {
	blockchainGuestAccount
}

func (this *RippleAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	return this.coin.GetDefaultFeeInfo()
}
