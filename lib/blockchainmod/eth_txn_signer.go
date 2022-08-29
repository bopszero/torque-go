package blockchainmod

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type EthereumTxnSigner struct {
	baseTxnSigner

	chainID    *big.Int
	privateKey *ecdsa.PrivateKey

	amount     decimal.Decimal
	srcAddress *common.Address
	dstAddress *common.Address

	isMove bool

	signedTxn *types.Transaction
}

func NewEthereumTxnSigner(client Client, feeInfo FeeInfo) *EthereumTxnSigner {
	if config.BlockchainUseTestnet {
		return NewEthereumTestRopstenTxnSigner(client, feeInfo)
	} else {
		return NewEthereumMainnetTxnSigner(client, feeInfo)
	}
}

func NewEthereumMainnetTxnSigner(client Client, feeInfo FeeInfo) *EthereumTxnSigner {
	return &EthereumTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainID: ChainIdEthereum,
	}
}

func NewEthereumTestRopstenTxnSigner(client Client, feeInfo FeeInfo) *EthereumTxnSigner {
	return &EthereumTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainID: ChainIdEthereumTestRopsten,
	}
}

func (this *EthereumTxnSigner) GetRaw() []byte {
	var txnBuffer bytes.Buffer
	comutils.PanicOnError(
		this.signedTxn.EncodeRLP(&txnBuffer),
	)

	return txnBuffer.Bytes()
}

func (this *EthereumTxnSigner) GetRawHex() string {
	return comutils.HexEncode(this.GetRaw())
}

func (this *EthereumTxnSigner) GetHash() string {
	return this.signedTxn.Hash().Hex()
}

func (this *EthereumTxnSigner) GetAmount() decimal.Decimal {
	return this.amount
}

func (this *EthereumTxnSigner) GetSrcAddress() string {
	return strings.ToLower(this.srcAddress.Hex())
}

func (this *EthereumTxnSigner) GetDstAddress() string {
	return strings.ToLower(this.dstAddress.Hex())
}

func (this *EthereumTxnSigner) GetEstimatedFee() (meta.CurrencyAmount, error) {
	feeAmount := meta.CurrencyAmount{
		Currency: constants.CurrencyEthereum,
		Value:    this.feeInfo.GetBaseLimitMaxValue(),
	}
	return feeAmount, nil
}

func (this *EthereumTxnSigner) SetNonce(nonce Nonce) error {
	if _, err := nonce.GetNumber(); err != nil {
		return err
	}
	return this.baseTxnSigner.SetNonce(nonce)
}

func (this *EthereumTxnSigner) SetSrc(privateKey string, _ string) error {
	privateKey = comutils.HexTrim(privateKey)

	privateKeyEcdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return utils.WrapError(err)
	}

	this.privateKey = privateKeyEcdsa

	publicKey := privateKeyEcdsa.Public().(*ecdsa.PublicKey)
	srcAddress := crypto.PubkeyToAddress(*publicKey)

	this.srcAddress = &srcAddress
	return nil
}

func (this *EthereumTxnSigner) SetDst(address string, amount decimal.Decimal) error {
	dstAddress, err := hexToEthereumAddress(address)
	if err != nil {
		return err
	}

	this.dstAddress = &dstAddress
	this.amount = amount

	return nil
}

func (this *EthereumTxnSigner) SetMoveDst(address string) error {
	dstAddress, err := hexToEthereumAddress(address)
	if err != nil {
		return err
	}
	this.dstAddress = &dstAddress
	this.isMove = true

	return nil
}

func (this *EthereumTxnSigner) Sign(canOverwrite bool) (err error) {
	if !canOverwrite && this.signedTxn != nil {
		return nil
	}

	if this.srcAddress == nil {
		return utils.IssueErrorf("ETH txn signer needs source address to sign")
	}
	if this.dstAddress == nil {
		return utils.IssueErrorf("ETH txn signer needs destination address to sign")
	}

	nonce, err := this.loadNonce()
	if err != nil {
		return
	}
	txn, err := this.newNormalTxn(nonce)
	if err != nil {
		return
	}
	signedTxn, err := types.SignTx(txn, types.NewEIP155Signer(this.chainID), this.privateKey)
	if err != nil {
		return
	}

	this.signedTxn = signedTxn
	return
}

func (this *EthereumTxnSigner) loadNonce() (_ uint64, err error) {
	if this.nonce != nil {
		return this.nonce.GetNumber()
	}

	nonce, err := this.client.GetNextNonce(this.GetSrcAddress())
	if err != nil {
		return
	}
	if err = this.SetNonce(nonce); err != nil {
		return
	}

	return this.nonce.GetNumber()
}

func (this *EthereumTxnSigner) GetSignedTxn() *types.Transaction {
	return this.signedTxn
}

func (this *EthereumTxnSigner) getBalance() (_ decimal.Decimal, err error) {
	if this.srcAddress == nil {
		err = utils.IssueErrorf("ETH txn signer needs source address to get balance")
		return
	}
	return this.client.GetBalance(this.GetSrcAddress())
}

func (this *EthereumTxnSigner) newNormalTxn(nonce uint64) (*types.Transaction, error) {
	var (
		gasPrice     = EthereumToWei(this.baseTxnSigner.feeInfo.GetBasePrice())
		gasLimit     = big.NewInt(int64(this.baseTxnSigner.feeInfo.LimitMaxQuantity))
		maxGasAmount = new(big.Int).Mul(gasPrice, gasLimit)

		selectedAmount *big.Int
		leftAmount     *big.Int
	)
	if this.isMove {
		balance, err := this.getBalance()
		if err != nil {
			return nil, err
		}
		leftAmount = new(big.Int).Sub(EthereumToWei(balance), maxGasAmount)
		selectedAmount = leftAmount
	} else {
		selectedAmount = EthereumToWei(this.amount)
		if this.isPreferOffline {
			leftAmount = big.NewInt(1) // Ignore balance validation
		} else {
			balance, err := this.getBalance()
			if err != nil {
				return nil, err
			}
			totalAmount := new(big.Int).Add(selectedAmount, maxGasAmount)
			leftAmount = new(big.Int).Sub(EthereumToWei(balance), totalAmount)
		}
	}
	if leftAmount.Cmp(big.NewInt(0)) < 0 {
		return nil, utils.WrapError(constants.ErrorBalanceNotEnough)
	}

	txn := types.NewTransaction(
		nonce, *this.dstAddress, selectedAmount,
		uint64(this.baseTxnSigner.feeInfo.LimitMaxQuantity), gasPrice,
		nil)
	return txn, nil
}
