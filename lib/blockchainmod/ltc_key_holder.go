package blockchainmod

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GenLitecoinMainnetKeyHolder() (KeyHolder, error) {
	return genLitecoinKeyHolder(&LitecoinChainConfig)
}

func GenLitecoinTestnetKeyHolder() (KeyHolder, error) {
	return genLitecoinKeyHolder(&LitecoinChainConfigTestnet)
}

func LoadLitecoinMainnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	return loadLitecoinKeyHolder(wifKey, &LitecoinChainConfig)
}

func LoadLitecoinTestnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	return loadLitecoinKeyHolder(wifKey, &LitecoinChainConfigTestnet)
}

func genLitecoinKeyHolder(chainConfig *chaincfg.Params) (KeyHolder, error) {
	wif, err := GenBitcoinLikeWIF(chainConfig)
	if err != nil {
		return nil, err
	}

	return newLitecoinKeyHolder(wif, chainConfig)
}

func loadLitecoinKeyHolder(wifKey string, chainConf *chaincfg.Params) (_ KeyHolder, err error) {
	wif, err := btcutil.DecodeWIF(wifKey)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return newLitecoinKeyHolder(wif, chainConf)
}

func newLitecoinKeyHolder(wif *btcutil.WIF, chainConfig *chaincfg.Params) (
	_ KeyHolder, err error,
) {
	address, err := GetBitcoinWifAddressP2PKH(wif, chainConfig)
	if err != nil {
		return nil, err
	}

	account := BitcoinLikeKeyHolder{
		wif:     wif,
		address: address,
	}
	return &account, nil
}
