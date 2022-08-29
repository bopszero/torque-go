package database

import (
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func IsDbError(err error) bool {
	if err == nil {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func IsDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *mysql.MySQLError:
		return err.(*mysql.MySQLError).Number == MySQLErrorCodeDuplicateEntry
	default:
		return false
	}
}
