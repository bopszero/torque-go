package transaction

import (
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/services/transaction/crontasks"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func StartCron() {
	cronSkip := cron.New(
		cron.WithLogger(cron.DefaultLogger),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)
	cronFree := cron.New(
		cron.WithLogger(cron.DefaultLogger),
	)

	cronSkip.AddFunc("@every 30s", crontasks.DepositCollectAndApprove)

	var (
		coinBitcoin     = blockchainmod.GetCoinNativeF(constants.CurrencyBitcoin)
		coinBitcoinCash = blockchainmod.GetCoinNativeF(constants.CurrencyBitcoinCash)
		coinLitecoin    = blockchainmod.GetCoinNativeF(constants.CurrencyLitecoin)
		coinEthereum    = blockchainmod.GetCoinNativeF(constants.CurrencyEthereum)
		coinTron        = blockchainmod.GetCoinNativeF(constants.CurrencyTron)
		coinRipple      = blockchainmod.GetCoinNativeF(constants.CurrencyRipple)
		coinUsdtERC20   = blockchainmod.GetCoinF(constants.CurrencyTetherUSD, blockchainmod.GetSystemNetworkEthereum())
		// coinUsdtTRC20   = blockchainmod.GetCoinF(constants.CurrencyTetherUSD, blockchainmod.GetSystemNetworkTron())
	)
	cronSkip.AddFunc("@every 1m", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinBitcoin) })
	cronSkip.AddFunc("@every 1m", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinBitcoinCash) })
	cronSkip.AddFunc("@every 1m", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinLitecoin) })
	cronSkip.AddFunc("@every 15s", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinEthereum) })
	cronSkip.AddFunc("@every 10s", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinTron) })
	cronSkip.AddFunc("@every 10s", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinRipple) })
	cronSkip.AddFunc("@every 15s", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinUsdtERC20) })
	// cronSkip.AddFunc("@every 15s", func() { crontasks.SystemWithdrawalExecuteCurrencyRequest(coinUsdtTRC20) })

	// Disable due to Trading team loss issue
	// TODO: Enable when everything goes to a normal state
	// cronSkip.AddFunc("@every 2m", func() { crontasks.DepositCrawBlocks(constants.CurrencyBitcoinCash) })
	// cronSkip.AddFunc("@every 2m", func() { crontasks.DepositCrawBlocks(constants.CurrencyBitcoinCashABC) })
	// cronSkip.AddFunc("@every 1m", func() { crontasks.DepositCrawBlocks(constants.CurrencyEthereum) })
	// cronSkip.AddFunc("@every 1m", func() { crontasks.DepositCrawBlocks(constants.CurrencyRipple) })
	// if !config.Test {
	// 	cronSkip.AddFunc("@every 2m", func() { crontasks.DepositCrawBlocks(constants.CurrencyBitcoin) })
	// 	cronSkip.AddFunc("@every 1m", func() { crontasks.DepositCrawBlocks(constants.CurrencyLitecoin) })
	// 	cronSkip.AddFunc("@every 1m", func() { crontasks.DepositCrawBlocks(constants.CurrencyTron) })
	// }

	cronSkip.Start()
	cronFree.Start()
	comutils.EchoWithTime("Started the Cron Service.")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, config.SignalInterrupt)

	osSignal := <-signalChannel
	comutils.EchoWithTime("Received signal `%s`, stopping the service.", osSignal)
	<-cronSkip.Stop().Done()
	<-cronFree.Stop().Done()
	comutils.EchoWithTime("Stopped the Cron Service.")
}
