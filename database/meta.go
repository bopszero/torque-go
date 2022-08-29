package database

import (
	"time"
)

const (
	AliasDefault    = "default"
	AliasMainMaster = "main.master"
	AliasMainSlave  = "main.slave"

	AliasInternalMaster = "internal.master"

	// AliasAuthMaster = "auth.master"
	// AliasAuthSlave  = "auth.slave"

	AliasWalletMaster = "wallet.master"
	AliasWalletSlave  = "wallet.slave"
)

const (
	DefaultPoolSize        = 5
	DefaultPoolMaxLifetime = 5 * time.Minute
	DefaultPoolMaxSize     = 50
)

const (
	MySQLErrorCodeDuplicateEntry = 1062
)
