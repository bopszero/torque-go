package blockchainmod

import (
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

var (
	globalConfig = comtypes.NewSingleton(func() interface{} {
		configMap := viper.GetStringMap(config.KeyBlockchainConfig)

		var conf Config
		comutils.PanicOnError(
			utils.DumpDataByJSON(&configMap, &conf),
		)
		if conf.EthereumContractDefaultGasLimit < EthereumStandardGasLimit {
			panic(utils.IssueErrorf(
				"Ethereum contract default gas limit %v is too low (cannot be lower than %v)",
				conf.EthereumContractDefaultGasLimit, EthereumStandardGasLimit,
			))
		}
		if conf.TetherUsdEthereumGasLimit < EthereumTokenUsdtGasLimitMin {
			panic(utils.IssueErrorf(
				"Ethereum Tether USD gas limit %v is too low (cannot be lower than %v)",
				conf.TetherUsdEthereumGasLimit, EthereumTokenUsdtGasLimitMin,
			))
		}

		return &conf
	})
)

type Config struct {
	EthereumContractDefaultGasLimit uint32 `json:"eth_contract_default_gas_limit"`
	TetherUsdEthereumGasLimit       uint32 `json:"usdt_eth_gas_limit"`
	TetherUsdTronEnergyLimit        uint32 `json:"usdt_trx_energy_limit"`
}

func GetConfig() *Config {
	return globalConfig.Get().(*Config)
}
