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

type BitcoinSingleTxnSigner struct {
	baseTxnSigner

	wif         *btcutil.WIF
	chainConfig *chaincfg.Params

	isMove     bool
	srcAddress btcutil.Address
	dstAddress btcutil.Address
	amount     decimal.Decimal

	signedTxn *wire.MsgTx
	signedFee decimal.Decimal
}

func NewBitcoinSingleTxnSigner(client Client, feeInfo FeeInfo) SingleTxnSigner {
	if config.BlockchainUseTestnet {
		return NewBitcoinTestnetSingleTxnSigner(client, feeInfo)
	} else {
		return NewBitcoinMainnetSingleTxnSigner(client, feeInfo)
	}
}

func NewBitcoinMainnetSingleTxnSigner(client Client, feeInfo FeeInfo) *BitcoinSingleTxnSigner {
	return &BitcoinSingleTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinChainConfig,
	}
}

func NewBitcoinTestnetSingleTxnSigner(client Client, feeInfo FeeInfo) *BitcoinSingleTxnSigner {
	return &BitcoinSingleTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinChainConfigTestnet,
	}
}

func (this *BitcoinSingleTxnSigner) GetRaw() []byte {
	var txnBuffer bytes.Buffer
	comutils.PanicOnError(
		this.signedTxn.Serialize(&txnBuffer),
	)

	return txnBuffer.Bytes()
}

func (this *BitcoinSingleTxnSigner) GetRawHex() string {
	return comutils.HexEncode(this.GetRaw())
}

func (this *BitcoinSingleTxnSigner) GetHash() string {
	return this.signedTxn.TxHash().String()
}

func (this *BitcoinSingleTxnSigner) GetAmount() decimal.Decimal {
	return this.amount
}

func (this *BitcoinSingleTxnSigner) GetSrcAddress() string {
	return this.srcAddress.EncodeAddress()
}

func (this *BitcoinSingleTxnSigner) GetDstAddress() string {
	return this.dstAddress.EncodeAddress()
}

func (this *BitcoinSingleTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: constants.CurrencyBitcoin,
		Value:    this.signedFee,
	}
	if this.signedTxn == nil {
		return feeAmount, utils.IssueErrorf("bitcoin txn must be signed before getting fee")
	}

	return feeAmount, nil
}

func (this *BitcoinSingleTxnSigner) SetSrc(privateKey string, hintAddress string) error {
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return utils.WrapError(err)
	}

	var parsedHintAddress btcutil.Address
	if hintAddress != "" {
		parsedHintAddress, err = btcutil.DecodeAddress(hintAddress, this.chainConfig)
		if err != nil {
			return utils.WrapError(err)
		}
	}
	srcAddress, err := ParseBitcoinWifAddress(wif, this.chainConfig, parsedHintAddress)
	if err != nil {
		return err
	}

	this.wif = wif
	this.srcAddress = srcAddress

	return nil
}

func (this *BitcoinSingleTxnSigner) SetDst(address string, amount decimal.Decimal) error {
	if this.isMove {
		return utils.IssueErrorf(
			"bitcoin txn signer single cannot set destination address in MOVE mode",
		)
	}

	dstAddress, err := btcutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return utils.WrapError(err)
	}
	this.dstAddress = dstAddress
	this.amount = amount

	return nil
}

func (this *BitcoinSingleTxnSigner) SetMoveDst(address string) error {
	dstAddress, err := btcutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return utils.WrapError(err)
	}
	this.dstAddress = dstAddress
	this.isMove = true

	return nil
}

func (this *BitcoinSingleTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	if !this.isMove && this.amount.IsZero() {
		err = utils.IssueErrorf(
			"bitcoin transaction with zero amount is not support | from_address=%v,to_address=%v",
			this.GetSrcAddress(), this.GetDstAddress(),
		)
		return
	}

	srcPayScript, err := txscript.PayToAddrScript(this.srcAddress)
	if err != nil {
		return utils.WrapError(err)
	}
	dstPayScript, err := txscript.PayToAddrScript(this.dstAddress)
	if err != nil {
		return utils.WrapError(err)
	}

	var (
		txn         = wire.NewMsgTx(wire.TxVersion)
		srcAddress  = this.GetSrcAddress()
		isSegWitTxn = IsBitcoinSegWitAddress(srcAddress)
	)

	utxOutputs, err := this.client.GetUtxOutputs(srcAddress, decimal.NewFromInt(MaxBalance))
	if err != nil {
		return
	}
	var (
		totalInputAmount = decimal.Zero
		inputUtxOutputs  = make([]UnspentTxnOutput, 0, len(utxOutputs))
	)
	for _, utxOutput := range utxOutputs {
		var utxoTxnHash chainhash.Hash
		err = chainhash.Decode(&utxoTxnHash, utxOutput.GetTxnHash())
		if err != nil {
			return utils.WrapError(err)
		}

		outpoint := wire.NewOutPoint(&utxoTxnHash, utxOutput.GetIndex())
		txn.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

		totalInputAmount = totalInputAmount.Add(utxOutput.GetAmount())
		inputUtxOutputs = append(inputUtxOutputs, utxOutput)
	}
	if len(inputUtxOutputs) == 0 {
		return utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	var feeAmount decimal.Decimal
	if this.isMove {
		var (
			txnEstimatedSize uint32
			inputCount       = uint32(len(inputUtxOutputs))
		)
		if isSegWitTxn {
			txnEstimatedSize = EstimateBitcoinSegwitTxnSize(inputCount, 1)
		} else {
			txnEstimatedSize = EstimateBitcoinLegacyTxnSize(inputCount, 1)
		}
		feeAmount = comutils.DecimalClamp(
			this.feeInfo.GetBasePrice().Mul(decimal.NewFromInt(int64(txnEstimatedSize))),
			this.feeInfo.GetBaseLimitMinValue(),
			this.feeInfo.GetBaseLimitMaxValue(),
		)
		var (
			moveAmount = totalInputAmount.Sub(feeAmount)
			outputMove = wire.NewTxOut(BitcoinToSatoshi(moveAmount), dstPayScript)
		)
		txn.AddTxOut(outputMove)
	} else {
		feeAmount = this.feeInfo.GetBaseLimitMaxValue()
		leftoverAmount := totalInputAmount.Sub(this.amount).Sub(feeAmount)
		if leftoverAmount.IsNegative() {
			err = constants.ErrorBalanceNotEnough
			return
		}
		if leftoverAmount.IsPositive() {
			outputKeep := wire.NewTxOut(BitcoinToSatoshi(leftoverAmount), srcPayScript)
			txn.AddTxOut(outputKeep)
		}
		outputTransfer := wire.NewTxOut(BitcoinToSatoshi(this.amount), dstPayScript)
		txn.AddTxOut(outputTransfer)
	}

	for idx, input := range txn.TxIn {
		if isSegWitTxn {
			var (
				inputUtxOutput = inputUtxOutputs[idx]
				inputValue     = BitcoinToSatoshi(inputUtxOutput.GetAmount())
			)
			witness, err := txscript.WitnessSignature(
				txn, txscript.NewTxSigHashes(txn),
				idx, inputValue, srcPayScript,
				txscript.SigHashAll, this.wif.PrivKey, true)
			if err != nil {
				return utils.WrapError(err)
			}
			input.Witness = witness
		} else {
			sigScript, err := txscript.SignatureScript(
				txn,
				idx, srcPayScript,
				txscript.SigHashAll, this.wif.PrivKey, true)
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
