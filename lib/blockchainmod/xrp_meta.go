package blockchainmod

import "github.com/shopspring/decimal"

const (
	RippleDecimalPlaces    = 6
	RippleSeedSize         = 16
	RippleNormalTxnFee     = 12
	RippleMinConfirmations = 12
	RippleTxnStatusSuccess = "success"
)

const (
	RippleKeyEcdsaDefaultSequence = 0

	RippleXAddressSize          = 2 + 20 + 1 + 8 // Prefix + Address + UseTag + Tag
	RippleXAddressTagMin        = 10000000
	RippleXAddressTagMax        = 100000000 - 1
	RippleXAddressTagTotalCount = RippleXAddressTagMax - RippleXAddressTagMin + 1
)

var (
	RippleReserveBalance = decimal.NewFromInt(20)
)

var (
	RippleSeedEd25519Prefix     = []byte{0x01, 0xE1, 0x4B}
	RippleXAddressPrefixMainnet = []byte{0x05, 0x44}
	RippleXAddressPrefixTestnet = []byte{0x04, 0x93}
)
