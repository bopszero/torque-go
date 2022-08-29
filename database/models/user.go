package models

import (
	"database/sql"
	"time"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	UserTableName = "users"

	UserColID            = "user_id"
	UserColUsername      = "username"
	UserColEmail         = "email"
	UserColEmailOriginal = "original_email"
	UserColCreateDate    = "date_created"
)

type User struct {
	ID             meta.UID `gorm:"column:user_id;primaryKey"`
	Code           string   `gorm:"column:user_code"`
	Username       string   `gorm:"column:username"`
	Password       string   `gorm:"column:password" json:"-"`
	FirstName      string   `gorm:"column:first_name"`
	LastName       string   `gorm:"column:last_name"`
	Email          string   `gorm:"column:email"`
	OriginalEmail  string   `gorm:"column:original_email"`
	ReferralCode   string   `gorm:"column:code"`
	TierType       int      `gorm:"column:type"`
	TierTypeStatus int      `gorm:"column:type_status"`
	TwoFaKey       string   `gorm:"column:tf_auth_key"`
	ParentID       meta.UID `gorm:"column:referred_by"`
	Status         string   `gorm:"column:status"`
	Nation         string   `gorm:"column:nation"`

	IsDeleted  sql.NullBool `gorm:"column:deleted"`
	CreateDate time.Time    `gorm:"column:date_created"`
}

func (User) TableName() string {
	return UserTableName
}
