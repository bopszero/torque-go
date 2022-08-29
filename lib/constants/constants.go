package constants

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

const (
	AmountMaxDecimalPlaces        = 18
	AmountTradingMaxDecimalPlaces = 8
)

const (
	CommonStatusActive   = 1
	CommonStatusInactive = -1

	CommonStatusCodeActive   = "active"
	CommonStatusCodeInactive = "inactive"
)

const (
	DateFormatISO      = "2006-01-02"
	DateFormatSlug     = "20060102"
	DateTimeFormatISO  = "2006-01-02 15:04:05"
	DateTimeFormatSlug = "20060102150405"
)

var (
	DecimalOneNegative = decimal.NewFromInt(-1)
	DecimalOne         = decimal.NewFromInt(1)
	DecimalTen         = decimal.NewFromInt(10)
)

const (
	DepositStatusInQueue              = "In Queue"
	DepositStatusPendingConfirmations = "Under Processing"
	DepositStatusPendingReinvest      = "Pending"
	DepositStatusApproved             = "Approved"
	DepositStatusRejected             = "Rejected"
)

const (
	DirectionTypeUnknown = meta.Direction(0)
	DirectionTypeSend    = meta.Direction(-1)
	DirectionTypeReceive = meta.Direction(1)
)

const (
	JudgeStatusPending  = "Pending"
	JudgeStatusApproved = "Approved"
	JudgeStatusRejected = "Rejected"
)

const (
	PromoCodeStatusAvailable = "ACTIVE"
	PromoCodeStatusUsed      = "USED"
)

const (
	SettingKeyBlackListIP                   = "ip:black_list"
	SettingKeyCurrencyNoticePattern         = "currency:notice:%v"
	SettingKeyIsEnableBanIP                 = "ip:is_enable_ban"
	SettingKeyKycJumioWhiteListIP           = "kyc:jumio:white_list_ip"
	SettingKeyKycLaunchTime                 = "kyc:launch_time"
	SettingKeyKycReminderDeadlineTime       = "kyc:reminder:deadline_time"
	SettingKeyKycReminderIsEnabled          = "kyc:reminder:is_enabled"
	SettingKeyKycSubmitLimit                = "kyc:submit_limit"
	SettingKeyKycWhiteListUID               = "kyc:white_list_uids"
	SettingKeyMetaFeatures                  = "meta:features"
	SettingKeySystemWithdrawalClientModeMap = "system:withdrawal:client_mode_map"
	SettingKeyTorqueTransferFee             = "torque_fee_amount"
	SettingKeyWhiteListIP                   = "ip:white_list"
)

const (
	ToggleOff     = -1
	ToggleOn      = 1
	ToggleCodeOn  = "1"
	ToggleCodeOff = "0"
)

const (
	WithdrawStatusPendingConfirm  = "Pending"
	WithdrawStatusCanceled        = "Canceled"
	WithdrawStatusPendingTransfer = "Processing"
	WithdrawStatusApproved        = "Approved"
	WithdrawStatusRejected        = "Rejected"
)

const (
	LockActionSendPersonal      = 1
	LockActionTORQTransfer      = 2
	LockActionTORQReallocate    = 3
	LockActionTORQPurchase      = 4
	LockActionConvertCoinToTORQ = 6
	LockActionTradingSend       = 5
)
