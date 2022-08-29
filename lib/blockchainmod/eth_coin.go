package blockchainmod

import (
	"fmt"
	"time"

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
		currency     = constants.CurrencyEthereum
		indexMainnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkEthereum,
		}
		indexTestnet = meta.BlockchainCurrencyIndex{
			Currency: currency,
			Network:  constants.BlockchainNetworkEthereumTestRopsten,
		}
		getterMainnet = func() Coin {
			ethCoin := newEthereumCoin(indexMainnet)
			return &ethCoin
		}
		getterTestnet = func() Coin {
			ethCoin := newEthereumCoin(indexTestnet)
			return &ethCoin
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

type EthereumCoin struct {
	baseCoin
}

func newEthereumCoin(currencyIdx meta.BlockchainCurrencyIndex) EthereumCoin {
	return EthereumCoin{baseCoin{
		currencyIdx: currencyIdx,
		networkMain: constants.BlockchainNetworkEthereum,
		networkTest: constants.BlockchainNetworkEthereumTestRopsten,
	}}
}

func (this EthereumCoin) GetTradingID() uint16 {
	return 3
}

func (this EthereumCoin) GetDecimalPlaces() uint8 {
	return EthereumDecimalPlaces
}

func (this EthereumCoin) GetModelNetwork() models.BlockchainNetworkInfo {
	return this.baseGetModelNetwork(this)
}

func (this EthereumCoin) GetDefaultFeeInfo() (feeInfo FeeInfo, err error) {
	cacheKey := fmt.Sprintf("blockchain:coin:fee_info:%v", this.GetIndexCode())
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
			feeInfo.SetLimitMaxQuantity(EthereumStandardGasLimit)
			return feeInfo, nil
		},
	)
	return
}

func (this EthereumCoin) GenTxnExplorerURL(txnHash string) string {
	if this.IsUsingMainnet() {
		return fmt.Sprintf(ExplorerEthereumMainnetTxnUrlPattern, txnHash)
	} else {
		return fmt.Sprintf(ExplorerEthereumTestRopstenTxnUrlPattern, txnHash)
	}
}

func (this EthereumCoin) NormalizeAddress(address string) (string, error) {
	parsedAddress, err := hexToEthereumAddress(address)
	if err != nil {
		return "", utils.IssueErrorf("malformed ETH address | address=%s,err=%s", address, err.Error())
	}
	return comutils.HexEncode0x(parsedAddress.Bytes()), nil
}

func (this EthereumCoin) NormalizeAddressLegacy(address string) (string, error) {
	return this.NormalizeAddress(address)
}

func (this EthereumCoin) NewKey() (KeyHolder, error) {
	return GenEthereumKeyHolder()
}

func (this EthereumCoin) LoadKey(privateKey string, _ string) (KeyHolder, error) {
	return LoadEthereumKeyHolder(privateKey)
}

func (this EthereumCoin) NewClientDefault() (Client, error) {
	if this.IsUsingMainnet() {
		return NewEthereumBitcoreMainnetClient()
	} else {
		return NewEthereumEtherscanTestRopstenSystemClient()
	}
}

func (this EthereumCoin) NewClientSpare() (Client, error) {
	if this.IsUsingMainnet() {
		return NewEthereumEtherscanMainnetSystemClient()
	} else {
		return NewEthereumEtherscanTestRopstenSystemClient()
	}
}

func (this EthereumCoin) NewAccountGuest(address string) (GuestAccount, error) {
	client, err := this.NewClientDefault()
	if err != nil {
		return nil, err
	}
	account := EthereumAccount{blockchainGuestAccount{
		coin:    this,
		client:  client,
		address: address,
	}}
	return &account, nil
}

func (this EthereumCoin) NewAccountOwner(privateKey string, hintAddress string) (
	_ OwnerAccount, err error,
) {
	return this.baseNewAccountOwner(this, privateKey, hintAddress)
}

func (this EthereumCoin) NewAccountSystem(ctx comcontext.Context, uid meta.UID) (
	_ SystemAccount, err error,
) {
	return this.baseNewAccountSystem(this, ctx, uid)
}

func (this EthereumCoin) NewTxnSignerSingle(client Client, offerFeeInfo *FeeInfo) (_ SingleTxnSigner, err error) {
	client, feeInfo, err := this.baseGetClientAndFeeInfo(this, client, offerFeeInfo)
	if err != nil {
		return
	}

	if this.IsUsingMainnet() {
		return NewEthereumMainnetTxnSigner(client, feeInfo), nil
	} else {
		return NewEthereumTestRopstenTxnSigner(client, feeInfo), nil
	}
}

func (this EthereumCoin) NewTxnSignerBulk(client Client, offerFeeInfo *FeeInfo) (BulkTxnSigner, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't support bulk transaction.")
}
