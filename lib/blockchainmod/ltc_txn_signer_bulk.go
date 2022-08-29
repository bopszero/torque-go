package blockchainmod

import "gitlab.com/snap-clickstaff/torque-go/config"

func NewLitecoinBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	if config.BlockchainUseTestnet {
		return NewLitecoinTestnetBulkTxnSigner(client, feeInfo)
	} else {
		return NewLitecoinMainnetBulkTxnSigner(client, feeInfo)
	}
}

func NewLitecoinMainnetBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	return &BitcoinBulkTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &LitecoinChainConfig,
	}
}

func NewLitecoinTestnetBulkTxnSigner(client Client, feeInfo FeeInfo) *BitcoinBulkTxnSigner {
	return &BitcoinBulkTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &LitecoinChainConfigTestnet,
	}
}
