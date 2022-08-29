package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/services/wallet"
	"gitlab.com/snap-clickstaff/torque-go/cmd/helpers"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Torque Wallet Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(
			comlocale.Init(config.LanguageCode, config.SignalReload))
		comutils.PanicOnError(database.Init())
	},
}

var walletServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Wallet HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		wallet.StartServer(host, port)
	},
}
var walletCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start Wallet Cron service",
	Run: func(cmd *cobra.Command, args []string) {
		wallet.StartCron()
	},
}
var walletConsumersCmd = &cobra.Command{
	Use:   "consumers",
	Short: "Start Wallet Consumers service",
	Run: func(cmd *cobra.Command, args []string) {
		wallet.StartConsumers()
	},
}

var walletPayoutCmd = &cobra.Command{
	Use:   "payout",
	Short: "Wallet Payout commands",
}
var walletPayoutCollectRewardsCmd = &cobra.Command{
	Use:   "collect_rewards",
	Short: "Collect Rewards (once a day)",
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			yesterdayTime := time.Now().Add(-24 * time.Hour)
			date = yesterdayTime.Format(constants.DateFormatISO)
		}

		helpers.WalletCollectRewards(date)
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	walletCmd.AddCommand(walletCronCmd)
	walletCmd.AddCommand(walletConsumersCmd)

	walletServeCmd.Flags().String("host", "localhost", "Host to bind")
	walletServeCmd.Flags().Int("port", 8080, "Port to bind")
	walletCmd.AddCommand(walletServeCmd)

	walletPayoutCollectRewardsCmd.Flags().String("date", "", "Date to collect")
	walletPayoutCmd.AddCommand(walletPayoutCollectRewardsCmd)
	walletCmd.AddCommand(walletPayoutCmd)
}
