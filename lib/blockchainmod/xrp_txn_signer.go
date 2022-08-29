package blockchainmod

import (
	"fmt"
	"math"

	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/data"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type RippleTxnSigner struct {
	baseTxnSigner

	key       crypto.Key
	keySeq    *uint32
	isMainnet bool

	amount     decimal.Decimal
	srcAddress *RippleXAddress
	dstAddress *RippleXAddress

	isMove bool

	signedTxn *data.Payment
}

func newRippleTxnSigner(client Client, feeInfo FeeInfo, isMainnet bool) *RippleTxnSigner {
	return &RippleTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		isMainnet: isMainnet,
	}
}

func NewRippleMainnetTxnSigner(client Client, feeInfo FeeInfo) *RippleTxnSigner {
	return newRippleTxnSigner(client, feeInfo, true)
}

func NewRippleTestnetTxnSigner(client Client, feeInfo FeeInfo) *RippleTxnSigner {
	return newRippleTxnSigner(client, feeInfo, false)
}

func (this *RippleTxnSigner) SetNonce(nonce Nonce) error {
	if _, err := nonce.GetNumber(); err != nil {
		return err
	}
	return this.baseTxnSigner.SetNonce(nonce)
}

func (this *RippleTxnSigner) GetRaw() []byte {
	_, bytes, err := data.Raw(this.signedTxn)
	comutils.PanicOnError(err)
	return bytes
}

func (this *RippleTxnSigner) GetHash() string {
	return this.signedTxn.Hash.String()
}

func (this *RippleTxnSigner) GetAmount() decimal.Decimal {
	return this.amount
}

func (this *RippleTxnSigner) GetSrcAddress() string {
	return this.srcAddress.String()
}

func (this *RippleTxnSigner) GetDstAddress() string {
	return this.srcAddress.String()
}

func (this *RippleTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: constants.CurrencyRipple,
		Value:    this.feeInfo.GetBaseLimitMaxValue(),
	}
	return feeAmount, nil
}

func (this *RippleTxnSigner) SetSrc(b58Seed string, hintAddress string) error {
	key, seq, err := LoadRippleSeedKey(b58Seed)
	if err != nil {
		return err
	}
	addressHash, err := crypto.AccountId(key, seq)
	if err != nil {
		return err
	}

	var tag *uint32
	if hintAddress != "" {
		hintXAddress, err := RippleParseXAddress(hintAddress, this.isMainnet)
		if err == nil {
			if hasTag, xTag := hintXAddress.GetTag(); hasTag {
				tag = &xTag
			}
		}
	}
	xAddress, err := NewRippleXAddress(addressHash.String(), tag, this.isMainnet)
	if err != nil {
		return err
	}
	this.key = key
	this.keySeq = seq
	this.srcAddress = &xAddress

	return nil
}

func (this *RippleTxnSigner) SetDst(address string, amount decimal.Decimal) error {
	xAddress, err := RippleParseXAddress(address, this.isMainnet)
	if err != nil {
		return err
	}

	this.dstAddress = &xAddress
	this.amount = amount

	return nil
}

func (this *RippleTxnSigner) SetMoveDst(address string) error {
	xAddress, err := RippleParseXAddress(address, this.isMainnet)
	if err != nil {
		return err
	}
	this.dstAddress = &xAddress
	this.isMove = true

	return nil
}

func (this *RippleTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	if this.srcAddress == nil {
		return utils.IssueErrorf("XRP txn signer needs source address to sign")
	}
	if this.dstAddress == nil {
		return utils.IssueErrorf("XRP txn signer needs destination address to sign")
	}

	nonce, err := this.loadNonce()
	if err != nil {
		return
	}
	txn, err := this.genTxn(nonce)
	if err != nil {
		return err
	}
	if err := data.Sign(&txn, this.key, this.keySeq); err != nil {
		return utils.WrapError(err)
	}

	this.signedTxn = &txn
	return nil
}

func (this *RippleTxnSigner) loadNonce() (uint32, error) {
	if this.nonce == nil {
		nonce, err := this.client.GetNextNonce(this.srcAddress.GetRootAddress())
		if err != nil {
			return 0, err
		}
		this.nonce = nonce
	}

	nonceUint64, err := this.nonce.GetNumber()
	if nonceUint64 > math.MaxUint32 {
		return 0, utils.IssueErrorf("Ripple nonce value %v is too high", nonceUint64)
	}

	return uint32(nonceUint64), err
}

func (this *RippleTxnSigner) genTxn(nonce uint32) (txn data.Payment, err error) {
	txn = data.Payment{
		TxBase: data.TxBase{TransactionType: data.PAYMENT},
	}
	var (
		totalFee       = this.feeInfo.GetBaseLimitMaxValue()
		selectedAmount decimal.Decimal
		leftAmount     decimal.Decimal
	)
	if this.isMove {
		balance, balanceErr := this.getBalance()
		if balanceErr != nil {
			err = balanceErr
			return
		}
		leftAmount = balance.Sub(totalFee)
		selectedAmount = leftAmount
	} else {
		selectedAmount = this.amount
		if this.isPreferOffline {
			leftAmount = comutils.DecimalOne // Ignore balance validation
		} else {
			balance, balanceErr := this.getBalance()
			if balanceErr != nil {
				err = balanceErr
				return
			}
			totalAmount := selectedAmount.Add(totalFee)
			leftAmount = balance.Sub(totalAmount)
		}
	}
	if leftAmount.IsNegative() {
		err = utils.WrapError(constants.ErrorBalanceNotEnough)
		return
	}

	fromAccount, err := data.NewAccountFromAddress(this.srcAddress.GetRootAddress())
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	toAccount, err := data.NewAccountFromAddress(this.dstAddress.GetRootAddress())
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	amountDesc := fmt.Sprintf("%v/XRP", selectedAmount)
	amount, err := data.NewAmount(amountDesc)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	feeValue, err := data.NewNativeValue(rippleToDrop(totalFee))
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	var (
		txnFlag = data.TxCanonicalSignature
	)
	txn.Sequence = nonce
	txn.Flags = &txnFlag
	txn.Account = *fromAccount
	txn.Destination = *toAccount
	txn.Amount = *amount
	txn.Fee = *feeValue
	if hasTag, tag := this.srcAddress.GetTag(); hasTag {
		txn.SourceTag = &tag
	}
	if hasTag, tag := this.dstAddress.GetTag(); hasTag {
		txn.DestinationTag = &tag
	}

	return txn, nil
}

func (this *RippleTxnSigner) getBalance() (_ decimal.Decimal, err error) {
	if this.srcAddress == nil {
		err = utils.IssueErrorf("XRP txn signer needs source address to get balance")
		return
	}
	return this.client.GetBalance(this.srcAddress.GetRootAddress())
}

func (this *RippleTxnSigner) GetSignedTxn() *data.Payment {
	return this.signedTxn
}
