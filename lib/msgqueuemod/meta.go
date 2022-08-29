package msgqueuemod

type (
	MessageType    int
	MessageHandler func(Message) error
)

const (
	MessageTypeOrderNotifCompleted = MessageType(1)
	MessageTypeOrderNotifFailed    = MessageType(2)
	MessageTypeKycSendRequestEmail = MessageType(10)

	MessageTypeTradingNotifDepositApproved  = MessageType(-1)
	MessageTypeTradingNotifWithdrawRejected = MessageType(-2)
)

const (
	QueueKeyAuth   = "auth"
	QueueKeyWallet = "wallet"
)

type Message struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}
