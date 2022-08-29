package blockchainmod

import (
	"encoding/binary"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func BitcoinToSatoshi(amount decimal.Decimal) int64 {
	satoshiAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyBitcoin,
			Value:    amount,
		},
		constants.CurrencySubBitcoinSatoshi,
	)
	return satoshiAmount.Value.IntPart()
}

func SatoshiToBitcoin(amount int64) decimal.Decimal {
	btcAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubBitcoinSatoshi,
			Value:    decimal.NewFromInt(amount),
		},
		constants.CurrencyBitcoin,
	)
	return btcAmount.Value
}

func EthereumToWei(amount decimal.Decimal) *big.Int {
	weiAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencyEthereum,
			Value:    amount,
		},
		constants.CurrencySubEthereumWei,
	)
	return weiAmount.Value.BigInt()
}

func WeiToEthereum(value *big.Int) decimal.Decimal {
	ethAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    decimal.NewFromBigInt(value, 0),
		},
		constants.CurrencyEthereum,
	)
	return ethAmount.Value
}

func GweiToWei(amount decimal.Decimal) *big.Int {
	weiAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumGwei,
			Value:    amount,
		},
		constants.CurrencySubEthereumWei,
	)
	return weiAmount.Value.BigInt()
}

func hex0xToUint64(hexString string) (uint64, error) {
	realHex := comutils.HexTrim(hexString)
	if len(realHex)%2 != 0 {
		realHex = "0" + realHex
	}

	valueBytes, err := comutils.HexDecode(realHex)
	if err != nil {
		return 0, err
	}
	valueBytes = common.LeftPadBytes(valueBytes, 8)

	return binary.BigEndian.Uint64(valueBytes), nil
}

func uint64ToHex0x(number uint64) string {
	return constants.HexHumanPrefix + strconv.FormatUint(number, 16)
}

func hexToEthereumAddress(addressHex string) (addr common.Address, err error) {
	if !common.IsHexAddress(addressHex) {
		err = utils.WrapError(constants.ErrorAddress)
		return
	}

	addr = common.HexToAddress(addressHex)
	return
}

func TimeMsToS(timeMs int64) int64 {
	return timeMs / 1000
}

func TimeNsToMs(timeNs int64) int64 {
	return timeNs / 1000000
}
