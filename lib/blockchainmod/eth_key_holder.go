package blockchainmod

import (
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

type EthereumKeyHolder struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (this *EthereumKeyHolder) GetPrivateKey() string {
	return comutils.HexEncode0x(crypto.FromECDSA(this.privateKey))
}

func (this *EthereumKeyHolder) GetPublicKey() string {
	return comutils.HexEncode0x(crypto.FromECDSAPub(this.publicKey))
}

func (this *EthereumKeyHolder) GetAddress() string {
	return strings.ToLower(crypto.PubkeyToAddress(*this.publicKey).Hex())
}

func GenEthereumKeyHolder() (KeyHolder, error) {
	privateKey, err := CreateEcdsaPrivateKey()
	if err != nil {
		return nil, err
	}

	return newEthereumKeyHolder(privateKey)
}

func LoadEthereumKeyHolder(privateKey string) (KeyHolder, error) {
	privateKeyEcdsa, err := crypto.HexToECDSA(comutils.HexTrim(privateKey))
	if err != nil {
		return nil, err
	}

	return newEthereumKeyHolder(privateKeyEcdsa)
}

func newEthereumKeyHolder(privateKey *ecdsa.PrivateKey) (KeyHolder, error) {
	account := EthereumKeyHolder{
		privateKey: privateKey,
		publicKey:  privateKey.Public().(*ecdsa.PublicKey),
	}
	return &account, nil
}
