package models

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	DepositUserAddressTableName = "deposit_user_address"
)

type DepositUserAddress struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	UID      meta.UID               `gorm:"column:uid"`
	Currency meta.Currency          `gorm:"column:currency"`
	Network  meta.BlockchainNetwork `gorm:"column:network"`
	Address  string                 `gorm:"column:address"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
}

func (DepositUserAddress) TableName() string {
	return DepositUserAddressTableName
}
