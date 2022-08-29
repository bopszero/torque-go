package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreTronTokenTransaction struct {
	BitcoreTronTransaction
	rc20Transfer RC20Transfer
	tokenMeta    TokenMeta
}

func (this *BitcoreTronTokenTransaction) initTokenMeta(tokenMeta TokenMeta) error {
	if this.rc20Transfer != nil {
		return nil
	}
	transfers, err := this.GetRC20Transfers(tokenMeta)
	if err != nil {
		return err
	}
	if len(transfers) == 0 {
		return utils.WrapError(constants.ErrorDataNotFound)
	}

	this.tokenMeta = tokenMeta
	this.rc20Transfer = transfers[0]
	return nil
}

func (this *BitcoreTronTokenTransaction) GetToAddress() (string, error) {
	return this.rc20Transfer.GetToAddress(), nil
}

func (this *BitcoreTronTokenTransaction) GetAmount() (decimal.Decimal, error) {
	return this.rc20Transfer.GetAmount(), nil
}
