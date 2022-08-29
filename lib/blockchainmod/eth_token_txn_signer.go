package blockchainmod

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type EthereumTokenTxnSigner struct {
	*EthereumTxnSigner

	tokenMeta TokenMeta
}

func NewEthereumTokenMainnetTxnSigner(
	tokenMeta TokenMeta,
	client Client, feeInfo FeeInfo,
) *EthereumTokenTxnSigner {
	return &EthereumTokenTxnSigner{
		EthereumTxnSigner: NewEthereumMainnetTxnSigner(client, feeInfo),
		tokenMeta:         tokenMeta,
	}
}

func NewEthereumTokenTestRopstenTxnSigner(
	tokenMeta TokenMeta,
	client Client, feeInfo FeeInfo,
) *EthereumTokenTxnSigner {
	return &EthereumTokenTxnSigner{
		EthereumTxnSigner: NewEthereumTestRopstenTxnSigner(client, feeInfo),
		tokenMeta:         tokenMeta,
	}
}

func (this *EthereumTokenTxnSigner) Sign(canOverwrite bool) (err error) {
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
	txn, err := this.newTokenTxn(nonce)
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

func (this *EthereumTokenTxnSigner) newTokenTxn(nonce uint64) (_ *types.Transaction, err error) {
	tokenAddress, err := hexToEthereumAddress(this.tokenMeta.Address)
	if err != nil {
		return
	}

	// Contract address: 0x722dd3F80BAC40c951b51BdD28Dd19d435762180
	// Function call: showMeTheMoney(address,uint256)

	var dataBuffer bytes.Buffer

	// transferFuncSignature := []byte("redeem(uint256)")
	// hash := sha3.NewLegacyKeccak256()
	// hash.Write(transferFuncSignature)
	// methodID := hash.Sum(nil)[:4]
	// _, err = dataBuffer.Write(methodID)
	// if err != nil {
	// 	return
	// }

	methodID, err := comutils.HexDecode0x(EthereumERC20TransferTokenMethodIdHex)
	comutils.PanicOnError(err)
	if _, err = dataBuffer.Write(methodID); err != nil {
		return
	}

	paddedAddress := common.LeftPadBytes(this.dstAddress.Bytes(), 32)
	if _, err = dataBuffer.Write(paddedAddress); err != nil {
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
			balance, err := this.getBalance()
			if err != nil {
				return nil, err
			}
			if selectedAmount.GreaterThan(balance) {
				return nil, utils.WrapError(constants.ErrorBalanceNotEnough)
			}
		}
	}

	amount := selectedAmount.Shift(int32(this.tokenMeta.DecimalPlaces))
	paddedAmount := common.LeftPadBytes(amount.BigInt().Bytes(), 32)
	if _, err = dataBuffer.Write(paddedAmount); err != nil {
		return
	}

	ethAmountZero := big.NewInt(0)
	gasPrice := GweiToWei(this.baseTxnSigner.feeInfo.Price)
	txn := types.NewTransaction(
		nonce,
		tokenAddress, ethAmountZero,
		uint64(this.baseTxnSigner.feeInfo.LimitMaxQuantity), gasPrice,
		dataBuffer.Bytes(),
	)

	return txn, nil
}
