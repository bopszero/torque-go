package withdrawalmod

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/msgqueuemod"
	"gitlab.com/snap-clickstaff/torque-go/lib/thirdpartymod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type WithdrawNotifParams struct {
	BalanceTxnID uint64 `json:"balance_txn_id"`
}

func init() {
	comutils.PanicOnError(
		msgqueuemod.RegisterHandler(
			msgqueuemod.MessageTypeTradingNotifWithdrawRejected,
			PushWithdrawRejectedHandler,
		),
	)
}

func PushWithdrawRejectedNotificationAsync(balanceTxn models.UserBalanceTxn) error {
	msg := msgqueuemod.NewMessageJsonF(
		msgqueuemod.MessageTypeTradingNotifWithdrawRejected,
		WithdrawNotifParams{BalanceTxnID: balanceTxn.ID},
	)
	queue, err := msgqueuemod.GetQueueWallet()
	if err != nil {
		return err
	}
	return msgqueuemod.PublishMessage(queue, msg)
}

func PushWithdrawRejectedHandler(msg msgqueuemod.Message) (err error) {
	var params WithdrawNotifParams
	if err := comutils.JsonDecode(msg.Data.(string), &params); err != nil {
		return utils.WrapError(err)
	}

	var (
		ctx = comcontext.NewContext()
		db  = database.GetDbSlave()
	)

	var balanceTxn models.UserBalanceTxn
	if err := db.First(&balanceTxn, &models.UserBalanceTxn{ID: params.BalanceTxnID}).Error; err != nil {
		return err
	}

	notifTitle, err := comlocale.TranslateKeyData(
		ctx,
		constants.TranslationKeyTradingNotifTitleWithdrawRejected,
		meta.O{"txn": balanceTxn})
	if err != nil {
		return err
	}

	client := thirdpartymod.GetPushServiceSystemClient()
	messageData := thirdpartymod.PushServiceMessageData{
		Title:   notifTitle,
		Message: "",

		Action:            thirdpartymod.ServicePushActionActivity,
		ActionDestination: thirdpartymod.ServicePushActionDestinationTradingTxn,
		ActionData: comutils.JsonEncodeF(meta.O{
			"id":   params.BalanceTxnID,
			"type": constants.LegacyBalanceTxnTypeCodeNameMap[balanceTxn.Type],
		}),
	}
	err = client.Push(ctx, balanceTxn.UserID, messageData)
	if err != nil {
		return
	}

	return nil
}
