package blockchainmod

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type EtherscanEthereumBlock struct {
	Hash          string                            `json:"hash"`
	HeightHex     string                            `json:"number"`
	ParentHash    string                            `json:"parentHash"`
	CreateTimeHex string                            `json:"timestamp"`
	Transactions  []EtherscanEthereumHexTransaction `json:"transactions"`
}

func (this *EtherscanEthereumBlock) GetCurrency() meta.Currency {
	return constants.CurrencyEthereum
}

func (this *EtherscanEthereumBlock) GetNetwork() meta.BlockchainNetwork {
	return constants.BlockchainNetworkEthereum
}

func (this *EtherscanEthereumBlock) GetTimeUnix() int64 {
	value, err := hex0xToUint64(this.CreateTimeHex)
	comutils.PanicOnError(err)

	return int64(value)
}

func (this *EtherscanEthereumBlock) GetHash() string {
	return this.Hash
}

func (this *EtherscanEthereumBlock) GetHeight() uint64 {
	height, err := hex0xToUint64(this.HeightHex)
	comutils.PanicOnError(err)

	return height
}

func (this *EtherscanEthereumBlock) GetTransactions() ([]Transaction, error) {
	txns := make([]Transaction, len(this.Transactions))

	for i := range this.Transactions {
		txn := &this.Transactions[i]
		txn.createTime = this.GetTimeUnix()
		txns[i] = txn
	}

	return txns, nil
}

func (this *EtherscanEthereumBlock) GetParentHash() string {
	return this.ParentHash
}

type EtherscanEthereumHexTransaction struct {
	baseTransaction

	createTime int64
	client     *EtherscanEthereumClient
	receipt    *etherscanEthereumTxnReceipt

	Hash           string `json:"hash"`
	BlockHash      string `json:"blockHash"`
	BlockHeightHex string `json:"blockNumber"`
	NonceHex       string `json:"nonce"`

	FromAddress  string `json:"from"`
	ToAddress    string `json:"to"`
	AmountWeiHex string `json:"value"`
	InputHex     string `json:"input"`

	GasPriceWeiHex string `json:"gasPrice"`
	GasLimitHex    string `json:"gas"`
}

func (this *EtherscanEthereumHexTransaction) getReceipt() (*etherscanEthereumTxnReceipt, error) {
	if this.receipt == nil {
		if this.client == nil {
			return nil, utils.IssueErrorf("cannot get txn fee on etherscan without client")
		}

		receipt, err := this.client.GetTxnReceipt(this.Hash)
		if err != nil {
			return nil, err
		}

		this.receipt = receipt
	}

	return this.receipt, nil
}

func (this *EtherscanEthereumHexTransaction) getLastestBlockHeightFast() (latestBlockHeight uint64) {
	cacheKey := fmt.Sprintf("blockchain:client:%v:block:latest", this.GetCurrency())

	err := comcache.GetOrCreate(
		comcache.GetDefaultCache(),
		cacheKey,
		15*time.Second,
		&latestBlockHeight,
		func() (interface{}, error) {
			return this.client.GetLatestBlockHeight()
		},
	)
	comutils.PanicOnError(err)

	return
}

func (this *EtherscanEthereumHexTransaction) SetClient(client *EtherscanEthereumClient) {
	this.client = client
}

func (this *EtherscanEthereumHexTransaction) GetCurrency() meta.Currency {
	return constants.CurrencyEthereum
}

func (this *EtherscanEthereumHexTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetBlockHash() == "" {
		return constants.BlockchainTxnStatusPending
	}

	receipt, err := this.getReceipt()
	comutils.PanicOnError(err)
	if !receipt.IsEmpty() && !receipt.IsSuccess() {
		return constants.BlockchainTxnStatusFailed
	}

	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *EtherscanEthereumHexTransaction) GetDirection() meta.Direction {
	switch this.ownerAddress {
	case this.FromAddress:
		return constants.DirectionTypeSend
	case this.ToAddress:
		return constants.DirectionTypeReceive
	default:
		return constants.DirectionTypeUnknown
	}
}

func (this *EtherscanEthereumHexTransaction) GetFee() meta.CurrencyAmount {
	receipt, err := this.getReceipt()
	comutils.PanicOnError(err)
	if receipt.IsEmpty() {
		return meta.CurrencyAmount{
			Currency: constants.CurrencyEthereum,
			Value:    decimal.Zero,
		}
	}

	gasUsed := receipt.GetGasUsed()
	gasPriceWei, err := hex0xToUint64(this.GasPriceWeiHex)
	comutils.PanicOnError(err)

	gasTotalWei := gasPriceWei * gasUsed
	return currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    comutils.NewDecimalF(gasTotalWei),
		},
		constants.CurrencyEthereum,
	)
}

func (this *EtherscanEthereumHexTransaction) GetHash() string {
	return this.Hash
}

func (this *EtherscanEthereumHexTransaction) GetBlockHash() string {
	return this.BlockHash
}

func (this *EtherscanEthereumHexTransaction) GetBlockHeight() uint64 {
	height, err := hex0xToUint64(this.BlockHeightHex)
	comutils.PanicOnError(err)

	return height
}

func (this *EtherscanEthereumHexTransaction) GetConfirmations() uint64 {
	blockHeight := this.GetBlockHeight()
	if blockHeight == 0 {
		return 0
	}

	latestBlockHeight := this.getLastestBlockHeightFast()
	if latestBlockHeight > blockHeight {
		return latestBlockHeight - blockHeight + 1
	}

	return 1
}

func (this *EtherscanEthereumHexTransaction) GetFromAddress() (string, error) {
	return this.FromAddress, nil
}

func (this *EtherscanEthereumHexTransaction) GetToAddress() (string, error) {
	return this.ToAddress, nil
}

func (this *EtherscanEthereumHexTransaction) GetAmount() (decimal.Decimal, error) {
	amount, err := hex0xToUint64(this.AmountWeiHex)
	if err != nil {
		return decimal.Zero, err
	}

	amountEth := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    comutils.NewDecimalF(amount),
		},
		constants.CurrencyEthereum,
	)
	return amountEth.Value, nil
}

func (this *EtherscanEthereumHexTransaction) GetInputDataHex() (string, error) {
	return this.InputHex, nil
}

func (this *EtherscanEthereumHexTransaction) GetTimeUnix() int64 {
	return this.createTime
}

func (this *EtherscanEthereumHexTransaction) GetInputs() ([]Input, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction inputs")
}

func (this *EtherscanEthereumHexTransaction) GetOutputs() ([]Output, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction outputs")
}

type EtherscanEthereumTransaction struct {
	baseTransaction

	Hash        string `json:"hash"`
	BlockHash   string `json:"blockHash"`
	BlockHeight uint64 `json:"blockNumber,string"`
	Nonce       uint64 `json:"nonce,string"`

	FromAddress   string          `json:"from"`
	ToAddress     string          `json:"to"`
	AmountWei     decimal.Decimal `json:"value"`
	InputHex      string          `json:"input"`
	Confirmations uint64          `json:"confirmations,string"`
	Time          int64           `json:"timeStamp,string"`

	GasPriceWei uint64 `json:"gasPrice,string"`
	GasLimit    uint32 `json:"gas,string"`
	GasUsed     uint32 `json:"gasUsed,string"`
}

func (this *EtherscanEthereumTransaction) GetCurrency() meta.Currency {
	return constants.CurrencyEthereum
}

func (this *EtherscanEthereumTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *EtherscanEthereumTransaction) GetDirection() meta.Direction {
	switch this.ownerAddress {
	case this.FromAddress:
		return constants.DirectionTypeSend
	case this.ToAddress:
		return constants.DirectionTypeReceive
	default:
		return constants.DirectionTypeUnknown
	}
}

func (this *EtherscanEthereumTransaction) GetFee() meta.CurrencyAmount {
	gasTotalWei := this.GasPriceWei * uint64(this.GasUsed)
	return currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    comutils.NewDecimalF(gasTotalWei),
		},
		constants.CurrencyEthereum,
	)
}

func (this *EtherscanEthereumTransaction) GetHash() string {
	return this.Hash
}

func (this *EtherscanEthereumTransaction) GetBlockHash() string {
	return this.BlockHash
}

func (this *EtherscanEthereumTransaction) GetBlockHeight() uint64 {
	return this.BlockHeight
}

func (this *EtherscanEthereumTransaction) GetConfirmations() uint64 {
	return this.Confirmations
}

func (this *EtherscanEthereumTransaction) GetFromAddress() (string, error) {
	return this.FromAddress, nil
}

func (this *EtherscanEthereumTransaction) GetToAddress() (string, error) {
	return this.ToAddress, nil
}

func (this *EtherscanEthereumTransaction) GetAmount() (decimal.Decimal, error) {
	amountEth := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    this.AmountWei,
		},
		constants.CurrencyEthereum,
	)
	return amountEth.Value, nil
}

func (this *EtherscanEthereumTransaction) GetInputDataHex() (string, error) {
	return this.InputHex, nil
}

func (this *EtherscanEthereumTransaction) GetTimeUnix() int64 {
	return this.Time
}

func (this *EtherscanEthereumTransaction) GetInputs() ([]Input, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction inputs")
}

func (this *EtherscanEthereumTransaction) GetOutputs() ([]Output, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction outputs")
}

type etherscanEthereumTxnReceipt struct {
	BlockHash       string `json:"blockHash"`
	BlockHeightHex  string `json:"blockNumber"`
	ContractAddress string `json:"contractAddress"`
	GasUsedHex      string `json:"gasUsed"`
	StatusHex       string `json:"status"`
}

func (this *etherscanEthereumTxnReceipt) IsEmpty() bool {
	return this.BlockHash == ""
}

func (this *etherscanEthereumTxnReceipt) GetGasUsed() uint64 {
	gasUsed, err := hex0xToUint64(this.GasUsedHex)
	comutils.PanicOnError(err)

	return gasUsed
}

func (this *etherscanEthereumTxnReceipt) IsSuccess() bool {
	value, err := hex0xToUint64(this.StatusHex)
	comutils.PanicOnError(err)

	return value == 1
}

// type EtherscanEthereumInternalTransaction struct {
// 	ToAddress ethCommon.Address
// 	Amount    decimal.Decimal
// }

// func NewEtherscanEthereumInternalTransaction(contractDataHex string) (isValid bool, _ *EtherscanEthereumInternalTransaction, err error) {
// 	contactDataBytes, err := comutils.HexDecode(comutils.HexTrim(contractDataHex))
// 	if err != nil {
// 		return
// 	}
// 	dataBuffer := bytes.NewBuffer(contactDataBytes)

// 	methodID := make([]byte, 4)
// 	_, err = dataBuffer.Read(methodID)
// 	if err != nil {
// 		return
// 	}
// 	methodIdHex := hexHumanEncode(methodID)

// 	isValid = methodIdHex == EtherscanEthereumMethodIdInternalTxn
// 	if !isValid {
// 		return
// 	}

// 	addressBytes := make([]byte, 32)
// 	_, err = dataBuffer.Read(addressBytes)
// 	if err != nil {
// 		return
// 	}
// 	address := ethCommon.BytesToAddress(addressBytes)

// 	valueBytes := make([]byte, 32)
// 	_, err = dataBuffer.Read(valueBytes)
// 	if err != nil {
// 		return
// 	}
// 	amountWei := big.NewInt(0).SetBytes(valueBytes)

// 	txn := EtherscanEthereumInternalTransaction{
// 		ToAddress: address,
// 		Amount:    WeiToEtherscanEthereum(amountWei),
// 	}

// 	return true, &txn, nil
// }
