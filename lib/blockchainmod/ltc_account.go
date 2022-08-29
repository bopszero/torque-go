package blockchainmod

type LitecoinAccount struct {
	blockchainGuestAccount
}

func (this *LitecoinAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	feeInfo, err = this.coin.GetDefaultFeeInfo()
	if err != nil {
		return
	}
	utxOutputCount, err := this.getAddressUtxOutputCountFast(this.GetAddress())
	if err != nil {
		return
	}

	txnSize := EstimateBitcoinLegacyTxnSize(utxOutputCount, BitcoinTxnSingleOutputCount)
	feeInfo.SetLimitMaxQuantity(txnSize)
	return feeInfo, nil
}
