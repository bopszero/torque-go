package models

import (
	"database/sql/driver"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database/dbfields"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	SystemWithdrawalAddressTableName = "system_withdrawal_address"
)

type SystemWithdrawalAddress struct {
	ID uint32 `gorm:"column:id;primaryKey"`

	Network  meta.BlockchainNetwork                   `gorm:"column:network"`
	Currency meta.Currency                            `gorm:"column:currency"`
	Address  string                                   `gorm:"column:address"`
	Key      SystemWithdrawalAddressKeyEncryptedField `gorm:"column:key"`
	Status   int8                                     `gorm:"column:status"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (SystemWithdrawalAddress) TableName() string {
	return SystemWithdrawalAddressTableName
}

type SystemWithdrawalAddressKeyEncryptedField struct {
	dbfields.EncryptedField
}

func (this *SystemWithdrawalAddressKeyEncryptedField) SetSecret() error {
	return this.EncryptedField.SetSecretHex(viper.GetString(config.KeySystemWithdrawalSecret))
}

func (this SystemWithdrawalAddressKeyEncryptedField) GetValue() (string, error) {
	if err := this.SetSecret(); err != nil {
		return "", err
	}
	return this.EncryptedField.GetValue()
}

func (this SystemWithdrawalAddressKeyEncryptedField) Value() (driver.Value, error) {
	if err := this.SetSecret(); err != nil {
		return nil, err
	}
	return this.EncryptedField.Value()
}

func (this *SystemWithdrawalAddressKeyEncryptedField) Scan(input interface{}) error {
	if err := this.SetSecret(); err != nil {
		return err
	}
	return this.EncryptedField.Scan(input)
}

func NewSystemWithdrawalAddressKeyEncryptedField(key string) SystemWithdrawalAddressKeyEncryptedField {
	return SystemWithdrawalAddressKeyEncryptedField{dbfields.NewEncryptedField(key)}
}
