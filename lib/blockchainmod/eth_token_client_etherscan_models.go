package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type EtherscanEthereumTokenBlock struct {
	EtherscanEthereumBlock

	tokenTxns []EtherscanEthereumTokenTransaction
}

func (this *EtherscanEthereumTokenBlock) setTxns(txns []EtherscanEthereumTokenTransaction) {
	this.tokenTxns = txns
}

func (this *EtherscanEthereumTokenBlock) GetTransactions() ([]Transaction, error) {
	txns := make([]Transaction, len(this.tokenTxns))
	for i := range this.tokenTxns {
		txns[i] = &this.tokenTxns[i]
	}
	return txns, nil
}

type EtherscanEthereumTokenTransaction struct {
	baseTransaction

	TokenSymbol string `json:"tokenSymbol"`
	Hash        string `json:"hash"`
	BlockHash   string `json:"blockHash"`
	BlockHeight uint64 `json:"blockNumber,string"`
	Nonce       uint64 `json:"nonce,string"`

	ContractAddress string          `json:"contractAddress"`
	FromAddress     string          `json:"from"`
	ToAddress       string          `json:"to"`
	AmountWei       decimal.Decimal `json:"value,string"`
	Confirmations   uint64          `json:"confirmations,string"`
	Time            int64           `json:"timeStamp,string"`

	GasPriceWei uint64 `json:"gasPrice,string"`
	GasLimit    uint32 `json:"gas,string"`
	GasUsed     uint32 `json:"gasUsed,string"`
}

func (this *EtherscanEthereumTokenTransaction) GetCurrency() meta.Currency {
	return meta.NewCurrency(this.TokenSymbol)
}

func (this *EtherscanEthereumTokenTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *EtherscanEthereumTokenTransaction) GetDirection() meta.Direction {
	switch this.ownerAddress {
	case this.FromAddress:
		return constants.DirectionTypeSend
	case this.ToAddress:
		return constants.DirectionTypeReceive
	default:
		return constants.DirectionTypeUnknown
	}
}

func (this *EtherscanEthereumTokenTransaction) GetFee() meta.CurrencyAmount {
	gasTotalWei := this.GasPriceWei * uint64(this.GasUsed)
	return currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    comutils.NewDecimalF(gasTotalWei),
		},
		constants.CurrencyEthereum,
	)
}

func (this *EtherscanEthereumTokenTransaction) GetHash() string {
	return this.Hash
}

func (this *EtherscanEthereumTokenTransaction) GetBlockHash() string {
	return this.BlockHash
}

func (this *EtherscanEthereumTokenTransaction) GetBlockHeight() uint64 {
	return this.BlockHeight
}

func (this *EtherscanEthereumTokenTransaction) GetConfirmations() uint64 {
	return this.Confirmations
}

func (this *EtherscanEthereumTokenTransaction) GetFromAddress() (string, error) {
	return this.FromAddress, nil
}

func (this *EtherscanEthereumTokenTransaction) GetToAddress() (string, error) {
	return this.ToAddress, nil
}

func (this *EtherscanEthereumTokenTransaction) GetAmount() (decimal.Decimal, error) {
	amountEth := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: constants.CurrencySubEthereumWei,
			Value:    comutils.NewDecimalF(this.AmountWei),
		},
		constants.CurrencyEthereum,
	)
	return amountEth.Value, nil
}

func (this *EtherscanEthereumTokenTransaction) GetInputDataHex() (string, error) {
	return "", nil
}

func (this *EtherscanEthereumTokenTransaction) GetTimeUnix() int64 {
	return this.Time
}

func (this *EtherscanEthereumTokenTransaction) GetInputs() ([]Input, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction inputs")
}

func (this *EtherscanEthereumTokenTransaction) GetOutputs() ([]Output, error) {
	return nil, utils.IssueErrorf("Ethereum doesn't have transaction outputs")
}
