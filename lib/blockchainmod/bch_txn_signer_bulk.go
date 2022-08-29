package blockchainmod

import (
	"bytes"

	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/txscript"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type bitcoinCashBulkTxnSignerInputSource struct {
	wif     *bchutil.WIF
	Address bchutil.Address
}

type bitcoinCashBulkTxnSignerOutput struct {
	Amount  decimal.Decimal
	Address bchutil.Address
}

type BitcoinCashBulkTxnSigner struct {
	baseTxnSigner
	coin Coin

	leftoverAddress bchutil.Address
	chainConfig     *chaincfg.Params

	isMove       bool
	inputSources []bitcoinCashBulkTxnSignerInputSource
	outputs      []bitcoinCashBulkTxnSignerOutput

	signedTxn *wire.MsgTx
	signedFee decimal.Decimal
}

func NewBitcoinCashMainnetBulkTxnSigner(coin Coin, client Client, feeInfo FeeInfo) *BitcoinCashBulkTxnSigner {
	return &BitcoinCashBulkTxnSigner{
		coin: coin,
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinCashChainConfig,
	}
}

func NewBitcoinCashTestnetBulkTxnSigner(coin Coin, client Client, feeInfo FeeInfo) *BitcoinCashBulkTxnSigner {
	return &BitcoinCashBulkTxnSigner{
		coin: coin,
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinCashChainConfigTestnet,
	}
}

func (this *BitcoinCashBulkTxnSigner) SetLeftoverAddress(address string) error {
	leftoverAddress, err := bchutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return utils.WrapError(err)
	}

	this.leftoverAddress = leftoverAddress
	return nil
}

func (this *BitcoinCashBulkTxnSigner) GetRaw() []byte {
	var txnBuffer bytes.Buffer
	comutils.PanicOnError(
		this.signedTxn.Serialize(&txnBuffer),
	)

	return txnBuffer.Bytes()
}

func (this *BitcoinCashBulkTxnSigner) GetRawHex() string {
	return comutils.HexEncode(this.GetRaw())
}

func (this *BitcoinCashBulkTxnSigner) GetHash() string {
	return this.signedTxn.TxHash().String()
}

func (this *BitcoinCashBulkTxnSigner) GetAmount() decimal.Decimal {
	panic(utils.IssueErrorf("not implemented"))
}

func (this *BitcoinCashBulkTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: this.coin.GetCurrency(),
		Value:    this.signedFee,
	}
	if this.signedTxn == nil {
		return feeAmount, utils.IssueErrorf("bitcoin cash txn must be signed before getting fee")
	}

	return feeAmount, nil
}

func (this *BitcoinCashBulkTxnSigner) AddSrc(privateKey string, hintAddress string) (err error) {
	wif, err := bchutil.DecodeWIF(privateKey)
	if err != nil {
		return utils.WrapError(err)
	}
	srcAddress, err := GetBitcoinCashWifAddressP2PKH(wif, this.chainConfig)
	if err != nil {
		return err
	}

	inputSource := bitcoinCashBulkTxnSignerInputSource{
		wif:     wif,
		Address: srcAddress,
	}
	this.inputSources = append(this.inputSources, inputSource)

	return nil
}

func (this *BitcoinCashBulkTxnSigner) AddDst(address string, amount decimal.Decimal) error {
	if this.isMove {
		return utils.IssueErrorf(
			"bitcoin cash txn signer bulk cannot add destination address in MOVE mode",
		)
	}

	dstAddress, err := bchutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return err
	}

	output := bitcoinCashBulkTxnSignerOutput{
		Amount:  amount,
		Address: dstAddress,
	}
	this.outputs = append(this.outputs, output)

	return nil
}

func (this *BitcoinCashBulkTxnSigner) SetMoveDst(address string) error {
	if err := this.SetLeftoverAddress(address); err != nil {
		return err
	}

	this.isMove = true
	return nil
}

func (this *BitcoinCashBulkTxnSigner) Sign(canOverwrite bool) (err error) {
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
		inputToSrcMap    = make(map[*wire.TxIn]bitcoinCashBulkTxnSignerInputSource)
		inputUtxOutputs  = make([]UnspentTxnOutput, 0, len(this.inputSources))
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
			input := wire.NewTxIn(outpoint, nil)

			txn.AddTxIn(input)

			totalInputAmount = totalInputAmount.Add(utxOutput.GetAmount())
			inputToSrcMap[input] = inputSrc
			inputUtxOutputs = append(inputUtxOutputs, utxOutput)
		}
	}
	if len(inputUtxOutputs) == 0 {
		return utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	txnEstimatedSize := EstimateBitcoinLegacyTxnSize(
		uint32(len(inputUtxOutputs)),
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
			"bitcoin cash bulk txn has a large keep amount to be wasted | fee_amount=%v,keep_amount=%v",
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
		return utils.IssueErrorf("bitcoin cash bulk txn has an empty output set")
	}

	for i, input := range txn.TxIn {
		var (
			inputSrc       = inputToSrcMap[input]
			inputUtxOutput = inputUtxOutputs[i]
			inputValue     = BitcoinToSatoshi(inputUtxOutput.GetAmount())
		)
		inputPayScript, err := txscript.PayToAddrScript(inputSrc.Address)
		if err != nil {
			return utils.WrapError(err)
		}

		sigScript, err := txscript.SignatureScript(
			txn,
			i, inputValue, inputPayScript,
			txscript.SigHashAll, inputSrc.wif.PrivKey, true)
		if err != nil {
			return utils.WrapError(err)
		}

		input.SignatureScript = sigScript
	}

	this.signedTxn = txn
	this.signedFee = feeAmount

	return nil
}
