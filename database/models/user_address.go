package models

import (
	"database/sql/driver"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserAddressTableName = "user_address"

	UserAddressColUID     = "uid"
	UserAddressColAddress = "address"
)

type UserAddress struct {
	ID      uint64                       `gorm:"column:id;primaryKey"`
	UID     meta.UID                     `gorm:"column:uid"`
	Network meta.BlockchainNetwork       `gorm:"column:network"`
	Address string                       `gorm:"column:address"`
	Key     UserAddressKeyEncryptedField `gorm:"column:key"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
}

func (UserAddress) TableName() string {
	return UserAddressTableName
}

type UserAddressKeyEncryptedField struct {
	dbfields.EncryptedField
}

func (this *UserAddressKeyEncryptedField) SetSecret() error {
	return this.EncryptedField.SetSecretHex(viper.GetString(config.KeyBlockchainSecret))
}

func (this *UserAddressKeyEncryptedField) GetValue() (string, error) {
	if err := this.SetSecret(); err != nil {
		return "", err
	}
	return this.EncryptedField.GetValue()
}

func (this *UserAddressKeyEncryptedField) GetValueF() string {
	value, err := this.GetValue()
	comutils.PanicOnError(err)

	return value
}

func (this UserAddressKeyEncryptedField) Value() (driver.Value, error) {
	if err := this.SetSecret(); err != nil {
		return nil, err
	}
	return this.EncryptedField.Value()
}

func (this *UserAddressKeyEncryptedField) Scan(input interface{}) error {
	if err := this.SetSecret(); err != nil {
		return err
	}
	return this.EncryptedField.Scan(input)
}

func NewUserAddressKeyEncryptedField(key string) UserAddressKeyEncryptedField {
	return UserAddressKeyEncryptedField{dbfields.NewEncryptedField(key)}
}
