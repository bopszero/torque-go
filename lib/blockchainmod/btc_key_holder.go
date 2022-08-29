package blockchainmod

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoinLikeKeyHolder struct {
	wif     *btcutil.WIF
	address btcutil.Address
}

func (this *BitcoinLikeKeyHolder) GetPrivateKey() string {
	return this.wif.String()
}

func (this *BitcoinLikeKeyHolder) GetPublicKey() string {
	return base58.Encode(this.wif.PrivKey.PubKey().SerializeCompressed())
}

func (this *BitcoinLikeKeyHolder) GetAddress() string {
	return this.address.String()
}

func GenBitcoinMainnetKeyHolder() (KeyHolder, error) {
	return genBitcoinLikeSegwitKeyHolder(&BitcoinChainConfig)
}

func GenBitcoinTestnetKeyHolder() (KeyHolder, error) {
	return genBitcoinLikeSegwitKeyHolder(&BitcoinChainConfigTestnet)
}

func LoadBitcoinMainnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	return loadBitcoinKeyHolder(wifKey, &BitcoinChainConfig, hintAddress)
}

func LoadBitcoinTestnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	return loadBitcoinKeyHolder(wifKey, &BitcoinChainConfigTestnet, hintAddress)
}

func loadBitcoinKeyHolder(
	wifKey string, chainConf *chaincfg.Params, hintAddress string,
) (_ KeyHolder, err error) {
	wif, err := btcutil.DecodeWIF(wifKey)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	var parsedHintAddress btcutil.Address
	if hintAddress != "" {
		parsedHintAddress, err = btcutil.DecodeAddress(hintAddress, chainConf)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	}

	return newBitcoinLikeSegWitKeyHolder(wif, chainConf, parsedHintAddress)
}
