package ordermod

import (
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

type OrderDB struct {
	*gorm.DB
}

func GetOrderDB(db *gorm.DB) *OrderDB {
	return &OrderDB{db}
}

func (this *OrderDB) GetOrder(ID uint64) (order models.Order, err error) {
	err = this.First(&order, ID).Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func (this *OrderDB) GetCurrencyOrder(currency meta.Currency, ID uint64) (
	order models.Order, err error,
) {
	err = this.
		First(&order, &models.Order{ID: ID, Currency: currency}).
		Error
	if err != nil {
		err = utils.WrapError(err)
	}
	return
}

func (this *OrderDB) FilterUserCurrency(uid meta.UID, currency meta.Currency) *OrderDB {
	return &OrderDB{
		this.
			Where(&models.Order{UID: uid, Currency: currency}).
			Where(dbquery.In(models.OrderColStatus, constants.OrderVisibleStatuses)),
	}
}

func (this *OrderDB) FilterBlockchainTxns() *OrderDB {
	return &OrderDB{
		this.
			Where(&models.Order{
				SrcChannelType: constants.ChannelTypeSrcBlockchainNetwork,
			}).
			Where(dbquery.In(
				models.OrderColDstChannelType, constants.BlockchainChannelTypes,
			)).
			Where(dbquery.NotEqual(models.OrderColDstChannelRef, "")),
	}
}

func (this *OrderDB) FilterBlockchainTxn(currency meta.Currency, hash string) *OrderDB {
	return &OrderDB{
		this.
			FilterBlockchainTxns().
			Where(&models.Order{
				Currency:      currency,
				DstChannelRef: hash,
			}),
	}
}
