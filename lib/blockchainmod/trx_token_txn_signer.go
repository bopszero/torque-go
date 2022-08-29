package blockchainmod

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tronaddr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type TronTokenTxnSigner struct {
	*TronTxnSigner

	tokenMeta TokenMeta
}

func NewTronTokenMainnetTxnSigner(
	tokenMeta TokenMeta,
	client Client, feeInfo FeeInfo,
) *TronTokenTxnSigner {
	return &TronTokenTxnSigner{
		TronTxnSigner: NewTronMainnetTxnSigner(client, feeInfo),
		tokenMeta:     tokenMeta,
	}
}

func NewTronTokenTestShastaTxnSigner(
	tokenMeta TokenMeta,
	client Client, feeInfo FeeInfo,
) *TronTokenTxnSigner {
	return &TronTokenTxnSigner{
		TronTxnSigner: NewTronTestShastaTxnSigner(client, feeInfo),
		tokenMeta:     tokenMeta,
	}
}

func (this *TronTokenTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	if this.srcAddress == nil {
		return utils.IssueErrorf("TRX txn signer needs source address to sign")
	}
	if this.dstAddress == nil {
		return utils.IssueErrorf("TRX txn signer needs destination address to sign")
	}

	refBlock, err := this.fetchRefBlock()
	if err != nil {
		return
	}
	txn, err := this.newTokenTxn(refBlock)
	if err != nil {
		return
	}
	txnRawData, err := proto.Marshal(txn.GetRawData())
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	txnHash := comutils.HashSha256(txnRawData)

	signature, err := crypto.Sign(txnHash, this.privateKey)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	txn.Signature = append(txn.Signature, signature)

	this.signedTxnHash = txnHash
	this.signedTxn = txn

	return
}

func (this *TronTokenTxnSigner) newTokenTxn(refBlock Block) (_ *core.Transaction, err error) {
	tokenAddress, err := tronaddr.Base58ToAddress(this.tokenMeta.Address)
	if err != nil {
		return
	}

	var contractData bytes.Buffer

	methodID, err := comutils.HexDecode(TronTRC20TransferTokenMethodIdHex)
	comutils.PanicOnError(err)
	if _, err = contractData.Write(methodID); err != nil {
		return
	}

	paddedAddress := common.LeftPadBytes(this.dstAddress.Bytes(), 32)
	if _, err = contractData.Write(paddedAddress); err != nil {
		return
	}

	var selectedAmount decimal.Decimal
	if this.isMove {
		selectedAmount, err = this.getBalance()
		if err != nil {
			return
		} else if selectedAmount.IsZero() {
			err = utils.WrapError(constants.ErrorAmountTooLow)
			return
		}
	} else {
		selectedAmount = this.amount
		if !this.isPreferOffline {
			balance, getErr := this.getBalance()
			if getErr != nil {
				err = getErr
				return
			}
			if selectedAmount.GreaterThan(balance) {
				err = utils.WrapError(constants.ErrorBalanceNotEnough)
				return
			}
		}
	}

	amount := selectedAmount.Shift(int32(this.tokenMeta.DecimalPlaces))
	paddedAmount := common.LeftPadBytes(amount.BigInt().Bytes(), 32)
	if _, err = contractData.Write(paddedAmount); err != nil {
		return
	}

	transfer := core.TriggerSmartContract{
		OwnerAddress:    *this.srcAddress,
		ContractAddress: tokenAddress,
		Data:            contractData.Bytes(),
	}
	transferAny, err := ptypes.MarshalAny(&transfer)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	contract := core.Transaction_Contract{
		Type:      core.Transaction_Contract_TriggerSmartContract,
		Parameter: transferAny,
	}

	txnData := this.newTxnData(refBlock)
	txnData.Contract = []*core.Transaction_Contract{&contract}

	txn := core.Transaction{
		RawData: &txnData,
	}
	return &txn, nil
}
