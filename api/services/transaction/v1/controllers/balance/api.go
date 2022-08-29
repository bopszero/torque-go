package balance

import (
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/torque-go/api"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/api/responses"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func GetUserBalance(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel GetBalanceRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	var balances []models.UserBalance
	if reqModel.Currency == "" {
		balances = tradingbalance.GetUserBalances(ctx, reqModel.UID)
	} else {
		balances = []models.UserBalance{
			tradingbalance.GetUserBalance(ctx, reqModel.UID, reqModel.Currency),
		}
	}

	var responseBalances []GetBalanceResponse
	for _, balance := range balances {
		repsonseBalance := GetBalanceResponse{
			UID:        balance.UID,
			Currency:   balance.Currency,
			Amount:     balance.Amount,
			UpdateTime: balance.UpdateTime,
		}
		responseBalances = append(responseBalances, repsonseBalance)
	}

	return responses.Ok(
		ctx,
		meta.O{
			"balances": responseBalances,
		},
	)
}

func AddTxn(c echo.Context) error {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		reqModel AddTxnRequest
	)
	if err := api.BindAndValidate(c, &reqModel); err != nil {
		return err
	}

	txnTypeMeta, ok := constants.TradingBalanceTypeMetaCodeMap[strings.ToLower(reqModel.TypeCode)]
	if !ok {
		return utils.IssueErrorf("unknown transaction type `%s`", reqModel.TypeCode)
	}

	txn, err := tradingbalance.AddTransaction(
		ctx,
		reqModel.Currency, reqModel.UID, reqModel.Amount,
		txnTypeMeta.ID, reqModel.Ref)
	if err != nil {
		return err
	}

	responseTxn := AddTxnResponse{
		ID:       txn.ID,
		Currency: txn.Currency,
		UID:      txn.UserID,
		Amount:   txn.Amount,
		Balance:  txn.Balance,
		TypeCode: constants.TradingBalanceTypeMetaMap[txn.Type].Code,
		Ref:      txn.Ref,
	}
	if txn.ParentID.Valid {
		parentID := uint64(txn.ParentID.Int64)
		responseTxn.ParentID = &parentID
	}

	return responses.Ok(
		ctx,
		meta.O{
			"txn": responseTxn,
		},
	)
}
