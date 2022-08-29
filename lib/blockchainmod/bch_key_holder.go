package blockchainmod

import (
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	"github.com/gcash/bchutil/base58"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoinCashKeyHolder struct {
	wif     *bchutil.WIF
	address bchutil.Address
}

func (this *BitcoinCashKeyHolder) GetPrivateKey() string {
	return this.wif.String()
}

func (this *BitcoinCashKeyHolder) GetPublicKey() string {
	return base58.Encode(this.wif.PrivKey.PubKey().SerializeCompressed())
}

func (this *BitcoinCashKeyHolder) GetAddress() string {
	return this.address.EncodeAddress()
}

func GenBitcoinCashMainnetKeyHolder() (KeyHolder, error) {
	return genBitcoinCashKeyHolder(&BitcoinCashChainConfig)
}

func GenBitcoinCashTestnetKeyHolder() (KeyHolder, error) {
	return genBitcoinCashKeyHolder(&BitcoinCashChainConfigTestnet)
}

func genBitcoinCashWIF(chainConfig *chaincfg.Params) (*bchutil.WIF, error) {
	privateKey, err := CreateBitcoinCashPrivateKey()
	if err != nil {
		return nil, err
	}

	return bchutil.NewWIF(privateKey, chainConfig, true)
}

func genBitcoinCashKeyHolder(chainConfig *chaincfg.Params) (KeyHolder, error) {
	wif, err := genBitcoinCashWIF(chainConfig)
	if err != nil {
		return nil, err
	}

	return newBitcoinCashKeyHolder(wif, chainConfig)
}

func LoadBitcoinCashMainnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	wif, err := bchutil.DecodeWIF(wifKey)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return newBitcoinCashKeyHolder(wif, &BitcoinCashChainConfig)
}

func LoadBitcoinCashTestnetKeyHolder(wifKey string, hintAddress string) (KeyHolder, error) {
	wif, err := bchutil.DecodeWIF(wifKey)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return newBitcoinCashKeyHolder(wif, &BitcoinCashChainConfigTestnet)
}

func newBitcoinCashKeyHolder(wif *bchutil.WIF, chainConfig *chaincfg.Params) (
	_ KeyHolder, err error,
) {
	address, err := GetBitcoinCashWifAddressP2PKH(wif, chainConfig)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	account := BitcoinCashKeyHolder{
		wif:     wif,
		address: address,
	}
	return &account, nil
}
