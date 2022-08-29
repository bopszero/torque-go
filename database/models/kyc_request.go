package models

import (
	"time"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	KycRequestTableName     = "kyc_request"
	KycRequestColEmail      = "email_original"
	KycRequestColStatus     = "status"
	KycRequestColCreateTime = "create_time"
)

type KycRequest struct {
	ID                 uint64                `gorm:"column:id;primaryKey"`
	Code               string                `gorm:"column:code"`
	UID                meta.UID              `gorm:"column:uid"`
	UserCode           string                `gorm:"column:user_code"`
	EmailOriginal      string                `gorm:"column:email_original"`
	FullName           string                `gorm:"column:full_name"`
	DOB                time.Time             `gorm:"column:dob"`
	Reference          string                `gorm:"column:reference"`
	Status             meta.KycRequestStatus `gorm:"column:status"`
	CloseUID           int64                 `gorm:"column:close_uid"`
	CloseTime          int64                 `gorm:"column:close_time"`
	Note               string                `gorm:"column:note"`
	NoteInternal       string                `gorm:"column:note_internal"`
	VerificationStatus string                `gorm:"column:verification_status"`
	AttemptWeight      int64                 `gorm:"column:attempt_weight"`
	Nationality        string                `gorm:"column:nationality"`
	Address            string                `gorm:"column:address"`
	CreateTime         int64                 `gorm:"column:create_time;autoCreateTime"`
	UpdateTime         int64                 `gorm:"column:update_time;autoUpdateTime"`
}

func (KycRequest) TableName() string {
	return KycRequestTableName
}
