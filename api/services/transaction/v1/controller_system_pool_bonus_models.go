package v1

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/bonuspoolmod"
)

type (
	PoolBonusLeaderCheckoutRequest struct {
		FromDate string `json:"from_date" validate:"required,datetime=2006-01-02"`
		ToDate   string `json:"to_date" validate:"required,datetime=2006-01-02"`
	}
	PoolBonusLeaderCheckoutResponse struct {
		ExecutionHash string                        `json:"execution_hash"`
		FromDate      string                        `json:"from_date"`
		ToDate        string                        `json:"to_date"`
		TotalAmount   decimal.Decimal               `json:"total_amount"`
		TierInfoList  []bonuspoolmod.LeaderTierInfo `json:"tier_info_list"`
	}
)

type (
	PoolBonusLeaderExecuteRequest struct {
		ExecutionHash string                        `json:"execution_hash" validate:"required"`
		FromDate      string                        `json:"from_date" validate:"required,datetime=2006-01-02"`
		ToDate        string                        `json:"to_date" validate:"required,datetime=2006-01-02"`
		TotalAmount   decimal.Decimal               `json:"total_amount" validate:"required"`
		TierInfoList  []bonuspoolmod.LeaderTierInfo `json:"tier_info_list" validate:"required,min=1,dive"`
	}
	PoolBonusLeaderExecuteResponse struct {
		Execution *models.LeaderBonusPoolExecution `json:"execution"`
	}
)
