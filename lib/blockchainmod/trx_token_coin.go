package blockchainmod

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type TronTokenCoin struct {
	TronCoin

	tokenMeta TokenMeta
}

func (this TronTokenCoin) genNotSupportCurrencyError() error {
	return utils.IssueErrorf(
		"blockchain module hasn't support TRX token `%v`",
		this.GetCurrency(),
	)
}

func (this TronTokenCoin) GetTradingID() uint16 {
	if this.GetCurrency() == constants.CurrencyTetherUSD {
		return 4
	}

	panic(this.genNotSupportCurrencyError())
}

func (this TronTokenCoin) GetNetworkCurrency() meta.Currency {
	return constants.CurrencyTron
}

func (this TronTokenCoin) GetTokenMeta() TokenMeta {
	return this.tokenMeta
}

func (this TronTokenCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:token:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		5*time.Minute,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err = GetTronFeeInfoViaMyCoinGateBitcore()
			if err != nil {
				return nil, err
			}
			feeInfo.SetLimitMaxQuantity(GetConfig().TetherUsdTronEnergyLimit)
			return feeInfo, nil
		},
	)
	return
}

func (this TronTokenCoin) NewClientDefault() (Client, error) {
	tokenMeta := this.GetTokenMeta()
	if this.IsUsingMainnet() {
		return NewTronTokenBitcoreMainnetClient(tokenMeta), nil
	} else {
		return NewTronTokenBitcoreTestShastaClient(tokenMeta), nil
	}
}

func (this TronTokenCoin) NewClientSpare() (Client, error) {
	return this.NewClientDefault()
}

func (this TronTokenCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}

	account := TronTokenAccount{TronAccount{
		blockchainGuestAccount{
			coin:    this,
			client:  client,
			address: address,
		},
	}}
	return &account, nil
}

func (this TronTokenCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this TronTokenCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this TronTokenCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (
	txnSigner SingleTxnSigner, err error,
) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		txnSigner = NewTronTokenMainnetTxnSigner(this.tokenMeta, client, feeInfo)
	} else {
		txnSigner = NewTronTokenTestShastaTxnSigner(this.tokenMeta, client, feeInfo)
	}

	return txnSigner, nil
}
