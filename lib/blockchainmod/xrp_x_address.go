package blockchainmod

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/rubblelabs/ripple/crypto"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type RippleXAddress struct {
	addressHash crypto.Hash
	tag         *uint32
	isMainnet   bool

	xAddress string
}

func NewRippleXAddress(address string, tag *uint32, forMainnet bool) (xAddress RippleXAddress, err error) {
	addressHash, err := crypto.NewRippleHashCheck(address, crypto.RIPPLE_ACCOUNT_ID)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	var (
		tagValue uint32
		tagRef   *uint32
	)
	if tag != nil {
		tagValue = *tag
		tagRef = &tagValue
	}

	xAddress = RippleXAddress{
		addressHash: addressHash,
		tag:         tagRef,
		isMainnet:   forMainnet,
	}
	return xAddress, nil
}

func LoadRippleXAddress(xAddressStr string) (xAddress RippleXAddress, err error) {
	err = xAddress.decode(xAddressStr)
	return
}

func (this RippleXAddress) GetRootAddress() string {
	if this.addressHash == nil {
		return ""
	} else {
		return this.addressHash.String()
	}
}

func (this RippleXAddress) GetRootTagAddress() string {
	if this.addressHash == nil {
		return ""
	}
	if this.tag == nil {
		return this.addressHash.String()
	} else {
		return fmt.Sprintf("%s:%d", this.addressHash.String(), *this.tag)
	}
}

func (this RippleXAddress) GetTag() (tagTag bool, tag uint32) {
	if this.tag == nil {
		return false, 0
	} else {
		return true, *this.tag
	}
}

func (this RippleXAddress) IsMainnet() bool {
	return this.isMainnet
}

func (this RippleXAddress) encode() string {
	var (
		prefix   []byte
		tagFlag  byte
		tagBytes = make([]byte, 8) // four zero bytes (reserved for 64-bit tags)
	)
	if this.isMainnet {
		prefix = RippleXAddressPrefixMainnet
	} else {
		prefix = RippleXAddressPrefixTestnet
	}
	if this.tag != nil {
		tagFlag = 1
		if *this.tag == 0 {
			randomTag := rand.Uint32()
			this.tag = &randomTag
		}
		binary.LittleEndian.PutUint32(tagBytes, *this.tag)
	} else {
		tagFlag = 0
	}

	addressBytes := make([]byte, RippleXAddressSize)
	copy(addressBytes[:2], prefix)
	copy(addressBytes[2:22], this.addressHash.Payload())
	addressBytes[22] = tagFlag
	copy(addressBytes[23:], tagBytes)

	return RippleBase58Encode(addressBytes)
}

func (this *RippleXAddress) decode(xAddress string) error {
	decodedBytes, err := RippleBase58Decode(xAddress)
	if err != nil {
		return utils.WrapError(err)
	}
	addressBytes := utils.BytesSlice(decodedBytes, 0, -4)
	if len(addressBytes) != RippleXAddressSize {
		return utils.WrapError(constants.ErrorAddress)
	}
	var (
		isMainnet   bool
		prefixBytes = addressBytes[:2]
	)
	switch {
	case bytes.Compare(prefixBytes, RippleXAddressPrefixMainnet) == 0:
		isMainnet = true
		break
	case bytes.Compare(prefixBytes, RippleXAddressPrefixTestnet) == 0:
		isMainnet = true
	default:
		return utils.WrapError(constants.ErrorAddress)
	}

	addressHash, err := crypto.NewAccountId(addressBytes[2:22])
	if err != nil {
		return utils.WrapError(err)
	}
	this.addressHash = addressHash

	switch addressBytes[22] {
	case 0:
		break
	case 1:
		tag := binary.LittleEndian.Uint32(addressBytes[23:])
		this.tag = &tag
	default:
		return utils.WrapError(constants.ErrorAddress)
	}

	this.isMainnet = isMainnet
	this.xAddress = RippleBase58Encode(addressBytes)

	return nil

}

func (this RippleXAddress) GetAddress() string {
	if this.xAddress == "" {
		this.xAddress = this.encode()
	}
	return this.xAddress
}

func (this RippleXAddress) String() string {
	return this.GetAddress()
}
