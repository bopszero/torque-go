package tradingtxn

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/currencymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func TransferP2P(
	ctx comcontext.Context,
	fromUID meta.UID, fromCurrency meta.Currency, fromAmount decimal.Decimal,
	toUID meta.UID, toCurrency meta.Currency, exchangeRate decimal.Decimal,
	feeAmount decimal.Decimal, note string,
) (p2pTransfer models.P2PTransfer, err error) {
	fromCurrencyInfo, err := currencymod.GetCurrencyInfoFast(fromCurrency)
	if err != nil || !currencymod.IsValidTradingInfo(fromCurrencyInfo) {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}
	toCurrencyInfo, err := currencymod.GetCurrencyInfoFast(toCurrency)
	if err != nil || !currencymod.IsValidTradingInfo(toCurrencyInfo) {
		err = utils.WrapError(constants.ErrorCurrency)
		return
	}

	now := time.Now()
	p2pTransfer = models.P2PTransfer{
		FromUID:      fromUID,
		FromCurrency: fromCurrency,
		FromAmount:   utils.NormalizeTradingAmount(fromAmount),
		ToUID:        toUID,
		ToCurrency:   toCurrency,
		ToAmount:     utils.NormalizeTradingAmount(fromAmount.Mul(exchangeRate)),
		ExchangeRate: exchangeRate,
		FeeAmount:    utils.NormalizeTradingAmount(feeAmount),
		Note:         note,
		Status:       constants.CommonStatusActive,

		CreateTime: now.Unix(),
		UpdateTime: now.Unix(),
	}
	err = database.TransactionFromContext(ctx, database.AliasMainMaster, func(dbTxn *gorm.DB) error {
		comutils.PanicOnError(
			dbTxn.Save(&p2pTransfer).Error,
		)

		totalFromAmount := p2pTransfer.FromAmount.Add(feeAmount)
		_, fromErr := tradingbalance.AddTransaction(
			ctx,
			fromCurrency,
			p2pTransfer.FromUID,
			totalFromAmount.Mul(constants.DecimalOneNegative),
			constants.TradingBalanceTypeMetaDrTransfer.ID,
			comutils.Stringify(p2pTransfer.ID),
		)
		if fromErr != nil {
			return fromErr
		}

		_, toErr := tradingbalance.AddTransaction(
			ctx,
			toCurrency,
			p2pTransfer.ToUID,
			p2pTransfer.ToAmount,
			constants.TradingBalanceTypeMetaCrTransfer.ID,
			comutils.Stringify(p2pTransfer.ID),
		)
		if toErr != nil {
			return toErr
		}

		// TODO: Add fee amount into somewhere like system balance.

		return nil
	})
	return
}
