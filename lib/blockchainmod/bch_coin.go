package blockchainmod

import (
	"fmt"
	"time"

	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func init() {
	var (
		currency     = constants.CurrencyBitcoinCash
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoinCash,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoinCashTestnet,
		}
		getterMainnet = func() Coin {
			return &BitcoinCashCoin{baseCoin{
				currencyIdx: indexMainnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}
		}
		getterTestnet = func() Coin {
			return &BitcoinCashCoin{baseCoin{
				currencyIdx: indexTestnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}
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

type BitcoinCashCoin struct {
	baseCoin
}

func (this BitcoinCashCoin) GetTradingID() uint16 {
	return 20
}

func (this BitcoinCashCoin) GetDecimalPlaces() uint8 {
	return BitcoinDecimalPlaces
}

func (this BitcoinCashCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}

func (this BitcoinCashCoin) GetChainConfig() *chaincfg.Params {
	if this.IsUsingMainnet() {
		return &BitcoinCashChainConfig
	} else {
		return &BitcoinCashChainConfigTestnet
	}
}

func (this BitcoinCashCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:coin:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		30*time.Second,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err = GetFeeResponseFromMyCoinGateBitcore(this.GetCurrency(), 3)
			if err != nil {
				return nil, err
			}

			feeInfo.LimitMinValue = BitcoinCashFeeMinimum
			return feeInfo, nil
		},
	)
	return
}

func (this BitcoinCashCoin) GetMinTxnAmount() decimal.Decimal {
	minAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(BitcoinTxnMinAmountSatoshi),
		},
		this.GetCurrency(),
	)
	return minAmount.Value
}

func (this BitcoinCashCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerBitcoinCashMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerBitcoinCashTestnetTxnUrlPattern, txnHash)
	}
}

func (this BitcoinCashCoin) NormalizeAddress(address string) (string, error) {
	parsedAddress, err := bchutil.DecodeAddress(address, this.GetChainConfig())
	if err != nil {
		return "", utils.IssueErrorf("malformed BCH address | address=%s,err=%s", address, err.Error())
	}
	return parsedAddress.EncodeAddress(), nil
}

func (this BitcoinCashCoin) NormalizeAddressLegacy(address string) (string, error) {
	return this.NormalizeAddress(address)
}

func (this BitcoinCashCoin) NewKey() (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return GenBitcoinCashMainnetKeyHolder()
	} else {
		return GenBitcoinCashTestnetKeyHolder()
	}
}

func (this BitcoinCashCoin) LoadKey(privateKey string, hintAddress string) (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return LoadBitcoinCashMainnetKeyHolder(privateKey, hintAddress)
	} else {
		return LoadBitcoinCashTestnetKeyHolder(privateKey, hintAddress)
	}
}

func (this BitcoinCashCoin) NewClientDefault() (Client, error) {
	var client BitcoreUtxoLikeClient
	if this.IsUsingMainnet() {
		client = NewBitcoreUtxoLikeClient(this.GetCurrency(), BitcoreChainMainnet)
	} else {
		client = NewBitcoreUtxoLikeClient(this.GetCurrency(), BitcoreChainTestnet)
	}
	return &client, nil
}

func (this BitcoinCashCoin) NewClientSpare() (Client, error) {
	return this.NewClientDefault()
}

func (this BitcoinCashCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := BitcoinCashAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this BitcoinCashCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this BitcoinCashCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this BitcoinCashCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewBitcoinCashMainnetSingleTxnSigner(this, client, feeInfo), nil
	} else {
		return NewBitcoinCashTestnetSingleTxnSigner(this, client, feeInfo), nil
	}
}

func (this BitcoinCashCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (_ BulkTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewBitcoinCashMainnetBulkTxnSigner(this, client, feeInfo), nil
	} else {
		return NewBitcoinCashTestnetBulkTxnSigner(this, client, feeInfo), nil
	}
}
