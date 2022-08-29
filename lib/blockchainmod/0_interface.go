package blockchainmod

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type (
	Coin interface {
		GetIndexCode() string

		GetCurrency() meta.Currency
		GetTradingID() uint16
		IsUsingMainnet() bool
		IsAvailable() bool
		GetNetwork() meta.BlockchainNetwork
		GetNetworkMain() meta.BlockchainNetwork
		GetNetworkTest() meta.BlockchainNetwork
		GetNetworkCurrency() meta.Currency
		GetDecimalPlaces() uint8

		GetModelInfo() models.CurrencyInfo
		GetModelInfoLegacy() models.LegacyCurrencyInfo
		GetModelNetwork() models.BlockchainNetworkInfo
		GetModelNetworkCurrency() models.NetworkCurrency
		GetDefaultFeeInfo() (FeeInfo, error)
		GetMinTxnAmount() decimal.Decimal
		GenTxnExplorerURL(txnHash string) string
		NormalizeAddress(address string) (string, error)
		NormalizeAddressLegacy(address string) (string, error)

		NewKey() (KeyHolder, error)
		LoadKey(privateKey string, hintAddress string) (KeyHolder, error)
		NewClientDefault() (Client, error)
		NewClientSpare() (Client, error)
		NewAccountGuest(address string) (GuestAccount, error)
		NewAccountOwner(privateKey string, hintAddress string) (OwnerAccount, error)
		NewAccountSystem(comcontext.Context, meta.UID) (SystemAccount, error)
		NewTxnSignerSingle(Client, *FeeInfo) (SingleTxnSigner, error)
		NewTxnSignerBulk(Client, *FeeInfo) (BulkTxnSigner, error)
	}
)

type (
	Block interface {
		GetCurrency() meta.Currency
		GetNetwork() meta.BlockchainNetwork
		GetHash() string
		GetHeight() uint64
		GetTimeUnix() int64
		GetParentHash() string

		GetTransactions() ([]Transaction, error)
	}

	Transaction interface {
		SetOwnerAddress(address string)

		GetCurrency() meta.Currency
		GetTypeCode() string
		GetLocalStatus() meta.BlockchainTxnStatus
		GetDirection() meta.Direction
		GetFee() meta.CurrencyAmount

		GetHash() string
		// GetBlockHash() string
		GetBlockHeight() uint64
		GetConfirmations() uint64

		GetFromAddress() (string, error)
		GetToAddress() (string, error)
		GetAmount() (decimal.Decimal, error)

		GetInputDataHex() (string, error)
		GetTimeUnix() int64

		GetInputs() ([]Input, error)
		GetOutputs() ([]Output, error)
		GetRC20Transfers(TokenMeta) ([]RC20Transfer, error)
	}

	Input interface {
		GetPrevOutHash() string
		GetPrevOutIndex() uint32

		GetPrevOutAddress() string
		GetPrevOutAmount() decimal.Decimal
	}

	Output interface {
		GetIndex() uint32
		GetAddress() string
		GetAmount() decimal.Decimal
	}

	UnspentTxnOutput interface {
		GetAddress() string
		GetTxnHash() string

		GetIndex() uint32
		GetAmount() decimal.Decimal
	}

	RC20Transfer interface {
		GetTokenMeta() TokenMeta
		GetFromAddress() string
		GetToAddress() string
		GetAmount() decimal.Decimal
	}

	Nonce interface {
		GetNumber() (uint64, error)
		GetValue() interface{}
		Next() (Nonce, error)
	}
)

type Client interface {
	GetBlock(blockHash string) (Block, error)
	GetBlockByHeight(height uint64) (Block, error)
	GetLatestBlock() (Block, error)

	GetBalance(address string) (decimal.Decimal, error)
	GetTxn(hash string) (Transaction, error)
	GetTxns(address string, paging meta.Paging) ([]Transaction, error)
	GetNextNonce(address string) (Nonce, error)
	GetUtxOutputs(address string, minAmount decimal.Decimal) ([]UnspentTxnOutput, error)

	PushTxnRaw(data []byte) error
}

type (
	GuestAccount interface {
		GetCurrency() meta.Currency
		GetAddress() string
		GetClient() Client

		GetBalance() (decimal.Decimal, error)
		GetTxn(hash string) (Transaction, error)
		GetTxns(paging meta.Paging) ([]Transaction, error)
		GetFeeInfoToAddress(toAddress string) (FeeInfo, error)

		PushTxn(data []byte) error
	}

	OwnerAccount interface {
		GuestAccount

		GetPublicKey() string
		GetPrivateKey() string

		GenerateSignedTxn(
			address string, amount decimal.Decimal, offerFeeInfo *FeeInfo,
		) (SingleTxnSigner, error)
	}

	SystemAccount interface {
		OwnerAccount

		GetUID() meta.UID
	}
)

type (
	KeyHolder interface {
		GetPrivateKey() string
		GetPublicKey() string
		GetAddress() string
	}
)

type (
	TxnSigner interface {
		SetClient(Client)

		GetRaw() []byte
		GetHash() string
		GetAmount() decimal.Decimal
		GetFeeInfo() FeeInfo
		GetEstimatedFee() (meta.CurrencyAmount, error)

		MarkAsPreferOffline()
		SetMoveDst(address string) error
		Sign(canOverwrite bool) error
	}

	BulkTxnSigner interface {
		TxnSigner

		SetLeftoverAddress(address string) error
		AddSrc(privateKey string, hintAddress string) error
		AddDst(address string, amount decimal.Decimal) error
	}

	SingleTxnSigner interface {
		TxnSigner

		GetSrcAddress() string
		GetDstAddress() string

		SetNonce(Nonce) error
		SetSrc(privateKey string, hintAddress string) error
		SetDst(address string, amount decimal.Decimal) error
	}
)
