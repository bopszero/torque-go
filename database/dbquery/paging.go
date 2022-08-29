package dbquery

import (
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gorm.io/gorm"
)

func Paging(db *gorm.DB, paging meta.Paging) *gorm.DB {
	if paging.Limit > 0 {
		db = db.Limit(int(paging.Limit))
	}

	if paging.Offset > 0 {
		db = db.Offset(int(paging.Offset))
	} else if paging.BeforeID > 0 {
		db = db.Where(Lt(models.CommonColID, paging.BeforeID))
	} else if paging.AfterID > 0 {
		db = db.Where(Gt(models.CommonColID, paging.BeforeID))
	}

	return db
}
