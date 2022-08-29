package blockchainmod

import (
	"github.com/gcash/bchd/bchec"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func CreateBitcoinCashPrivateKey() (*bchec.PrivateKey, error) {
	return bchec.NewPrivateKey(bchec.S256())
}

func GetBitcoinCashWifAddressP2PKH(wif *bchutil.WIF, chainConfig *chaincfg.Params) (
	*bchutil.AddressPubKeyHash, error,
) {
	wifPubKeyBytes := wif.PrivKey.PubKey().SerializeCompressed()
	wifPkHash := bchutil.Hash160(wifPubKeyBytes)
	address, err := bchutil.NewAddressPubKeyHash(wifPkHash, chainConfig)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return address, nil
}
