package blockchainmod

type TronAccount struct {
	blockchainGuestAccount
}

func (this *TronAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	// TODO: Implement new account 10K fee
	return this.coin.GetDefaultFeeInfo()
}
