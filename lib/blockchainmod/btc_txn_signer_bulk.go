package blockchainmod

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type bitcoinBulkTxnSignerInputSource struct {
	wif     *btcutil.WIF
	Address btcutil.Address
}

type bitcoinBulkTxnSignerOutput struct {
	Amount  decimal.Decimal
	Address btcutil.Address
}

type BitcoinBulkTxnSigner struct {
	baseTxnSigner

	leftoverAddress btcutil.Address
	chainConfig     *chaincfg.Params

	isMove       bool
	inputSources []bitcoinBulkTxnSignerInputSource
	outputs      []bitcoinBulkTxnSignerOutput

	signedTxn *wire.MsgTx
	signedFee decimal.Decimal
}

func NewBitcoinBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	if config.BlockchainUseTestnet {
		return NewBitcoinTestnetBulkTxnSigner(client, feeInfo)
	} else {
		return NewBitcoinMainnetBulkTxnSigner(client, feeInfo)
	}
}

func NewBitcoinMainnetBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	return &BitcoinBulkTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinChainConfig,
	}
}

func NewBitcoinTestnetBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	return &BitcoinBulkTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinChainConfigTestnet,
	}
}

func (this *BitcoinBulkTxnSigner) SetLeftoverAddress(address string) error {
	leftoverAddress, err := btcutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return utils.WrapError(err)
	}

	this.leftoverAddress = leftoverAddress
	return nil
}

func (this *BitcoinBulkTxnSigner) GetRaw() []byte {
	var txnBuffer bytes.Buffer
	comutils.PanicOnError(
		this.signedTxn.Serialize(&txnBuffer),
	)

	return txnBuffer.Bytes()
}

func (this *BitcoinBulkTxnSigner) GetRawHex() string {
	return comutils.HexEncode(this.GetRaw())
}

func (this *BitcoinBulkTxnSigner) GetHash() string {
	return this.signedTxn.TxHash().String()
}

func (this *BitcoinBulkTxnSigner) GetAmount() decimal.Decimal {
	panic(utils.IssueErrorf("not implemented"))
}

func (this *BitcoinBulkTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: constants.CurrencyBitcoin,
		Value:    this.signedFee,
	}
	if this.signedTxn == nil {
		return feeAmount, utils.IssueErrorf("bitcoin txn must be signed before getting fee")
	}

	return feeAmount, nil
}

func (this *BitcoinBulkTxnSigner) AddSrc(privateKey string, hintAddress string) (err error) {
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return utils.WrapError(err)
	}

	var hintAddressObj btcutil.Address
	if hintAddress != "" {
		hintAddressObj, err = btcutil.DecodeAddress(hintAddress, this.chainConfig)
		if err != nil {
			return utils.WrapError(err)
		}
	}
	srcAddress, err := ParseBitcoinWifAddress(wif, this.chainConfig, hintAddressObj)
	if err != nil {
		return err
	}

	inputSource := bitcoinBulkTxnSignerInputSource{
		wif:     wif,
		Address: srcAddress,
	}
	this.inputSources = append(this.inputSources, inputSource)

	return nil
}

func (this *BitcoinBulkTxnSigner) AddDst(address string, amount decimal.Decimal) error {
	if this.isMove {
		return utils.IssueErrorf(
			"bitcoin txn signer bulk cannot add destination address in MOVE mode",
		)
	}

	dstAddress, err := btcutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return err
	}

	output := bitcoinBulkTxnSignerOutput{
		Amount:  amount,
		Address: dstAddress,
	}
	this.outputs = append(this.outputs, output)

	return nil
}

func (this *BitcoinBulkTxnSigner) SetMoveDst(address string) error {
	if err := this.SetLeftoverAddress(address); err != nil {
		return err
	}

	this.isMove = true
	return nil
}

func (this *BitcoinBulkTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	txn := wire.NewMsgTx(wire.TxVersion)

	var outputKeep *wire.TxOut
	if this.leftoverAddress != nil {
		leftoverPayScript, err := txscript.PayToAddrScript(this.leftoverAddress)
		if err != nil {
			return utils.WrapError(err)
		}

		outputKeep = wire.NewTxOut(0, leftoverPayScript)
		txn.AddTxOut(outputKeep)
	}

	outputTotalAmount := decimal.Zero
	for _, output := range this.outputs {
		outputPayScript, err := txscript.PayToAddrScript(output.Address)
		if err != nil {
			return utils.WrapError(err)
		}

		txn.AddTxOut(wire.NewTxOut(BitcoinToSatoshi(output.Amount), outputPayScript))
		outputTotalAmount = outputTotalAmount.Add(output.Amount)
	}

	var (
		totalInputAmount = decimal.Zero
		inputToSrcMap    = make(map[*wire.TxIn]bitcoinBulkTxnSignerInputSource)
		inputUtxOutputs  = make([]UnspentTxnOutput, 0, len(this.inputSources))
		inputCountLegacy uint32
		inputCountSegWit uint32
	)
	for _, inputSrc := range this.inputSources {
		address := inputSrc.Address.EncodeAddress()
		utxOutputs, err := this.client.GetUtxOutputs(
			address,
			decimal.NewFromInt(MaxBalance))
		if err != nil {
			return err
		}
		for _, utxOutput := range utxOutputs {
			var utxoTxnHash chainhash.Hash
			err = chainhash.Decode(&utxoTxnHash, utxOutput.GetTxnHash())
			if err != nil {
				return utils.WrapError(err)
			}

			outpoint := wire.NewOutPoint(&utxoTxnHash, utxOutput.GetIndex())
			input := wire.NewTxIn(outpoint, nil, nil)

			txn.AddTxIn(input)

			totalInputAmount = totalInputAmount.Add(utxOutput.GetAmount())
			inputToSrcMap[input] = inputSrc
			inputUtxOutputs = append(inputUtxOutputs, utxOutput)
		}

		utxOutputCount := uint32(len(utxOutputs))
		if IsBitcoinSegWitAddress(address) {
			inputCountSegWit += utxOutputCount
		} else {
			inputCountLegacy += utxOutputCount
		}
	}
	if len(inputUtxOutputs) == 0 {
		return utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	txnEstimatedSize := EstimateBitcoinMixedTxnSize(
		inputCountSegWit, inputCountLegacy,
		uint32(len(txn.TxOut)))
	feeAmount := comutils.DecimalClamp(
		this.feeInfo.GetBasePrice().Mul(decimal.NewFromInt(int64(txnEstimatedSize))),
		this.feeInfo.GetBaseLimitMinValue(),
		this.feeInfo.GetBaseLimitMaxValue(),
	)
	leftoverAmount := totalInputAmount.Sub(outputTotalAmount).Sub(feeAmount)
	if leftoverAmount.IsNegative() {
		return utils.WrapError(constants.ErrorBalanceNotEnough)
	}
	if outputKeep == nil && leftoverAmount.GreaterThanOrEqual(feeAmount) {
		return utils.IssueErrorf(
			"bitcoin bulk txn has a large keep amount to be wasted | fee_amount=%v,keep_amount=%v",
			feeAmount, leftoverAmount,
		)
	}
	if outputKeep != nil {
		outputKeep.Value = BitcoinToSatoshi(leftoverAmount)

		if leftoverAmount.IsZero() {
			txn.TxOut = txn.TxOut[1:]
		}
	}
	if len(txn.TxOut) < 1 {
		return utils.IssueErrorf("bitcoin bulk txn has an empty output set")
	}

	for i, input := range txn.TxIn {
		var (
			inputSrc       = inputToSrcMap[input]
			inputUtxOutput = inputUtxOutputs[i]
		)
		inputPayScript, err := txscript.PayToAddrScript(inputSrc.Address)
		if err != nil {
			return utils.WrapError(err)
		}

		if IsBitcoinSegWitAddress(inputUtxOutput.GetAddress()) {
			inputValue := BitcoinToSatoshi(inputUtxOutput.GetAmount())
			witness, err := txscript.WitnessSignature(
				txn, txscript.NewTxSigHashes(txn),
				i, inputValue, inputPayScript,
				txscript.SigHashAll, inputSrc.wif.PrivKey, true)
			if err != nil {
				return utils.WrapError(err)
			}
			input.Witness = witness
		} else {
			sigScript, err := txscript.SignatureScript(
				txn,
				i, inputPayScript,
				txscript.SigHashAll, inputSrc.wif.PrivKey, true)
			if err != nil {
				return utils.WrapError(err)
			}

			input.SignatureScript = sigScript
		}
	}

	this.signedTxn = txn
	this.signedFee = feeAmount

	return nil
}
