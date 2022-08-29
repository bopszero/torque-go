package blockchainmod

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BitcoreEthereumBlock struct {
	BitcoreBalanceLikeBlock
}

type BitcoreEthereumTransaction struct {
	BitcoreBalanceLikeTransaction

	GasLimit int64                             `json:"gasLimit"`
	GasPrice int64                             `json:"gasPrice"`
	Receipt  BitcoreEthereumTransactionReceipt `json:"receipt"`
}

func (this *BitcoreEthereumTransaction) isConfirmed() bool {
	return this.GetConfirmations() >= EthereumMinConfirmations
}

func (this *BitcoreEthereumTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if !this.isConfirmed() {
		if this.isConflicted() {
			return constants.BlockchainTxnStatusFailed
		} else {
			return constants.BlockchainTxnStatusPending
		}
	}
	if this.Receipt.IsEmpty() {
		return constants.BlockchainTxnStatusPending
	} else if this.Receipt.IsSuccess {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusFailed
	}
}

func (this *BitcoreEthereumTransaction) GetDirection() meta.Direction {
	switch strings.ToLower(this.ownerAddress) {

	case strings.ToLower(this.FromAddress):
		return constants.DirectionTypeSend

	case strings.ToLower(this.ToAddress):
		return constants.DirectionTypeReceive

	default:
		if this.GetCurrency() != constants.CurrencyEthereum {
			return constants.DirectionTypeReceive
		}

		return constants.DirectionTypeUnknown
	}
}

func (this *BitcoreEthereumTransaction) GetFee() meta.CurrencyAmount {
	fee := meta.CurrencyAmount{
		Currency: this.currencyUnit,
	}
	if !this.Receipt.IsEmpty() {
		fee.Value = decimal.NewFromInt(this.GasPrice * this.Receipt.GasUsed)
	}
	return currencymod.ConvertAmountF(fee, this.currencyBase)
}

func (this *BitcoreEthereumTransaction) GetRC20Transfers(tokenMeta TokenMeta) ([]RC20Transfer, error) {
	var anyRC20Transfers []RC20Transfer
	for _, log := range this.Receipt.Logs {
		transfer, ok := log.GetERC20TokenTransfer(tokenMeta)
		if !ok {
			continue
		}
		anyRC20Transfers = append(anyRC20Transfers, &transfer)
	}
	return anyRC20Transfers, nil
}

type BitcoreEthereumTransactionReceipt struct {
	BlockHash   string                          `json:"blockHash"`
	BlockHeight int64                           `json:"blockNumber"`
	FromAddress string                          `json:"from"`
	ToAddress   string                          `json:"to"`
	GasUsed     int64                           `json:"gasUsed"`
	IsSuccess   bool                            `json:"status"`
	TxnHash     string                          `json:"transactionHash"`
	Logs        []BitcoreEthereumTransactionLog `json:"logs"`
	// LogsBloom   string `json:"logsBloom"`
}

func (this *BitcoreEthereumTransactionReceipt) IsEmpty() bool {
	return this.BlockHeight == 0
}

type BitcoreEthereumTransactionLog struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber uint64   `json:"blockNumber"`
	TxnHash     string   `json:"transactionHash"`
	TxnIndex    uint64   `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	LogIndex    int64    `json:"logIndex"`
	IsRemoved   bool     `json:"removed"`
}

func (this BitcoreEthereumTransactionLog) IsERC20TokenTransfer(tokenMeta TokenMeta) bool {
	if !utils.IsSameStringCI(this.Address, tokenMeta.Address) {
		return false
	}
	if len(this.Topics) != 3 {
		return false
	}
	if this.Topics[0] != EthereumERC20TrasnferTokenLogTopicHex {
		return false
	}

	amountBytes, err := comutils.HexDecode0x(this.Data)
	if err != nil {
		return false
	}
	amount := decimal.NewFromBigInt(
		new(big.Int).SetBytes(amountBytes),
		-int32(tokenMeta.DecimalPlaces))
	if amount.IsNegative() {
		return false
	}

	return true
}

func (this BitcoreEthereumTransactionLog) GetERC20TokenTransfer(tokenMeta TokenMeta) (
	_ BitcoreEthereumTransactionLogERC20Transfer, _ bool,
) {
	if !this.IsERC20TokenTransfer(tokenMeta) {
		return
	}

	amountBytes, err := comutils.HexDecode0x(this.Data)
	comutils.PanicOnError(err)
	amount := decimal.NewFromBigInt(
		new(big.Int).SetBytes(amountBytes),
		-int32(tokenMeta.DecimalPlaces))
	fromAddressBytes, err := comutils.HexDecode0x(this.Topics[1])
	if err != nil {
		return
	}
	toAddressBytes, err := comutils.HexDecode0x(this.Topics[2])
	if err != nil {
		return
	}

	var (
		fromAddress = common.BytesToAddress(fromAddressBytes)
		toAddress   = common.BytesToAddress(toAddressBytes)
	)
	transfer := BitcoreEthereumTransactionLogERC20Transfer{
		TokenMeta:   tokenMeta,
		MethodHash:  this.Topics[0],
		FromAddress: fromAddress.String(),
		ToAddress:   toAddress.String(),
		Amount:      amount,
	}
	return transfer, true
}

type BitcoreEthereumTransactionLogERC20Transfer struct {
	TokenMeta   TokenMeta
	MethodHash  string
	FromAddress string
	ToAddress   string
	Amount      decimal.Decimal
}

func (this *BitcoreEthereumTransactionLogERC20Transfer) GetTokenMeta() TokenMeta {
	return this.TokenMeta
}

func (this *BitcoreEthereumTransactionLogERC20Transfer) GetFromAddress() string {
	return this.FromAddress
}

func (this *BitcoreEthereumTransactionLogERC20Transfer) GetToAddress() string {
	return this.ToAddress
}

func (this *BitcoreEthereumTransactionLogERC20Transfer) GetAmount() decimal.Decimal {
	return this.Amount
}
