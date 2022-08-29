package models

const (
	JumioScanTableName = "jumio_scan"
	JumioScanColStatus = "status"
)

type JumioScan struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	Reference   string `gorm:"column:scan_reference"`
	RequestCode string `gorm:"column:kyc_request_code"`
	UserCode    string `gorm:"column:user_code"`
	Status      int64  `gorm:"column:status"`
	CreateTime  int64  `gorm:"column:create_time;autoCreateTime"`
	UpdateTime  int64  `gorm:"column:update_time;autoUpdateTime"`
}

func (JumioScan) TableName() string {
	return JumioScanTableName
}
