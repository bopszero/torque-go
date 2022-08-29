package blockchainmod

type EthereumTokenAccount struct {
	blockchainGuestAccount
}

func (this *EthereumTokenAccount) GetFeeInfoToAddress(toAddress string) (
	feeInfo FeeInfo, err error,
) {
	return this.coin.GetDefaultFeeInfo()
}
