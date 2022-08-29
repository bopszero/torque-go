package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

type BitcoinCashAccount struct {
	blockchainGuestAccount
}

func (this *BitcoinCashAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	feeInfo, err = this.coin.GetDefaultFeeInfo()
	if err != nil {
		return
	}
	utxOutputCount, err := this.getAddressUtxOutputCountFast(this.GetAddress())
	if err != nil {
		return
	}

	var (
		txnSize = EstimateBitcoinLegacyTxnSize(utxOutputCount, BitcoinTxnSingleOutputCount)
		txnFee  = feeInfo.Price.Mul(decimal.NewFromInt(int64(txnSize)))
	)
	feeInfo.LimitMaxQuantity = txnSize
	feeInfo.SetLimitMaxValue(txnFee, constants.CurrencySubBitcoinSatoshi)

	return feeInfo, nil
}
