package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type tBitcoreBaseBlock struct {
	txns []Transaction

	Currency          meta.Currency `json:"chain"`
	Network           string        `json:"network"`
	Hash              string        `json:"hash"`
	Height            int64         `json:"height"`
	Size              uint64        `json:"size"`
	PreviousBlockHash string        `json:"previousBlockHash"`
	TransactionCount  uint32        `json:"transactionCount"`
	Confirmations     int64         `json:"confirmations"`
	Timestamp         int64         `json:"timestamp"`
}

func (this *tBitcoreBaseBlock) setTxns(txns []Transaction) {
	this.txns = txns
}

func (this *tBitcoreBaseBlock) isMainChain() bool {
	return this.Network == BitcoreChainMainnet
}

func (this *tBitcoreBaseBlock) GetCurrency() meta.Currency {
	return this.Currency
}

func (this *tBitcoreBaseBlock) GetNetwork() meta.BlockchainNetwork {
	coin := GetCoinNativeF(this.Currency)
	if this.Network == BitcoreChainMainnet {
		return coin.GetNetworkMain()
	} else {
		return coin.GetNetworkTest()
	}
}

func (this *tBitcoreBaseBlock) GetTimeUnix() int64 {
	return TimeMsToS(this.Timestamp)
}

func (this *tBitcoreBaseBlock) GetHash() string {
	return this.Hash
}

func (this *tBitcoreBaseBlock) GetHeight() uint64 {
	if this.Height < 0 {
		return 0
	} else {
		return uint64(this.Height)
	}
}

func (this *tBitcoreBaseBlock) GetParentHash() string {
	return this.PreviousBlockHash
}

func (this *tBitcoreBaseBlock) GetTransactions() ([]Transaction, error) {
	return this.txns, nil
}

type tBitcoreBaseTransaction struct {
	baseTransaction
	currencyUnit meta.Currency
	currencyBase meta.Currency

	Hash          string          `json:"txid"`
	Network       string          `json:"network"`
	Currency      meta.Currency   `json:"chain"`
	BlockHeight   int64           `json:"blockHeight"`
	BlockHash     string          `json:"blockHash"`
	Fee           decimal.Decimal `json:"fee"`
	Value         decimal.Decimal `json:"value"`
	Nonce         uint64          `json:"nonce"`
	Confirmations int64           `json:"confirmations"`
	Timestamp     int64           `json:"timestamp"`
}

func (this *tBitcoreBaseTransaction) SetCurrencyPair(baseCurrency, unitCurrency meta.Currency) {
	this.currencyBase = baseCurrency
	this.currencyUnit = unitCurrency
}

func (this *tBitcoreBaseTransaction) isMainChain() bool {
	return this.Network == BitcoreChainMainnet
}

func (this *tBitcoreBaseTransaction) isConflicted() bool {
	return this.BlockHeight == BitcoreBlockHeightAsConflict
}

func (this *tBitcoreBaseTransaction) isConfirmed() bool {
	return this.Confirmations > 0
}

func (this *tBitcoreBaseTransaction) GetCurrency() meta.Currency {
	return this.Currency
}

func (this *tBitcoreBaseTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.isConfirmed() {
		return constants.BlockchainTxnStatusSucceeded
	} else if this.isConflicted() {
		return constants.BlockchainTxnStatusFailed
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *tBitcoreBaseTransaction) GetFee() meta.CurrencyAmount {
	feeAmount := meta.CurrencyAmount{
		Currency: this.currencyUnit,
		Value:    this.Fee,
	}
	return currencymod.ConvertAmountF(feeAmount, this.currencyBase)
}

func (this *tBitcoreBaseTransaction) GetHash() string {
	return this.Hash
}

func (this *tBitcoreBaseTransaction) GetConfirmations() uint64 {
	if this.Confirmations < 0 {
		return 0
	} else {
		return uint64(this.Confirmations)
	}
}

func (this *tBitcoreBaseTransaction) GetBlockHeight() uint64 {
	if this.BlockHeight < 0 {
		return 0
	} else {
		return uint64(this.BlockHeight)
	}
}

func (this *tBitcoreBaseTransaction) GetAmount() (decimal.Decimal, error) {
	currencyAmount := currencymod.ConvertAmountF(
		meta.CurrencyAmount{
			Currency: this.currencyUnit,
			Value:    this.Value,
		},
		this.currencyBase,
	)
	return currencyAmount.Value, nil
}

func (this *tBitcoreBaseTransaction) GetTimeUnix() int64 {
	return TimeMsToS(this.Timestamp)
}
