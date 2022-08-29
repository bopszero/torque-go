package blockchainmod

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"

	"github.com/rubblelabs/ripple/crypto"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func RippleBase58Encode(bytes []byte) string {
	return crypto.Base58Encode(bytes, crypto.ALPHABET)
}

func RippleBase58Decode(text string) ([]byte, error) {
	textBytes, err := crypto.Base58Decode(text, crypto.ALPHABET)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return textBytes, nil
}

func genRippleEd25519Seed() string {
	var (
		randReader = rand.Reader
		dataBytes  = make([]byte, RippleSeedSize)
	)
	_, err := io.ReadFull(randReader, dataBytes)
	comutils.PanicOnError(err)

	var (
		fullSize  = len(RippleSeedEd25519Prefix) + RippleSeedSize
		seedBytes = make([]byte, 0, fullSize)
	)
	seedBytes = append(seedBytes, RippleSeedEd25519Prefix...)
	seedBytes = append(seedBytes, dataBytes...)

	return RippleBase58Encode(seedBytes)
}

func LoadRippleSeedKey(seed string) (crypto.Key, *uint32, error) {
	hash, err := crypto.NewRippleHash(seed)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	if hash.Version() == crypto.RIPPLE_FAMILY_SEED {
		var sequence uint32 = RippleKeyEcdsaDefaultSequence
		key, err := crypto.NewECDSAKey(hash.Payload())
		return key, &sequence, err
	} else {
		key, err := loadRippleEd25519SeedKey(seed)
		return key, nil, err
	}
}

func loadRippleEd25519SeedKey(seed string) (crypto.Key, error) {
	seedDecodedBytes, err := RippleBase58Decode(seed)
	if err != nil {
		return nil, err
	}
	var (
		seedBytes  = utils.BytesSlice(seedDecodedBytes, 0, -4)
		seedPrefix = seedBytes[:len(RippleSeedEd25519Prefix)]
	)
	if bytes.Compare(seedPrefix, RippleSeedEd25519Prefix) != 0 {
		return nil, utils.IssueErrorf(
			"Ripple ed25519 seed has a invalid prefix %v",
			comutils.HexEncode(seedPrefix),
		)
	}
	seedBody := seedBytes[len(seedPrefix):]
	if len(seedBody) != RippleSeedSize {
		return nil, utils.IssueErrorf(
			"Ripple ed25519 seed has a invalid length %v",
			len(seedBody),
		)
	}
	return crypto.NewEd25519Key(seedBody)
}

func RippleParseXAddress(inputAddress string, forMainnet bool) (xAddress RippleXAddress, err error) {
	addressParts := strings.Split(inputAddress, ":")
	switch len(addressParts) {
	case 1:
		address := addressParts[0]
		if _, err = crypto.NewRippleHashCheck(address, crypto.RIPPLE_ACCOUNT_ID); err == nil {
			xAddress, err = NewRippleXAddress(address, nil, forMainnet)
		} else {
			xAddress, err = LoadRippleXAddress(address)
		}
	case 2:
		var (
			address = addressParts[0]
			tagStr  = addressParts[1]
		)
		if tagStr == "" {
			xAddress, err = NewRippleXAddress(address, nil, forMainnet)
		} else {
			tag, tagErr := comutils.ParseUint32(tagStr)
			if tagErr != nil {
				err = utils.WrapError(tagErr)
				return
			}
			xAddress, err = NewRippleXAddress(address, &tag, forMainnet)
		}
	default:
		err = utils.WrapError(constants.ErrorAddress)
	}
	if err != nil {
		err = utils.IssueErrorf("malformed XRP address | address=%s,err=%s", inputAddress, err.Error())
	}
	return
}

func rippleToDrop(value decimal.Decimal) int64 {
	amount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyRipple,
			Value:    value,
		},
		constants.CurrencySubRippleDrop,
	)
	return amount.Value.IntPart()
}
