package bonuspoolmod

import (
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

type LeaderDB struct {
	*gorm.DB
}

func LeaderGetDB(db *gorm.DB) *LeaderDB {
	return &LeaderDB{db}
}

func (this *LeaderDB) GetAndLockExecution(ID uint32, status int8) (exec models.LeaderBonusPoolExecution, err error) {
	query := dbquery.SelectForUpdate(this.DB)
	if status > 0 {
		query = query.Where(&models.LeaderBonusPoolExecution{Status: status})
	}
	err = query.First(&exec, ID).Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func (this *LeaderDB) GetExecutionDetails(executionID uint32) (details []models.LeaderBonusPoolDetail, err error) {
	err = this.
		Where(&models.LeaderBonusPoolDetail{ExecutionID: executionID}).
		Find(&details).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func (this *LeaderDB) GetDetail(ID uint64) (detail models.LeaderBonusPoolDetail, err error) {
	err = this.
		First(&detail, ID).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}
