package blockchainmod

type BitcoinAccount struct {
	blockchainGuestAccount
}

func (this *BitcoinAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	feeInfo, err = this.coin.GetDefaultFeeInfo()
	if err != nil {
		return
	}
	utxOutputCount, err := this.getAddressUtxOutputCountFast(this.GetAddress())
	if err != nil {
		return
	}

	var txnSize uint32
	if IsBitcoinSegWitAddress(this.GetAddress()) {
		txnSize = EstimateBitcoinSegwitTxnSize(utxOutputCount, BitcoinTxnSingleOutputCount)
	} else {
		txnSize = EstimateBitcoinLegacyTxnSize(utxOutputCount, BitcoinTxnSingleOutputCount)
	}
	feeInfo.SetLimitMaxQuantity(txnSize)

	return feeInfo, nil
}
