package blockchainmod

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

type TronKeyHolder struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (this *TronKeyHolder) GetPrivateKey() string {
	return comutils.HexEncode(crypto.FromECDSA(this.privateKey))
}

func (this *TronKeyHolder) GetPublicKey() string {
	return comutils.HexEncode(crypto.FromECDSAPub(this.publicKey))
}

func (this *TronKeyHolder) GetAddress() string {
	return address.PubkeyToAddress(*this.publicKey).String()
}

func GenTronKeyHolder() (KeyHolder, error) {
	privateKey, err := CreateEcdsaPrivateKey()
	if err != nil {
		return nil, err
	}

	return newTronKeyHolder(privateKey)
}

func LoadTronKeyHolder(privateKey string) (KeyHolder, error) {
	privateKeyEcdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return newTronKeyHolder(privateKeyEcdsa)
}

func newTronKeyHolder(privateKey *ecdsa.PrivateKey) (KeyHolder, error) {
	account := TronKeyHolder{
		privateKey: privateKey,
		publicKey:  privateKey.Public().(*ecdsa.PublicKey),
	}
	return &account, nil
}
