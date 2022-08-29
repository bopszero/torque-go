package blockchainmod

import (
	"fmt"

	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func init() {
	var (
		currency     = constants.CurrencyRipple
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkRipple,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkRippleTestnet,
		}
		getterMainnet = func() Coin {
			rippleCoin := newRippleCoin(indexMainnet)
			return &rippleCoin
		}
		getterTestnet = func() Coin {
			rippleCoin := newRippleCoin(indexTestnet)
			return &rippleCoin
		}
	)
	comutils.PanicOnError(
		RegisterNativeCoinLoader(
			currency,
			func() Coin {
				if config.BlockchainUseTestnet {
					return getterTestnet()
				} else {
					return getterMainnet()
				}
			},
		),
	)
	comutils.PanicOnError(
		RegisterCoinLoader(indexMainnet, getterMainnet),
	)
	comutils.PanicOnError(
		RegisterCoinLoader(indexTestnet, getterTestnet),
	)
}

type RippleCoin struct {
	baseCoin
}

func newRippleCoin(currencyIdx meta.BlockchainCurrencyIndex) RippleCoin {
	return RippleCoin{baseCoin{
		currencyIdx: currencyIdx,
		networkMain: constants.BlockchainNetworkRipple,
		networkTest: constants.BlockchainNetworkRippleTestnet,
	}}
}

func (this RippleCoin) GetTradingID() uint16 {
	return 51
}

func (this RippleCoin) GetDecimalPlaces() uint8 {
	return RippleDecimalPlaces
}

func (this RippleCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}

func (this RippleCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	feeInfo, err = GetRippleFeeInfoViaMyCoinGateBitcore()
	if err != nil {
		return
	}
	feeInfo.SetLimitMaxQuantity(1)
	return
}

func (this RippleCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerRippleMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerRippleTestnetTxnUrlPattern, txnHash)
	}
}

func (this RippleCoin) NormalizeAddress(address string) (string, error) {
	xAddress, err := RippleParseXAddress(address, this.IsUsingMainnet())
	if err != nil {
		return "", err
	} else {
		return xAddress.String(), nil
	}
}

func (this RippleCoin) NormalizeAddressLegacy(address string) (string, error) {
	xAddress, err := RippleParseXAddress(address, this.IsUsingMainnet())
	if err != nil {
		return "", err
	} else {
		return xAddress.GetRootTagAddress(), nil
	}
}

func (this RippleCoin) NewKey() (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return GenRippleMainnetKeyHolder()
	} else {
		return GenRippleTestnetKeyHolder()
	}
}

func (this RippleCoin) LoadKey(privateKey string, hintAddress string) (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return LoadRippleMainnetKeyHolder(privateKey, hintAddress)
	} else {
		return LoadRippleTestnetKeyHolder(privateKey, hintAddress)
	}
}

func (this RippleCoin) NewClientDefault() (Client, error) {
	if this.IsUsingMainnet() {
		return NewRippleBitcoreMainnetClient(), nil
	} else {
		return NewRippleBitcoreTestnetClient(), nil
	}
}

func (this RippleCoin) NewClientSpare() (Client, error) {
	return this.NewClientDefault()
}

func (this RippleCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := RippleAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this RippleCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this RippleCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this RippleCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewRippleMainnetTxnSigner(client, feeInfo), nil
	} else {
		return NewRippleTestnetTxnSigner(client, feeInfo), nil
	}
}

func (this RippleCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (_ BulkTxnSigner, err error) {
	return nil, utils.IssueErrorf("Ripple doesn't support bulk transaction.")
}
