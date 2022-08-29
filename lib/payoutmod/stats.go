package payoutmod

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GenPayoutStats(dateStr string) (stats models.PayoutStats, err error) {
	if _, err = comutils.TimeParse(constants.DateFormatISO, dateStr); err != nil {
		err = utils.WrapError(err)
		return
	}
	stats.Date = dateStr

	db := database.GetDbSlave()

	profitQuery := db.
		Model(&models.DailyProfit{}).
		Select(dbquery.Sum(models.DailyProfitColAmount)).
		Where(&models.DailyProfit{Date: dateStr, IsDeleted: models.NewBool(false)})
	profitRow := profitQuery.Row()
	comutils.PanicOnError(
		profitQuery.Error,
	)
	comutils.PanicOnError(
		profitRow.Scan(&stats.TotalProfitAmount),
	)

	affiliateQuery := db.
		Model(&models.AffiliateCommission{}).
		Select(dbquery.Sum(models.AffiliateCommissionColAmount)).
		Where(&models.AffiliateCommission{Date: dateStr, IsDeleted: models.NewBool(false)})
	affiliateRow := affiliateQuery.Row()
	comutils.PanicOnError(
		affiliateQuery.Error,
	)
	var actualAffiliateAmount decimal.Decimal
	comutils.PanicOnError(
		affiliateRow.Scan(&actualAffiliateAmount),
	)
	systemPayableAffliliateAmount := stats.TotalProfitAmount.Mul(decimal.NewFromFloat(DailyProfitAffiliateRate))
	stats.RemainingAmountAffiliate = systemPayableAffliliateAmount.Sub(actualAffiliateAmount)
	if stats.RemainingAmountAffiliate.IsNegative() {
		err = utils.IssueErrorf(
			"payout affiliate amount is greater than system payable amount | payable=%v,paid=%v",
			systemPayableAffliliateAmount, actualAffiliateAmount,
		)
		return
	}

	leaderRewardQuery := db.
		Model(&models.LeaderReward{}).
		Select(dbquery.Sum(models.LeaderRewardColAmount)).
		Where(&models.LeaderReward{Date: dateStr, IsDeleted: models.NewBool(false)})
	leaderRewardRow := leaderRewardQuery.Row()
	comutils.PanicOnError(
		leaderRewardQuery.Error,
	)
	var actualLeaderRewardAmount decimal.Decimal
	comutils.PanicOnError(
		leaderRewardRow.Scan(&actualLeaderRewardAmount),
	)
	systemPayableLeaderRewardAmount := stats.TotalProfitAmount.Mul(decimal.NewFromFloat(DailyProfitLeaderRewardRate))
	stats.RemainingAmountLeaderReward = systemPayableLeaderRewardAmount.Sub(actualLeaderRewardAmount)
	if stats.RemainingAmountLeaderReward.IsNegative() {
		err = utils.IssueErrorf(
			"payout leader reward amount is greater than system payable amount | payable=%v,paid=%v",
			systemPayableLeaderRewardAmount, actualLeaderRewardAmount,
		)
		return
	}

	stats.CreateTime = time.Now().Unix()
	return
}
