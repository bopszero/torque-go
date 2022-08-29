package blockchainmod

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func CreateBitcoinPrivateKey() (*btcec.PrivateKey, error) {
	key, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return key, nil
}

func LoadBitcoinPrivateKey(data []byte) *btcec.PrivateKey {
	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), data)
	return pk
}

func GetBitcoinWifAddressP2PKH(wif *btcutil.WIF, chainConfig *chaincfg.Params) (
	*btcutil.AddressPubKeyHash, error,
) {
	wifPubKeyBytes := wif.PrivKey.PubKey().SerializeCompressed()
	wifPkHash := btcutil.Hash160(wifPubKeyBytes)
	address, err := btcutil.NewAddressPubKeyHash(wifPkHash, chainConfig)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return address, nil
}

func GetBitcoinWifAddressP2WPKH(wif *btcutil.WIF, chainConfig *chaincfg.Params) (
	*btcutil.AddressWitnessPubKeyHash, error,
) {
	wifPubKeyBytes := wif.PrivKey.PubKey().SerializeCompressed()
	witnessProgram := btcutil.Hash160(wifPubKeyBytes)
	address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, chainConfig)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return address, nil
}

func newBitcoinLikeSegWitKeyHolder(
	wif *btcutil.WIF, chainConfig *chaincfg.Params, hintAddress btcutil.Address,
) (_ KeyHolder, err error) {
	var address btcutil.Address
	if hintAddress != nil {
		address, err = ParseBitcoinWifAddress(wif, chainConfig, hintAddress)
	} else {
		address, err = GetBitcoinWifAddressP2WPKH(wif, chainConfig)
	}
	if err != nil {
		return nil, err
	}

	account := BitcoinLikeKeyHolder{
		wif:     wif,
		address: address,
	}
	return &account, nil
}

func GenBitcoinLikeWIF(chainConfig *chaincfg.Params) (*btcutil.WIF, error) {
	privateKey, err := CreateBitcoinPrivateKey()
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(privateKey, chainConfig, true)
}

func genBitcoinLikeSegwitKeyHolder(chainConfig *chaincfg.Params) (KeyHolder, error) {
	wif, err := GenBitcoinLikeWIF(chainConfig)
	if err != nil {
		return nil, err
	}

	return newBitcoinLikeSegWitKeyHolder(wif, chainConfig, nil)
}
