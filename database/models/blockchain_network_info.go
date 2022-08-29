package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	BlockchainNetworkInfoTableName = "blockchain_network"

	BlockchainNetworkInfoColLatestBlockHeight = "latest_block_height"
)

type BlockchainNetworkInfo struct {
	ID uint16 `gorm:"column:id;primaryKey"`

	Network                 meta.BlockchainNetwork `gorm:"column:code"`
	Currency                meta.Currency          `gorm:"column:currency"`
	Name                    string                 `gorm:"column:name"`
	TokenTransferCodeName   string                 `gorm:"column:token_transfer_code_name"`
	LatestBlockHeight       uint64                 `gorm:"column:latest_block_height"`
	DepositMinConfirmations uint64                 `gorm:"column:deposit_min_confirmations"`

	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (BlockchainNetworkInfo) TableName() string {
	return BlockchainNetworkInfoTableName
}
