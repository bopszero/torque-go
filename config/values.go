package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	AppName   string
	SecretKey string

	Region       string
	LanguageCode string
	TimeZone     string
	Currency     string
	BaseDomain   string

	Env   string
	Debug bool
	Test  bool

	BlockchainUseTestnet bool
)

func getViperNonEmptyString(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("config `%s` cannot be empty", key))
	}

	return value
}

func initValues() {
	AppName = getViperNonEmptyString(KeyAppName)
	SecretKey = getViperNonEmptyString(KeySecretKey)

	Region = getViperNonEmptyString(KeyRegion)
	LanguageCode = getViperNonEmptyString(KeyLanguageCode)
	TimeZone = getViperNonEmptyString(KeyTimeZone)
	Currency = getViperNonEmptyString(KeyCurrency)

	Env = getViperNonEmptyString(KeyEnv)
	Debug = viper.GetBool(KeyDebug)
	Test = viper.GetBool(KeyTest)

	BlockchainUseTestnet = viper.GetBool(KeyBlockchainUseTestnet)
}
