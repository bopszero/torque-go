package models

import "gitlab.com/snap-clickstaff/torque-go/lib/meta"

const (
	SystemForwardingAddressTableName = "forward_address"

	SystemForwardingAddressColUID     = "uid"
	SystemForwardingAddressColAddress = "address"
)

type SystemForwardingAddress struct {
	ID      uint64                       `gorm:"column:id;primaryKey"`
	Network meta.BlockchainNetwork       `gorm:"column:network"`
	Address string                       `gorm:"column:address"`
	Key     UserAddressKeyEncryptedField `gorm:"column:key"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
}

func (SystemForwardingAddress) TableName() string {
	return SystemForwardingAddressTableName
}
