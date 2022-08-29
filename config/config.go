package config

import (
	"os"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comutils"
)

func Init() {
	InitWithPath("")
}

func InitWithPath(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else if envConfigFilePath := os.Getenv("CONFIG_FILE"); envConfigFilePath != "" {
		viper.SetConfigFile(envConfigFilePath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		comutils.EchoWithTime("Config file: %v\n", viper.ConfigFileUsed())
		setDefaults()
		fillEnvVars()
		initValues()

	} else {
		comutils.EchoWithTime("No config file is used. Error: %s.\n", err.Error())
	}
}

func setDefaults() {
	viper.SetDefault(KeyEnv, EnvDev)
	viper.SetDefault(KeyDebug, true)
	viper.SetDefault(KeyTest, true)

	viper.SetDefault(KeyServerGracefulTimeout, 15)
}

func fillEnvVars() {
	if tz := viper.GetString(KeyTimeZone); tz != "" {
		comutils.PanicOnError(
			os.Setenv("TZ", tz),
		)

	}
}
