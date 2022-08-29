package depositcrawler

import "gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"

const (
	ConfirmationsUpdateMax = 30
)

type Crawler interface {
	ConsumeBlock(blockchainmod.Block) error
	ConsumeTxn(blockchainmod.Transaction) error

	GetScanBlockHeights() ([]uint64, error)
}
