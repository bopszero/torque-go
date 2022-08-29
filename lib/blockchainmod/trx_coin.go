package blockchainmod

import (
	"fmt"
	"time"

	tronaddr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"gitlab.com/snap-clickstaff/go-common/comcache"
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
		currency     = constants.CurrencyTron
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkTron,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkTronTestShasta,
		}
		getterMainnet = func() Coin {
			tronCoin := newTronCoin(indexMainnet)
			return &tronCoin
		}
		getterTestnet = func() Coin {
			tronCoin := newTronCoin(indexTestnet)
			return &tronCoin
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

type TronCoin struct {
	baseCoin
}

func newTronCoin(currencyIdx meta.BlockchainCurrencyIndex) TronCoin {
	return TronCoin{baseCoin{
		currencyIdx: currencyIdx,
		networkMain: constants.BlockchainNetworkTron,
		networkTest: constants.BlockchainNetworkTronTestShasta,
	}}
}

func (this TronCoin) GetTradingID() uint16 {
	return 50
}

func (this TronCoin) GetDecimalPlaces() uint8 {
	return TronDecimalPlaces
}

func (this TronCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}

func (this TronCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:coin:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		5*time.Minute,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err := GetTronFeeInfoViaMyCoinGateBitcore()
			if err != nil {
				return nil, err
			}
			feeInfo.SetLimitMaxValue(TronFeeMaximum, constants.CurrencySubTronSun)
			return feeInfo, nil
		},
	)
	return
}

func (this TronCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerTronMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerTronTestShastaTxnUrlPattern, txnHash)
	}
}

func (this TronCoin) NormalizeAddress(address string) (string, error) {
	parsedAddress, err := tronaddr.Base58ToAddress(address)
	if err == nil && parsedAddress.Bytes()[0] != tronaddr.TronBytePrefix {
		err = fmt.Errorf("invalid tron address prefix")
	}
	if err != nil {
		return "", utils.IssueErrorf("malformed TRX address | address=%s,err=%s", address, err.Error())
	}

	return parsedAddress.String(), nil
}

func (this TronCoin) NormalizeAddressLegacy(address string) (string, error) {
	return this.NormalizeAddress(address)
}

func (this TronCoin) NewKey() (KeyHolder, error) {
	return GenTronKeyHolder()
}

func (this TronCoin) LoadKey(privateKey string, hintAddress string) (KeyHolder, error) {
	return LoadTronKeyHolder(privateKey)
}

func (this TronCoin) NewClientDefault() (Client, error) {
	if this.IsUsingMainnet() {
		return NewTronBitcoreMainnetClient(), nil
	} else {
		return NewTronBitcoreTestShastaClient(), nil
	}
}

func (this TronCoin) NewClientSpare() (Client, error) {
	return this.NewClientDefault()
}

func (this TronCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := TronAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this TronCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this TronCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this TronCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewTronMainnetTxnSigner(client, feeInfo), nil
	} else {
		return NewTronTestShastaTxnSigner(client, feeInfo), nil
	}
}

func (this TronCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (_ BulkTxnSigner, err error) {
	return nil, utils.IssueErrorf("Tron doesn't support bulk transaction.")
}
