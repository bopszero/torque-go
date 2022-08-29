package blockchainmod

import (
	"crypto/ecdsa"
	"encoding/binary"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	tronaddr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type TronTxnSigner struct {
	baseTxnSigner

	privateKey *ecdsa.PrivateKey

	amount     decimal.Decimal
	srcAddress *tronaddr.Address
	dstAddress *tronaddr.Address

	isMove bool

	signedTxnHash []byte
	signedTxn     *core.Transaction
}

func NewTronTxnSigner(client Client, feeInfo FeeInfo) *TronTxnSigner {
	if config.BlockchainUseTestnet {
		return NewTronTestShastaTxnSigner(client, feeInfo)
	} else {
		return NewTronMainnetTxnSigner(client, feeInfo)
	}
}

func NewTronMainnetTxnSigner(client Client, feeInfo FeeInfo) *TronTxnSigner {
	return &TronTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
	}
}

func NewTronTestShastaTxnSigner(client Client, feeInfo FeeInfo) *TronTxnSigner {
	return &TronTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
	}
}

func (this *TronTxnSigner) GetRaw() []byte {
	txnBytes, err := proto.Marshal(this.signedTxn)
	comutils.PanicOnError(err)

	return txnBytes
}

func (this *TronTxnSigner) GetHash() string {
	return comutils.HexEncode(this.signedTxnHash)
}

func (this *TronTxnSigner) GetAmount() decimal.Decimal {
	return this.amount
}

func (this *TronTxnSigner) GetSrcAddress() string {
	return this.srcAddress.String()
}

func (this *TronTxnSigner) GetDstAddress() string {
	return this.dstAddress.String()
}

func (this *TronTxnSigner) GetEstimatedFee() (amount meta.CurrencyAmount, err error) {
	if this.signedTxn == nil {
		err = utils.IssueErrorf("Tron txn must be signed before getting fee")
		return
	}

	amount = meta.CurrencyAmount{
		Currency: constants.CurrencySubTronSun,
		Value:    decimal.NewFromInt(this.signedTxn.GetRawData().FeeLimit),
	}
	return currencymod.ConvertAmount(amount, constants.CurrencyTron)
}

func (this *TronTxnSigner) SetNonce(nonce Nonce) error {
	blockObj := nonce.GetValue()
	if _, ok := blockObj.(Block); !ok {
		return utils.IssueErrorf("Tron nonce must be a Block (not %T)", blockObj)
	}
	return this.baseTxnSigner.SetNonce(nonce)
}

func (this *TronTxnSigner) SetSrc(privateKeyHex string, _ string) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return utils.WrapError(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	srcAddress := tronaddr.PubkeyToAddress(*publicKey)

	this.privateKey = privateKey
	this.srcAddress = &srcAddress

	return nil
}

func (this *TronTxnSigner) SetDst(address string, amount decimal.Decimal) error {
	dstAddress, err := tronaddr.Base58ToAddress(address)
	if err != nil {
		return err
	}

	this.dstAddress = &dstAddress
	this.amount = amount

	return nil
}

func (this *TronTxnSigner) SetMoveDst(address string) error {
	dstAddress, err := tronaddr.Base58ToAddress(address)
	if err != nil {
		return err
	}
	this.dstAddress = &dstAddress
	this.isMove = true

	return nil
}

func (this *TronTxnSigner) Sign(canOverwrite bool) (err error) {
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
	txn, err := this.newNormalTxn(refBlock)
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

func (this *TronTxnSigner) fetchRefBlock() (Block, error) {
	if this.nonce != nil {
		blockObj := this.nonce.GetValue()
		return blockObj.(Block), nil
	}

	block, err := this.client.GetLatestBlock()
	if err != nil {
		return nil, err
	}
	tronNonce := NewTronFrozenBlockNonce(block)
	if err := this.SetNonce(tronNonce); err != nil {
		return nil, err
	}

	return block, nil
}

func (this *TronTxnSigner) newTxnData(refBlock Block) (txnData core.TransactionRaw) {
	var (
		txnTime          = time.Now()
		expireTime       = txnTime.Add(TronTxnLifetime)
		blockNumber      = int64(refBlock.GetHeight())
		blockNumberBytes = make([]byte, binary.MaxVarintLen64)
		blockHashBytes   = comutils.HexDecodeF(refBlock.GetHash())
	)
	binary.BigEndian.PutUint64(blockNumberBytes, uint64(blockNumber))

	txnData = core.TransactionRaw{
		RefBlockNum:   blockNumber,
		RefBlockBytes: blockNumberBytes[6:8],
		RefBlockHash:  blockHashBytes[8:16],

		Timestamp:  tronMakeErrorTimeMs(TimeNsToMs(txnTime.UnixNano()), TronTxnErrorRangeTime),
		Expiration: tronMakeErrorTimeMs(TimeNsToMs(expireTime.UnixNano()), TronTxnErrorRangeExpireTime),
		FeeLimit:   this.feeInfo.LimitMaxValue.IntPart(),
	}
	return
}

func (this *TronTxnSigner) getBalance() (_ decimal.Decimal, err error) {
	if this.srcAddress == nil {
		err = utils.IssueErrorf("TRX txn signer needs source address to get balance")
		return
	}
	return this.client.GetBalance(this.GetSrcAddress())
}

func (this *TronTxnSigner) newNormalTxn(refBlock Block) (_ *core.Transaction, err error) {
	var (
		maxFee         = this.feeInfo.LimitMaxValue.IntPart()
		selectedAmount int64
		leftAmount     int64
	)
	if this.isMove {
		balance, err := this.getBalance()
		if err != nil {
			return nil, err
		}
		leftAmount = tronToSun(balance) - maxFee
		selectedAmount = leftAmount
	} else {
		selectedAmount = tronToSun(this.amount)
		if this.isPreferOffline {
			leftAmount = 1 // Ignore balance validation
		} else {
			balance, err := this.getBalance()
			if err != nil {
				return nil, err
			}
			totalAmount := selectedAmount + maxFee
			leftAmount = tronToSun(balance) - totalAmount
		}
	}
	if leftAmount < 0 {
		return nil, utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	transfer := core.TransferContract{
		OwnerAddress: *this.srcAddress,
		ToAddress:    *this.dstAddress,
		Amount:       selectedAmount,
	}
	transferAny, err := ptypes.MarshalAny(&transfer)
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	contract := core.Transaction_Contract{
		Type:      core.Transaction_Contract_TransferContract,
		Parameter: transferAny,
	}

	txnData := this.newTxnData(refBlock)
	txnData.Contract = []*core.Transaction_Contract{&contract}

	txn := core.Transaction{
		RawData: &txnData,
	}
	return &txn, nil
}
