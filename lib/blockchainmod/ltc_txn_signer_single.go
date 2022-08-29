package blockchainmod

import "gitlab.com/snap-clickstaff/torque-go/config"

func NewLitecoinSingleTxnSigner(client Client, feeInfo FeeInfo) *BitcoinSingleTxnSigner {
	if config.BlockchainUseTestnet {
		return NewLitecoinTestnetSingleTxnSigner(client, feeInfo)
	} else {
		return NewLitecoinMainnetSingleTxnSigner(client, feeInfo)
	}
}

func NewLitecoinMainnetSingleTxnSigner(client Client, feeInfo FeeInfo) *BitcoinSingleTxnSigner {
	return &BitcoinSingleTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &LitecoinChainConfig,
	}
}

func NewLitecoinTestnetSingleTxnSigner(client Client, feeInfo FeeInfo) *BitcoinSingleTxnSigner {
	return &BitcoinSingleTxnSigner{
		baseTxnSigner: baseTxnSigner{
			client:  client,
			feeInfo: feeInfo,
		},
		chainConfig: &LitecoinChainConfigTestnet,
	}
}
