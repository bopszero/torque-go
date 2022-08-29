package balancemod

import (
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gorm.io/gorm"
)

type BalanceDB struct {
	*gorm.DB
}

func GetBalanceDB(db *gorm.DB) *BalanceDB {
	return &BalanceDB{db}
}

func (this *BalanceDB) FilterProfitReinvest(torqueTxnID uint64) *BalanceDB {
	return &BalanceDB{
		this.
			Model(&models.TorqueTxn{}).
			Where(
				&models.TorqueTxn{
					ID:         torqueTxnID,
					IsReinvest: models.NewBool(true),
					IsDeleted:  models.NewBool(false),
				},
			),
	}
}

func (this *BalanceDB) FilterProfitWithdraw(torqueTxnID uint64) *BalanceDB {
	return &BalanceDB{
		this.
			Model(&models.TorqueTxn{}).
			Where(
				&models.TorqueTxn{
					ID:         torqueTxnID,
					IsReinvest: models.NewBool(false),
					IsDeleted:  models.NewBool(false),
				},
			),
	}
}
