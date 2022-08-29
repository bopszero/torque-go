package blockchainmod

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type baseCoin struct {
	currencyIdx meta.BlockchainCurrencyIndex
	networkMain meta.BlockchainNetwork
	networkTest meta.BlockchainNetwork
}

func (this baseCoin) GetIndexCode() string {
	return fmt.Sprintf("%v@%v", this.GetCurrency(), this.GetNetwork())
}

func (this baseCoin) String() string {
	return this.GetIndexCode()
}

func (this baseCoin) GetCurrency() meta.Currency {
	return this.currencyIdx.Currency
}

func (this baseCoin) GetNetwork() meta.BlockchainNetwork {
	return this.currencyIdx.Network
}

func (this baseCoin) GetNetworkMain() meta.BlockchainNetwork {
	return this.networkMain
}

func (this baseCoin) GetNetworkTest() meta.BlockchainNetwork {
	return this.networkTest
}

func (this baseCoin) IsUsingMainnet() bool {
	return this.GetNetwork() == this.GetNetworkMain()
}

func (this baseCoin) IsAvailable() bool {
	return this.GetModelNetworkCurrency().Priority > 0
}

func (this baseCoin) GetNetworkCurrency() meta.Currency {
	return this.GetCurrency()
}

func (this baseCoin) GetModelInfo() models.CurrencyInfo {
	return currencymod.GetCurrencyInfoFastF(this.GetCurrency())
}

func (this baseCoin) GetModelInfoLegacy() models.LegacyCurrencyInfo {
	return currencymod.GetLegacyCurrencyInfoFastF(this.GetCurrency())
}

func (baseCoin) baseGetModelNetwork(that Coin) models.BlockchainNetworkInfo {
	networkInfo, err := currencymod.GetBlockchainNetworkInfoFast(that.GetNetwork())
	comutils.PanicOnError(err)
	return networkInfo
}

func (this baseCoin) GetModelNetworkCurrency() models.NetworkCurrency {
	return currencymod.GetNetworkCurrencyInfoFastF(this.GetCurrency(), this.GetNetwork())
}

func (this baseCoin) GetMinTxnAmount() decimal.Decimal {
	return decimal.Zero
}

func (baseCoin) baseGetClientAndFeeInfo(that Coin, client Client, offerFeeInfo *FeeInfo) (
	_ Client, _ FeeInfo, err error,
) {
	if client == nil {
		client, err = that.NewClientDefault()
		if err != nil {
			return
		}
	}

	var feeInfo FeeInfo
	if offerFeeInfo != nil {
		feeInfo = *offerFeeInfo
	} else if feeInfo, err = that.GetDefaultFeeInfo(); err != nil {
		return
	}

	return client, feeInfo, nil
}

func (baseCoin) baseNewAccountOwner(coin Coin, privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	keyHolder, err := coin.LoadKey(privateKey, hintAddress)
	if err != nil {
		return
	}
	guestAccount, err := coin.NewAccountGuest(keyHolder.GetAddress())
	if err != nil {
		return
	}
	account := blockchainOwnerAccount{
		GuestAccount: guestAccount,
		coin:         coin,
		keyHolder:    keyHolder,
	}
	return &account, nil
}

func (baseCoin) baseNewAccountSystem(coin Coin, ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	addressInfo, err := GetUserAddressFast(ctx, uid, coin)
	if err != nil {
		return
	}
	keyValue, err := addressInfo.Key.GetValue()
	if err != nil {
		return
	}
	ownerAccount, err := coin.NewAccountOwner(keyValue, addressInfo.Address)
	if err != nil {
		return
	}
	if addressInfo.Address != ownerAccount.GetAddress() && !config.Debug {
		err = utils.IssueErrorf(
			"user address doesn't belong to private key | uid=%v,currency=%v,address=%v",
			uid, coin.GetCurrency(), addressInfo.Address,
		)
		return
	}

	account := blockchainSystemAccount{
		OwnerAccount: ownerAccount,
		addressInfo:  addressInfo,
	}
	return &account, nil
}
