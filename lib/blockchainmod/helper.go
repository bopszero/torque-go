package blockchainmod

import (
	"crypto/ecdsa"
	"math"
	"reflect"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func CreateEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func GetBlockchainOrderStatus(status meta.BlockchainTxnStatus) meta.OrderStatus {
	switch status {
	case constants.BlockchainTxnStatusPending:
		return constants.OrderStatusHandleDst
	case constants.BlockchainTxnStatusSucceeded:
		return constants.OrderStatusCompleted
	default:
		return constants.OrderStatusFailed
	}
}

func isSochainTestnet(network string) bool {
	return strings.HasSuffix(
		strings.ToLower(network),
		soChainNetworkTestSuffix,
	)
}

func getSochainNetworkCurrency(network string) meta.Currency {
	if isSochainTestnet(network) {
		return meta.NewCurrencyU(network[:len(network)-4])
	} else {
		return meta.NewCurrencyU(network)
	}
}

func EstimateBitcoinLegacyTxnSize(inputCount uint32, outputCount uint32) uint32 {
	return EstimateBitcoinMixedTxnSize(0, inputCount, outputCount)
}

func EstimateBitcoinSegwitTxnSize(inputCount uint32, outputCount uint32) uint32 {
	return EstimateBitcoinMixedTxnSize(inputCount, 0, outputCount)
}

func EstimateBitcoinMixedTxnSize(
	inputCountSegwit uint32, inputCountLegacy uint32, outputCount uint32,
) uint32 {
	return BitcoinTxnSizeMeta +
		uint32(math.Ceil(float64(inputCountSegwit)*BitcoinTxnSegWitSizeInput)) +
		inputCountLegacy*BitcoinTxnLegacySizeInput +
		outputCount*BitcoinTxnSizeOutput
}

func IsNetworkCoin(coin Coin) bool {
	return coin.GetCurrency() == coin.GetNetworkCurrency()
}

func IsNetworkCurrency(currency meta.Currency) bool {
	coin, err := GetCoinNative(currency)
	if err != nil {
		return false
	}
	return IsNetworkCoin(coin)
}

func IsBitcoinSegWitAddress(address string) bool {
	oneIndex := strings.LastIndexByte(address, '1')
	if oneIndex < 2 {
		return false
	}
	prefix := address[:oneIndex+1]
	return chaincfg.IsBech32SegwitPrefix(prefix)
}

func ParseBitcoinWifAddress(
	wif *btcutil.WIF, chainConf *chaincfg.Params, hintAddress btcutil.Address,
) (address btcutil.Address, err error) {
	if hintAddress == nil {
		address, err = GetBitcoinWifAddressP2WPKH(wif, chainConf)
		if err != nil {
			err = utils.WrapError(err)
		}
		return
	}

	switch addressT := hintAddress.(type) {
	case *btcutil.AddressPubKey, *btcutil.AddressPubKeyHash:
		address, err = GetBitcoinWifAddressP2PKH(wif, chainConf)
		break
	case *btcutil.AddressWitnessPubKeyHash:
		address, err = GetBitcoinWifAddressP2WPKH(wif, chainConf)
		break
	default:
		err = utils.IssueErrorf(
			"bitcoin transaction signer doesn't support this type of address | address=%v,type=%v",
			hintAddress, reflect.TypeOf(addressT).Name,
		)
		break
	}
	if err != nil {
		return
	}
	if address.EncodeAddress() != hintAddress.EncodeAddress() {
		err = utils.IssueErrorf(
			"bitcoin transaction signer encounters a mismatched private key and address | expected_address=%v,actual_address=%v",
			address.EncodeAddress(), hintAddress.EncodeAddress(),
		)
		return
	}

	return address, nil
}

func GetTokenMetaEthereumTetherUSD() TokenMeta {
	if config.BlockchainUseTestnet {
		return EthereumTokenMetaTestRopstenTetherUSD
	} else {
		return EthereumTokenMetaMainnetTetherUSD
	}
}

func GetTokenMetaTronTetherUSD() TokenMeta {
	if config.BlockchainUseTestnet {
		return TronTokenMetaTestShastaTetherUSD
	} else {
		return TronTokenMetaMainnetTetherUSD
	}
}

func GetSystemNetworkEthereum() meta.BlockchainNetwork {
	if config.BlockchainUseTestnet {
		return constants.BlockchainNetworkEthereumTestRopsten
	} else {
		return constants.BlockchainNetworkEthereum
	}
}

func GetSystemNetworkTron() meta.BlockchainNetwork {
	if config.BlockchainUseTestnet {
		return constants.BlockchainNetworkTronTestShasta
	} else {
		return constants.BlockchainNetworkTron
	}
}
