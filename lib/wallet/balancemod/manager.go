package balancemod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func AddP2pTransfer(
	ctx comcontext.Context,
	currency meta.Currency, amount decimal.Decimal,
	fromUID meta.UID, toUID meta.UID,
	feeAmount decimal.Decimal, note string,
) (p2pTransfer models.P2PTransfer, err error) {
	now := time.Now()
	amount = currencymod.NormalizeAmount(currency, amount)
	p2pTransfer = models.P2PTransfer{
		FromUID:      fromUID,
		FromCurrency: currency,
		FromAmount:   amount,
		ToUID:        toUID,
		ToCurrency:   currency,
		ToAmount:     amount,
		ExchangeRate: constants.DecimalOne,
		FeeAmount:    utils.NormalizeTradingAmount(feeAmount),
		Note:         note,
		Status:       constants.CommonStatusActive,

		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		// TODO: Add fee amount into somewhere like system balance.
		return dbTxn.Save(&p2pTransfer).Error
	})
	return
}

func RedeemPromoCode(ctx comcontext.Context, uid meta.UID, code string) (*models.PromoCodeRedemption, error) {
	var (
		promoCode  models.PromoCode
		redemption models.PromoCodeRedemption
	)
	err := database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbquery.SelectForUpdate(dbTxn).
			First(
				&promoCode,
				&models.PromoCode{
					Code:      code,
					Status:    constants.PromoCodeStatusAvailable,
					IsDeleted: models.NewBool(false),
				},
			).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = constants.ErrorDataNotFound
			} else {
				err = constants.ErrorSystem
			}

			return err
		}

		err = dbTxn.
			First(&redemption, &models.PromoCodeRedemption{UID: uid, IsDeleted: models.NewBool(false)}).
			Error
		if database.IsDbError(err) {
			return err
		}
		if redemption.ID != 0 {
			return constants.ErrorDataExists
		}

		portfolioValueUsd := tradingbalance.CalcUserCoinBalanceUsdValue(uid)
		if portfolioValueUsd.LessThan(PromoCodeMinBalanceThreshold) {
			return constants.ErrorAmountTooLowWithValue.WithData(meta.O{
				"threshold": PromoCodeMinBalanceThreshold,
				"currency":  constants.CurrencyUSD,
			})
		}

		now := time.Now()
		redemption = models.PromoCodeRedemption{
			PromoCodeID: promoCode.ID,
			UID:         uid,
			CreditValue: PromoCodeCreditValue,
			IsDeleted:   models.NewBool(false),

			CreateDate: now,
			UpdateDate: now,
		}
		if err = dbTxn.Create(&redemption).Error; err != nil {
			return err
		}

		promoCode.Status = constants.PromoCodeStatusUsed
		promoCode.UpdateDate = now

		return dbTxn.Save(&promoCode).Error
	})
	if err != nil {
		return nil, err
	}

	return &redemption, nil
}
