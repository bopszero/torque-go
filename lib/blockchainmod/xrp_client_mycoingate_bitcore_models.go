package blockchainmod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type BitcoreRippleTransaction struct {
	BitcoreBalanceLikeTransaction

	ToTag *uint32 `json:"destinationTag"`
}

func (this *BitcoreRippleTransaction) isConfirmed() bool {
	return this.GetConfirmations() >= RippleMinConfirmations
}

func (this *BitcoreRippleTransaction) SetOwnerAddress(address string) {
	xAddress, err := RippleParseXAddress(address, this.isMainChain())
	comutils.PanicOnError(err)
	this.BitcoreBalanceLikeTransaction.SetOwnerAddress(xAddress.GetRootAddress())
}

func (this *BitcoreRippleTransaction) IsSuccess() bool {
	return this.Status == RippleTxnStatusSuccess
}

func (this *BitcoreRippleTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if !this.isConfirmed() || this.isEmptyStatus() {
		return constants.BlockchainTxnStatusPending
	}
	if this.IsSuccess() {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusFailed
	}
}

func (this *BitcoreRippleTransaction) GetToAddress() (string, error) {
	if this.ToTag == nil {
		return this.ToAddress, nil
	} else {
		return fmt.Sprintf("%s:%d", this.ToAddress, *this.ToTag), nil
	}
}

func (this *BitcoreRippleTransaction) GetToTag() (*uint32, error) {
	return this.ToTag, nil
}
