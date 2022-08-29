package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type SoChainBitcoinBlock struct {
	client     Client
	cachedTxns []Transaction

	Network       string `json:"network"`
	Hash          string `json:"block_hash"`
	Height        uint64 `json:"block_no"`
	Time          int64  `json:"time"`
	Confirmations uint64 `json:"confirmations"`
	Size          uint32 `json:"size"`
	PrevBlockHash string `json:"previous_blockhash"`

	TxnHashes []string `json:"txs"`
}

func (this *SoChainBitcoinBlock) setClient(client Client) {
	this.client = client
}

func (this *SoChainBitcoinBlock) GetCurrency() meta.Currency {
	return getSochainNetworkCurrency(this.Network)
}

func (this *SoChainBitcoinBlock) GetNetwork() meta.BlockchainNetwork {
	switch this.Network {

	case soChainNetworkBitcoin:
		return constants.BlockchainNetworkBitcoin
	case soChainNetworkBitcoin + soChainNetworkTestSuffix:
		return constants.BlockchainNetworkBitcoinTestnet

	case soChainNetworkLitecoin:
		return constants.BlockchainNetworkLitecoin
	case soChainNetworkLitecoin + soChainNetworkTestSuffix:
		return constants.BlockchainNetworkLitecoinTestnet

	default:
		return "unknown"
	}
}

func (this *SoChainBitcoinBlock) GetHash() string {
	return this.Hash
}

func (this *SoChainBitcoinBlock) GetHeight() uint64 {
	return this.Height
}

func (this *SoChainBitcoinBlock) GetTimeUnix() int64 {
	return this.Time
}

func (this *SoChainBitcoinBlock) GetParentHash() string {
	return this.PrevBlockHash
}

func (this *SoChainBitcoinBlock) GetTransactions() ([]Transaction, error) {
	if this.cachedTxns == nil {
		txns := make([]Transaction, 0, len(this.TxnHashes))
		for _, hash := range this.TxnHashes {
			txn, err := this.client.GetTxn(hash)
			if err != nil {
				return nil, err
			}
			txns = append(txns, txn)
		}
		this.cachedTxns = txns
	}

	return this.cachedTxns, nil
}

type SoChainBitcoinTransaction struct {
	baseTransaction

	cacheDirection *meta.Direction
	cacheAmount    *decimal.Decimal

	Network       string `json:"network"`
	BlockHash     string `json:"blockhash"`
	BlockHeight   uint64 `json:"block_no"`
	Confirmations uint64 `json:"confirmations"`

	Hash      string          `json:"txid"`
	Time      int64           `json:"time"`
	Size      uint32          `json:"size"`
	AmoutSent decimal.Decimal `json:"sent_value"`
	AmoutFee  decimal.Decimal `json:"fee"`
	Hex       string          `json:"tx_hex"`

	Inputs  []SoChainBitcoinInput  `json:"inputs"`
	Outputs []SoChainBitcoinOutput `json:"outputs"`
}

func (this *SoChainBitcoinTransaction) GetCurrency() meta.Currency {
	return getSochainNetworkCurrency(this.Network)
}

func (this *SoChainBitcoinTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}

}

func (this *SoChainBitcoinTransaction) getDirection() meta.Direction {
	if this.ownerAddress == "" {
		return constants.DirectionTypeUnknown
	}

	for _, input := range this.Inputs {
		if input.GetPrevOutAddress() == this.ownerAddress {
			return constants.DirectionTypeSend
		}
	}

	return constants.DirectionTypeReceive
}

func (this *SoChainBitcoinTransaction) GetDirection() meta.Direction {
	if this.cacheDirection != nil {
		return *this.cacheDirection
	}

	direction := this.getDirection()
	this.cacheDirection = &direction

	return direction
}

func (this *SoChainBitcoinTransaction) GetFee() meta.CurrencyAmount {
	return meta.CurrencyAmount{
		Currency: this.GetCurrency(),
		Value:    this.AmoutFee,
	}
}

func (this *SoChainBitcoinTransaction) GetHash() string {
	return this.Hash
}

func (this *SoChainBitcoinTransaction) GetBlockHash() string {
	return this.BlockHash
}

func (this *SoChainBitcoinTransaction) GetBlockHeight() uint64 {
	panic(utils.IssueErrorf("not supported yet - SoChainBitcoinTransaction.GetBlockHeight"))
}

func (this *SoChainBitcoinTransaction) GetConfirmations() uint64 {
	return this.Confirmations
}

func (this *SoChainBitcoinTransaction) GetFromAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeSend {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction from address")
}

func (this *SoChainBitcoinTransaction) GetToAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeReceive {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction to address")
}

func (this *SoChainBitcoinTransaction) GetAmount() (decimal.Decimal, error) {
	if this.cacheAmount != nil {
		return *this.cacheAmount, nil
	}

	if this.ownerAddress == "" {
		return decimal.Zero, utils.IssueErrorf("cannot determine bitcoin transaction amount without owner address")
	}

	amount := decimal.Zero
	for _, input := range this.Inputs {
		if input.GetPrevOutAddress() == this.ownerAddress {
			amount = amount.Sub(input.GetPrevOutAmount())
		}
	}
	for _, output := range this.Outputs {
		if output.GetAddress() == this.ownerAddress {
			amount = amount.Add(output.GetAmount())
		}
	}

	amount = amount.Abs()
	this.cacheAmount = &amount

	return amount, nil
}

func (this *SoChainBitcoinTransaction) GetInputDataHex() (string, error) {
	return "", nil
}

func (this *SoChainBitcoinTransaction) GetTimeUnix() int64 {
	return this.Time
}

func (this *SoChainBitcoinTransaction) GetInputs() ([]Input, error) {
	inputs := make([]Input, len(this.Inputs))

	for i := range this.Inputs {
		inputs[i] = &this.Inputs[i]
	}

	return inputs, nil
}

func (this *SoChainBitcoinTransaction) GetOutputs() ([]Output, error) {
	outputs := make([]Output, len(this.Outputs))

	for i := range this.Outputs {
		outputs[i] = &this.Outputs[i]
	}

	return outputs, nil
}

type SoChainBitcoinInputOutpoint struct {
	Hash  string `json:"txid"`
	Index uint32 `json:"output_no"`
}

type SoChainBitcoinInput struct {
	network string

	Address      string                      `json:"address"`
	Value        decimal.Decimal             `json:"value"`
	Outpoint     SoChainBitcoinInputOutpoint `json:"received_from"`
	SigScriptHex string                      `json:"script_hex"`
}

func (this *SoChainBitcoinInput) GetPrevOutHash() string {
	return this.Outpoint.Hash
}

func (this *SoChainBitcoinInput) GetPrevOutIndex() uint32 {
	return this.Outpoint.Index
}

func (this *SoChainBitcoinInput) GetPrevOutAddress() string {
	return this.Address
}

func (this *SoChainBitcoinInput) GetPrevOutAmount() decimal.Decimal {
	return this.Value
}

type SoChainBitcoinOutput struct {
	Index   uint32          `json:"output_no"`
	Address string          `json:"address"`
	Value   decimal.Decimal `json:"value"`
}

func (this *SoChainBitcoinOutput) GetIndex() uint32 {
	return this.Index
}

func (this *SoChainBitcoinOutput) GetAddress() string {
	return this.Address
}

func (this *SoChainBitcoinOutput) GetAmount() decimal.Decimal {
	return this.Value
}

type SoChainBitcoinUtxOutput struct {
	address string

	TxnHash       string          `json:"txid"`
	Index         uint32          `json:"output_no"`
	Value         decimal.Decimal `json:"value"`
	Confirmations int64           `json:"confirmations"`
}

func (this *SoChainBitcoinUtxOutput) SetAddress(address string) {
	this.address = address
}

func (this *SoChainBitcoinUtxOutput) GetAddress() string {
	return this.address
}

func (this *SoChainBitcoinUtxOutput) GetTxnHash() string {
	return this.TxnHash
}

func (this *SoChainBitcoinUtxOutput) GetIndex() uint32 {
	return this.Index
}

type SoChainBitcoinDisplayTransaction struct {
	baseTransaction

	currency meta.Currency

	Hash          string `json:"txid"`
	BlockHeight   uint64 `json:"block_no"`
	Confirmations uint64 `json:"confirmations"`
	Time          int64  `json:"time"`

	Incoming SoChainBitcoinDisplayTxnIncoming `json:"incoming"`
	Outgoing SoChainBitcoinDisplayTxnOutgoing `json:"outgoing"`
}

func (this *SoChainBitcoinUtxOutput) GetAmount() decimal.Decimal {
	return this.Value
}

type SoChainBitcoinDisplayTxnIncoming struct {
	OutputNo uint32                           `json:"output_no"`
	Value    decimal.Decimal                  `json:"value"`
	Spent    SoChainBitcoinDisplayTxnOutPoint `json:"spent"`
	Inputs   []SoChainBitcoinDisplayTxnInput  `json:"inputs"`
}

type SoChainBitcoinDisplayTxnInput struct {
	Index        uint32                           `json:"input_no"`
	Address      string                           `json:"address"`
	ReceivedFrom SoChainBitcoinDisplayTxnOutPoint `json:"received_from"`
}

type SoChainBitcoinDisplayTxnOutPoint struct {
	Hash     string `json:"txid"`
	OutputNo uint32 `json:"output_no"`
}

type SoChainBitcoinDisplayTxnOutgoing struct {
	Value   decimal.Decimal                  `json:"value"`
	Outputs []SoChainBitcoinDisplayTxnOutput `json:"outputs"`
}

type SoChainBitcoinDisplayTxnOutput struct {
	Index   uint32                           `json:"output_no"`
	Address string                           `json:"address"`
	Value   decimal.Decimal                  `json:"value"`
	Spent   SoChainBitcoinDisplayTxnOutPoint `json:"spent"`
}

func (this *SoChainBitcoinDisplayTransaction) setCurrency(currency meta.Currency) {
	this.currency = currency
}

func (this *SoChainBitcoinDisplayTransaction) GetCurrency() meta.Currency {
	return this.currency
}

func (this *SoChainBitcoinDisplayTransaction) GetLocalStatus() meta.BlockchainTxnStatus {
	if this.GetConfirmations() > 0 {
		return constants.BlockchainTxnStatusSucceeded
	} else {
		return constants.BlockchainTxnStatusPending
	}
}

func (this *SoChainBitcoinDisplayTransaction) GetDirection() meta.Direction {
	if this.Outgoing.Value.IsZero() {
		return constants.DirectionTypeReceive
	} else {
		return constants.DirectionTypeSend
	}
}

func (this *SoChainBitcoinDisplayTransaction) GetFee() meta.CurrencyAmount {
	return meta.CurrencyAmount{
		Currency: constants.CurrencySubBitcoinSatoshi,
		Value:    decimal.Zero, // Unknown
	}
}

func (this *SoChainBitcoinDisplayTransaction) GetHash() string {
	return this.Hash
}

func (this *SoChainBitcoinDisplayTransaction) GetBlockHeight() uint64 {
	return this.BlockHeight
}

func (this *SoChainBitcoinDisplayTransaction) GetConfirmations() uint64 {
	return this.Confirmations
}

func (this *SoChainBitcoinDisplayTransaction) GetFromAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeSend {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction from address")
}

func (this *SoChainBitcoinDisplayTransaction) GetToAddress() (string, error) {
	if this.GetDirection() == constants.DirectionTypeReceive {
		return this.ownerAddress, nil
	}

	return "", utils.IssueErrorf("cannot determine bitcoin transaction to address")
}

func (this *SoChainBitcoinDisplayTransaction) GetAmount() (decimal.Decimal, error) {
	if this.GetDirection() == constants.DirectionTypeSend {
		return this.Outgoing.Value, nil
	} else {
		return this.Incoming.Value, nil
	}
}

func (this *SoChainBitcoinDisplayTransaction) GetInputDataHex() (string, error) {
	return "", nil
}

func (this *SoChainBitcoinDisplayTransaction) GetTimeUnix() int64 {
	return this.Time
}

func (this *SoChainBitcoinDisplayTransaction) GetInputs() ([]Input, error) {
	return nil, utils.IssueErrorf("not implemented - SoChainBitcoinDisplayTransaction.GetInputs")
}

func (this *SoChainBitcoinDisplayTransaction) GetOutputs() ([]Output, error) {
	return nil, utils.IssueErrorf("not implemented - SoChainBitcoinDisplayTransaction.GetOutputs")
}
