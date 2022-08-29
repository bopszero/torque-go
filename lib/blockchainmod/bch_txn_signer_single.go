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

type BitcoinCashSingleTxnSigner struct {
	baseTxnSigner
	coin Coin

	wif         *bchutil.WIF
	chainConfig *chaincfg.Params

	isMove     bool
	srcAddress bchutil.Address
	dstAddress bchutil.Address
	amount     decimal.Decimal

	signedTxn *wire.MsgTx
	signedFee decimal.Decimal
}

func NewBitcoinCashMainnetSingleTxnSigner(
	coin Coin, client Client, feeInfo FeeInfo,
) *BitcoinCashSingleTxnSigner {
	return &BitcoinCashSingleTxnSigner{
		coin: coin,
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinCashChainConfig,
	}
}

func NewBitcoinCashTestnetSingleTxnSigner(
	coin Coin, client Client, feeInfo FeeInfo,
) *BitcoinCashSingleTxnSigner {
	return &BitcoinCashSingleTxnSigner{
		coin: coin,
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &BitcoinCashChainConfigTestnet,
	}
}

func (this *BitcoinCashSingleTxnSigner) GetRaw() []byte {
	var txnBuffer bytes.Buffer
	comutils.PanicOnError(
		this.signedTxn.Serialize(&txnBuffer),
	)

	return txnBuffer.Bytes()
}

func (this *BitcoinCashSingleTxnSigner) GetRawHex() string {
	return comutils.HexEncode(this.GetRaw())
}

func (this *BitcoinCashSingleTxnSigner) GetHash() string {
	return this.signedTxn.TxHash().String()
}

func (this *BitcoinCashSingleTxnSigner) GetAmount() decimal.Decimal {
	return this.amount
}

func (this *BitcoinCashSingleTxnSigner) GetSrcAddress() string {
	return this.srcAddress.EncodeAddress()
}

func (this *BitcoinCashSingleTxnSigner) GetDstAddress() string {
	return this.dstAddress.EncodeAddress()
}

func (this *BitcoinCashSingleTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: this.coin.GetCurrency(),
		Value:    this.signedFee,
	}
	if this.signedTxn == nil {
		return feeAmount, utils.IssueErrorf("bitcoin cash txn must be signed before getting fee")
	}

	return feeAmount, nil
}

func (this *BitcoinCashSingleTxnSigner) SetSrc(privateKey string, _ string) error {
	wif, err := bchutil.DecodeWIF(privateKey)
	if err != nil {
		return utils.WrapError(err)
	}
	srcAddress, err := GetBitcoinCashWifAddressP2PKH(wif, this.chainConfig)
	if err != nil {
		return err
	}

	this.wif = wif
	this.srcAddress = srcAddress

	return nil
}

func (this *BitcoinCashSingleTxnSigner) SetDst(address string, amount decimal.Decimal) (err error) {
	dstAddress, err := bchutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return
	}

	this.dstAddress = dstAddress
	this.amount = amount

	return
}

func (this *BitcoinCashSingleTxnSigner) SetMoveDst(address string) error {
	dstAddress, err := bchutil.DecodeAddress(address, this.chainConfig)
	if err != nil {
		return utils.WrapError(err)
	}
	this.dstAddress = dstAddress
	this.isMove = true

	return nil
}

func (this *BitcoinCashSingleTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	if !this.isMove && this.amount.IsZero() {
		err = utils.IssueErrorf(
			"bitcoin cash transaction with zero amount is not support | from_address=%v,to_address=%v",
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

	txn := wire.NewMsgTx(wire.TxVersion)

	utxOutputs, err := this.client.GetUtxOutputs(
		this.GetSrcAddress(),
		decimal.NewFromInt(MaxBalance))
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
		txn.AddTxIn(wire.NewTxIn(outpoint, nil))

		totalInputAmount = totalInputAmount.Add(utxOutput.GetAmount())
		inputUtxOutputs = append(inputUtxOutputs, utxOutput)
	}
	if len(inputUtxOutputs) == 0 {
		return utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	var feeAmount decimal.Decimal
	if this.isMove {
		var (
			inputCount       = uint32(len(inputUtxOutputs))
			txnEstimatedSize = EstimateBitcoinLegacyTxnSize(inputCount, 1)
		)
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
		var (
			inputUtxOutput = inputUtxOutputs[idx]
			inputValue     = BitcoinToSatoshi(inputUtxOutput.GetAmount())
		)
		sigScript, err := txscript.SignatureScript(
			txn,
			idx, inputValue, srcPayScript,
			txscript.SigHashAll, this.wif.PrivKey, true)
		if err != nil {
			return utils.WrapError(err)
		}
		input.SignatureScript = sigScript
	}

	this.signedTxn = txn
	this.signedFee = feeAmount

	return nil
}
