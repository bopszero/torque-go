package blockchainmod

import (
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
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
		currency     = constants.CurrencyBitcoin
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoin,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkBitcoinTestnet,
		}
		getterMainnet = func() Coin {
			return &BitcoinCoin{baseCoin{
				currencyIdx: indexMainnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}
		}
		getterTestnet = func() Coin {
			return &BitcoinCoin{baseCoin{
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

type BitcoinCoin struct {
	baseCoin
}

func (this BitcoinCoin) GetTradingID() uint16 {
	return 2
}

func (this BitcoinCoin) GetDecimalPlaces() uint8 {
	return BitcoinDecimalPlaces
}

func (this BitcoinCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}

func (this BitcoinCoin) GetChainConfig() *chaincfg.Params {
	if this.IsUsingMainnet() {
		return &BitcoinChainConfig
	} else {
		return &BitcoinChainConfigTestnet
	}
}

func (this BitcoinCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:coin:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		30*time.Second,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err = GetBitcoinFeeInfoViaMyCoinGateBitcore()
			if err != nil {
				return nil, err
			}

			feeInfo.LimitMinValue = BitcoinFeeMinimum
			return feeInfo, nil
		},
	)
	return
}

func (this BitcoinCoin) GetMinTxnAmount() decimal.Decimal {
	minAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(BitcoinTxnMinAmountSatoshi),
		},
		this.GetCurrency(),
	)
	return minAmount.Value
}

func (this BitcoinCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerBitcoinMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerBitcoinTestnetTxnUrlPattern, txnHash)
	}
}

func (this BitcoinCoin) NormalizeAddress(address string) (string, error) {
	chainConfig := this.GetChainConfig()

	parsedAddress, err := btcutil.DecodeAddress(address, chainConfig)
	if err != nil {
		return "", utils.IssueErrorf("malformed BTC address | address=%s,err=%s", address, err.Error())
	}

	parsedAddressText := parsedAddress.EncodeAddress()
	if IsBitcoinSegWitAddress(parsedAddressText) &&
		!strings.HasPrefix(parsedAddressText, chainConfig.Bech32HRPSegwit+"1") {
		return "", utils.IssueErrorf("malformed BTC SegWit address | address=%s", address)
	}

	return parsedAddress.EncodeAddress(), nil
}

func (this BitcoinCoin) NormalizeAddressLegacy(address string) (string, error) {
	return this.NormalizeAddress(address)
}

func (this BitcoinCoin) NewKey() (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return GenBitcoinMainnetKeyHolder()
	} else {
		return GenBitcoinTestnetKeyHolder()
	}
}

func (this BitcoinCoin) LoadKey(privateKey string, hintAddress string) (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return LoadBitcoinMainnetKeyHolder(privateKey, hintAddress)
	} else {
		return LoadBitcoinTestnetKeyHolder(privateKey, hintAddress)
	}
}

func (this BitcoinCoin) NewClientDefault() (Client, error) {
	if this.IsUsingMainnet() {
		return NewBitcoinBitcoreMainnetClient(), nil
	} else {
		return NewBitcoinSoChainTestnetClient(), nil
	}
}

func (this BitcoinCoin) NewClientSpare() (Client, error) {
	if this.IsUsingMainnet() {
		return NewBitcoinSoChainMainnetClient(), nil
	} else {
		return NewBitcoinSoChainTestnetClient(), nil
	}
}

func (this BitcoinCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := BitcoinAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this BitcoinCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this BitcoinCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this BitcoinCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewBitcoinMainnetSingleTxnSigner(client, feeInfo), nil
	} else {
		return NewBitcoinTestnetSingleTxnSigner(client, feeInfo), nil
	}
}

func (this BitcoinCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (_ BulkTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewBitcoinMainnetBulkTxnSigner(client, feeInfo), nil
	} else {
		return NewBitcoinTestnetBulkTxnSigner(client, feeInfo), nil
	}
}
