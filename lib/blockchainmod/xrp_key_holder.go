package blockchainmod

import (
	"github.com/rubblelabs/ripple/crypto"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

type RippleKeyHolder struct {
	b58Seed   string
	sequence  *uint32
	key       crypto.Key
	isMainnet bool
	tag       *uint32
}

func (this *RippleKeyHolder) SetTag(tag uint32) {
	this.tag = &tag
}

func (this *RippleKeyHolder) GetKey() crypto.Key {
	return this.key
}

func (this *RippleKeyHolder) GetPrivateKey() string {
	return this.b58Seed
}

func (this *RippleKeyHolder) GetPublicKey() string {
	return comutils.HexEncode(this.key.Public(this.sequence))
}

func (this *RippleKeyHolder) GetLegacyAddress() string {
	hash, err := crypto.AccountId(this.key, this.sequence)
	comutils.PanicOnError(err)
	return hash.String()
}

func (this *RippleKeyHolder) GetAddress() string {
	xAddress, err := NewRippleXAddress(this.GetLegacyAddress(), this.tag, this.isMainnet)
	comutils.PanicOnError(err)
	return xAddress.String()
}

func genRippleKeyHolder(forMainnet bool) (*RippleKeyHolder, error) {
	newSeed := genRippleEd25519Seed()
	keyHolder, err := newRippleKeyHolder(newSeed, forMainnet)
	if err != nil {
		return nil, err
	}

	return keyHolder, nil
}

func GenRippleMainnetKeyHolder() (*RippleKeyHolder, error) {
	return genRippleKeyHolder(true)
}

func GenRippleTestnetKeyHolder() (*RippleKeyHolder, error) {
	return genRippleKeyHolder(false)
}

func loadRippleKeyHolder(b58Seed string, hintAddress string, forMainnet bool) (*RippleKeyHolder, error) {
	keyHolder, err := newRippleKeyHolder(b58Seed, forMainnet)
	if err != nil {
		return nil, err
	}
	if hintAddress == "" {
		return keyHolder, nil
	}

	xAddress, err := RippleParseXAddress(hintAddress, forMainnet)
	if err != nil {
		return nil, err
	}
	if hasTag, addressTag := xAddress.GetTag(); hasTag {
		keyHolder.SetTag(addressTag)
	}
	return keyHolder, nil
}

func LoadRippleMainnetKeyHolder(b58Seed string, hintAddress string) (*RippleKeyHolder, error) {
	return loadRippleKeyHolder(b58Seed, hintAddress, true)
}

func LoadRippleTestnetKeyHolder(b58Seed string, hintAddress string) (*RippleKeyHolder, error) {
	return loadRippleKeyHolder(b58Seed, hintAddress, false)
}

func newRippleKeyHolder(b58Seed string, forMainnet bool) (*RippleKeyHolder, error) {
	seedKey, seq, err := LoadRippleSeedKey(b58Seed)
	if err != nil {
		return nil, err
	}

	account := RippleKeyHolder{
		b58Seed:   b58Seed,
		sequence:  seq,
		key:       seedKey,
		isMainnet: forMainnet,
	}
	return &account, nil
}
