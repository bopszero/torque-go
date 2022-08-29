package config

import (
	"syscall"
)

const (
	EnvLive    = "live"
	EnvStaging = "staging"
	EnvTest    = "test"
	EnvDev     = "dev"
)

const (
	CacheAliasDefault = "default"
	CacheAliasMemory  = "memory"
	CacheAliasRemote  = "remote"
)

const (
	SignalInterrupt = syscall.Signal(0x2) // os.Interrupt
	SignalReload    = syscall.Signal(0xa) // unix.SIGUSR1
	SignalTerminate = syscall.Signal(0xf) // unix.SIGTERM
)

const (
	KeyAppName                = "APP_NAME"
	KeySecretKey              = "SECRET_KEY"
	KeyBlockchainSecret       = "BLOCKCHAIN_SECRET"
	KeySystemWithdrawalSecret = "SYSTEM_WITHDRAWAL_SECRET"

	KeyRegion       = "REGION"
	KeyLanguageCode = "LANGUAGE_CODE"
	KeyTimeZone     = "TIME_ZONE"
	KeyCurrency     = "CURRENCY"

	KeyEnv   = "ENV"
	KeyDebug = "DEBUG"
	KeyTest  = "TEST"

	KeyServerGracefulTimeout = "SERVER_GRACEFUL_TIMEOUT"

	KeyDbConnectionMap = "DB_CONNECTION_MAP"
	KeyDbPoolOptions   = "DB_POOL_OPTIONS"
	KeyCacheMap        = "CACHE_MAP"

	KeyAuthLegacySecret   = "AUTH_SECRET" // Deprecated
	KeyAuthAccessSecret   = "AUTH_ACCESS_SECRET"
	KeyAuthAccessTimeout  = "AUTH_ACCESS_TIMEOUT"
	KeyAuthRefreshSecret  = "AUTH_REFRESH_SECRET"
	KeyAuthRefreshTimeout = "AUTH_REFRESH_TIMEOUT"

	KeyLogFile     = "LOG_FILE"
	KeySentryDSN   = "SENTRY_DSN"
	KeyMsgQueueDSN = "MSG_QUEUE_DSN"
	KeyRedLockDSN  = "REDLOCK_DSN"
	KeyJumioDSN    = "JUMIO_DSN"
	KeyEmailDSN    = "EMAIL_DSN"

	KeyTeleBotHealthToken       = "TELE_BOT_HEALTH_TOKEN"
	KeyTeleBotHealthRecipientID = "TELE_BOT_HEALTH_RECIPIENT_ID"
	KeyTreeRefreshDuration      = "TREE_REFRESH_DURATION"
	KeyBlockchainUseTestnet     = "BLOCKCHAIN_USE_TESTNET"
	KeyBlockchainConfig         = "BLOCKCHAIN_CONFIG"
	KeyBonusPoolLeaderConfig    = "BONUS_POOL_LEADER_CONFIG"

	KeyTorquePurchaseCompanyAddressMap = "TORQUE_PURCHASE_COMPANY_ADDRESS_MAP"
	KeySystemForwardingConfigMap       = "SYSTEM_FORWARDING_CONFIG_MAP"
	KeySystemWithdrawalConfigMap       = "SYSTEM_WITHDRAWAL_CONFIG_MAP"

	KeyServiceWebBaseURL           = "SERVICE_WEB_URI"
	KeyServiceWebMacSecret         = "SERVICE_WEB_MAC_SECRET"
	KeyServiceTradingSecret        = "SERVICE_TRADING_MAC_SECRET"
	KeyServiceTreeBaseURL          = "SERVICE_TREE_BASE_URL"
	KeyServiceTransactionMacSecret = "SERVICE_TRANSACTION_MAC_SECRET"
	KeyServiceWalletMacSecret      = "SERVICE_WALLET_MAC_SECRET"
	KeyServiceAuthMacSecret        = "SERVICE_AUTH_MAC_SECRET"
	KeyServicePushBaseURL          = "SERVICE_PUSH_URI"
	KeyServicePushMacSecret        = "SERVICE_PUSH_MAC_SECRET"

	KeyApiSnapNodeKey       = "API_SNAP_NODE_KEY"
	KeyApiBlockchainInfoKey = "API_BLOCKCHAIN_INFO_KEY"
	KeyApiEtherscanKey      = "API_ETHERSCAN_KEY"
	KeyApiBlockCypherKey    = "API_BLOCK_CYPHER_KEY"
	KeyApiCryptoCompareKey  = "API_CRYPTO_COMPARE_KEY"
	KeyApiBinanceKey        = "API_BINANCE_KEY"
	KeyApiBinanceSecret     = "API_BINANCE_SECRET"

	KeyGeoIPDBPath = "GEO_IP_DB_PATH"
)

const (
	ContextKeyRequest   = "request"
	ContextKeyUser      = "user"
	ContextKeyUID       = "uid"
	ContextKeyLocalizer = "localizer"
	ContextKeyJWT       = "jwt"

	ContextKeyDebugUID = "debug:uid"
)

const (
	HttpHeaderAcceptLanguage = "Accept-Language"
	HttpHeaderUserAgent      = "User-Agent"
	HttpHeaderCacheControl   = "Cache-Control"
	HttpHeaderCfConnectingIP = "Cf-Connecting-Ip"
	HttpHeaderXForwardedFor  = "X-Forwarded-For"

	HttpHeaderTorqueLanguage  = "Torque-Language"
	HttpHeaderApiResponseCode = "API-Response-Code"
)
