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
		currency     = constants.CurrencyLitecoin
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkLitecoin,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkLitecoinTestnet,
		}
		getterMainnet = func() Coin {
			return &LitecoinCoin{baseCoin{
				currencyIdx: indexMainnet,
				networkMain: indexMainnet.Network,
				networkTest: indexTestnet.Network,
			}}
		}
		getterTestnet = func() Coin {
			return &LitecoinCoin{baseCoin{
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

type LitecoinCoin struct {
	baseCoin
}

func (this LitecoinCoin) GetTradingID() uint16 {
	return 1
}

func (this LitecoinCoin) GetDecimalPlaces() uint8 {
	return BitcoinDecimalPlaces
}

func (this LitecoinCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}
func (this LitecoinCoin) GetChainConfig() *chaincfg.Params {
	if this.IsUsingMainnet() {
		return &LitecoinChainConfig
	} else {
		return &LitecoinChainConfigTestnet
	}
}

func (this LitecoinCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:coin:fee_info:%v", this.GetIndexCode())
	err = comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		30*time.Second,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err = GetLitecoinFeeInfoViaMyCoinGateBitcore()
			if err != nil {
				return nil, err
			}
			feeInfo.LimitMinValue = LitecoinFeeMinimum
			feeInfo.SetLimitMaxValue(LitecoinFeeMaximum, constants.CurrencySubBitcoinSatoshi)
			return feeInfo, nil
		},
	)
	return
}

func (this LitecoinCoin) GetMinTxnAmount() decimal.Decimal {
	minAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(LitecoinTxnMinAmountSatoshi),
		},
		constants.CurrencyLitecoin,
	)
	return minAmount.Value
}

func (this LitecoinCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerLitecoinMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerLitecoinTestnetTxnUrlPattern, txnHash)
	}
}

func (this LitecoinCoin) NormalizeAddress(address string) (string, error) {
	chainConfig := this.GetChainConfig()

	parsedAddress, err := btcutil.DecodeAddress(address, chainConfig)
	if err != nil {
		return "", utils.IssueErrorf("malformed LTC address | address=%s,err=%s", address, err.Error())
	}

	parsedAddressText := parsedAddress.EncodeAddress()
	if IsBitcoinSegWitAddress(parsedAddressText) &&
		!strings.HasPrefix(parsedAddressText, chainConfig.Bech32HRPSegwit+"1") {
		return "", utils.IssueErrorf("malformed LTC SegWit address | address=%s", address)
	}

	return parsedAddress.EncodeAddress(), nil
}

func (this LitecoinCoin) NormalizeAddressLegacy(address string) (string, error) {
	return this.NormalizeAddress(address)
}

func (this LitecoinCoin) NewKey() (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return GenLitecoinMainnetKeyHolder()
	} else {
		return GenLitecoinTestnetKeyHolder()
	}
}

func (this LitecoinCoin) LoadKey(privateKey string, hintAddress string) (KeyHolder, error) {
	if this.IsUsingMainnet() {
		return LoadLitecoinMainnetKeyHolder(privateKey, hintAddress)
	} else {
		return LoadLitecoinTestnetKeyHolder(privateKey, hintAddress)
	}
}

func (this LitecoinCoin) NewClientDefault() (Client, error) {
	if this.IsUsingMainnet() {
		return NewLitecoinBitcoreMainnetClient(), nil
	} else {
		return NewLitecoinSoChainTestnetClient(), nil
	}
}

func (this LitecoinCoin) NewClientSpare() (Client, error) {
	if this.IsUsingMainnet() {
		return NewLitecoinSoChainMainnetClient(), nil
	} else {
		return NewLitecoinSoChainTestnetClient(), nil
	}
}

func (this LitecoinCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := LitecoinAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this LitecoinCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this LitecoinCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this LitecoinCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewLitecoinMainnetSingleTxnSigner(client, feeInfo), nil
	} else {
		return NewLitecoinTestnetSingleTxnSigner(client, feeInfo), nil
	}
}

func (this LitecoinCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (_ BulkTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewLitecoinMainnetBulkTxnSigner(client, feeInfo), nil
	} else {
		return NewLitecoinTestnetBulkTxnSigner(client, feeInfo), nil
	}
}
