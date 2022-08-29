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

type EthereumTokenCoin struct {
	EthereumCoin

	tokenMeta TokenMeta
}

func (this EthereumTokenCoin) genNotSupportCurrencyError() error {
	return utils.IssueErrorf(
		"blockchain module hasn't support ETH token `%v`",
		this.GetCurrency(),
	)
}

func (this EthereumTokenCoin) GetTradingID() uint16 {
	if this.GetCurrency() == constants.CurrencyTetherUSD {
		return 4
	}

	panic(this.genNotSupportCurrencyError())
}

func (this EthereumTokenCoin) GetNetworkCurrency() meta.Currency {
	return constants.CurrencyEthereum
}

func (this EthereumTokenCoin) GetTokenMeta() TokenMeta {
	return this.tokenMeta
}

func (this EthereumTokenCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:token:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		20*time.Second,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err = GetEthereumFeeInfoViaEtherscan()
			if err != nil {
				return nil, err
			}
			feeInfo.SetLimitMaxQuantity(GetConfig().TetherUsdEthereumGasLimit)
			return feeInfo, nil
		},
	)
	return
}

func (this EthereumTokenCoin) NewClientDefault() (Client, error) {
	tokenMeta := this.GetTokenMeta()
	if this.IsUsingMainnet() {
		return NewEthereumTokenBitcoreMainnetClient(tokenMeta)
	} else {
		return NewEthereumTokenEtherscanTestRopstenSystemClient(tokenMeta)
	}
}

func (this EthereumTokenCoin) NewClientSpare() (Client, error) {
	tokenMeta := this.GetTokenMeta()
	if this.IsUsingMainnet() {
		return NewEthereumTokenEtherscanMainnetSystemClient(tokenMeta)
	} else {
		return NewEthereumTokenEtherscanTestRopstenSystemClient(tokenMeta)
	}
}

func (this EthereumTokenCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}

	account := EthereumTokenAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this EthereumTokenCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this EthereumTokenCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this EthereumTokenCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (
	txnSigner SingleTxnSigner, err error,
) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		txnSigner = NewEthereumTokenMainnetTxnSigner(this.tokenMeta, client, feeInfo)
	} else {
		txnSigner = NewEthereumTokenTestRopstenTxnSigner(this.tokenMeta, client, feeInfo)
	}

	return txnSigner, nil
}
