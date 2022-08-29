package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	ChannelInfoTableName = "channel"
)

type ChannelInfo struct {
	Type        meta.ChannelType `gorm:"column:type;primaryKey"`
	Name        string           `gorm:"column:name"`
	Description string           `gorm:"column:description"`

	MinTxnAmount            decimal.Decimal                    `gorm:"column:min_txn_amount"`
	MaxTxnAmount            decimal.Decimal                    `gorm:"column:max_txn_amount"`
	IsAvailable             sql.NullBool                       `gorm:"column:is_available"`
	BlockchainNetworkConfig ChannelInfoBlockchainNetworkConfig `gorm:"column:blockchain_network_config"`

	CreateTime int64 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime int64 `gorm:"column:update_time;autoUpdateTime"`
}

func (ChannelInfo) TableName() string {
	return ChannelInfoTableName
}

type ChannelInfoBlockchainNetworkConfig struct {
	CurrencyAvailabilityMap    map[meta.Currency]bool                       `json:"currency_availability_map"`
	CurrencyMarkupPriceMap     map[meta.Currency]meta.AmountMarkup          `json:"currency_markup_price_map"`
	CurrencyAmountThresholdMap map[meta.Currency]ChannelInfoAmountThreshold `json:"currency_amount_threshold_map"`
}

type ChannelInfoAmountThreshold struct {
	Min decimal.Decimal `json:"min"`
	Max decimal.Decimal `json:"max"`
}

func (this *ChannelInfoBlockchainNetworkConfig) Scan(value interface{}) error {
	var valueBytes []byte
	switch valueWithType := value.(type) {
	case string:
		valueBytes = []byte(valueWithType)
	case []byte:
		valueBytes = valueWithType
	default:
		return fmt.Errorf("could not convert value '%+v' to byte array", value)
	}

	if len(valueBytes) == 0 {
		return nil
	}

	return json.Unmarshal(valueBytes, this)
}

func (this ChannelInfoBlockchainNetworkConfig) Value() (driver.Value, error) {
	return json.Marshal(this)
}
