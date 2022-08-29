package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/services/transaction"
	"gitlab.com/snap-clickstaff/torque-go/api/services/transaction/crontasks"
	"gitlab.com/snap-clickstaff/torque-go/cmd/helpers"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/payoutmod"
)

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Torque Trading Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(
			comlocale.Init(config.LanguageCode, config.SignalReload))
		comutils.PanicOnError(database.Init())
	},
}

var transactionServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Trading HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		transaction.StartServer(host, port)
	},
}
var transactionCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start Trading Cron service",
	Run: func(cmd *cobra.Command, args []string) {
		transaction.StartCron()
	},
}

var transactionPayoutCmd = &cobra.Command{
	Use:   "payout",
	Short: "Process Payout",
}
var transactionPayoutGenerateStartCoinBalancesCmd = &cobra.Command{
	Use:   "gen_start_balances",
	Short: "Generate start Balances (once a day)",
	Run: func(cmd *cobra.Command, args []string) {
		yesterdayTime := comutils.TimeRoundDate(time.Now().Add(-24 * time.Hour))
		yesterdateDateStr := yesterdayTime.Format(constants.DateFormatISO)

		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = yesterdateDateStr
		}

		dateTime, err := comutils.TimeParse(constants.DateFormatISO, date)
		comutils.PanicOnError(err)
		if dateTime.After(yesterdayTime) {
			panic(fmt.Errorf("can only generate start date balances to date %v (not %v)", yesterdateDateStr, date))
		}

		helpers.TradingPayoutGenUserStartCoinBalances(date)
	},
}
var transactionPayoutGenerateStatsCmd = &cobra.Command{
	Use:   "gen_stats",
	Short: "Generate stats (once a day)",
	Run: func(cmd *cobra.Command, args []string) {
		yesterdayTime := comutils.TimeRoundDate(time.Now().Add(-24 * time.Hour))
		yesterdateDateStr := yesterdayTime.Format(constants.DateFormatISO)

		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = yesterdateDateStr
		}

		dateTime, err := comutils.TimeParse(constants.DateFormatISO, date)
		comutils.PanicOnError(err)
		if dateTime.After(yesterdayTime) {
			panic(fmt.Errorf("can only generate stats to date %v (not %v)", yesterdateDateStr, date))
		}

		db := database.GetDbMaster()
		stats, err := payoutmod.GenPayoutStats(date)
		comutils.PanicOnError(err)
		comutils.PanicOnError(
			db.Save(&stats).Error,
		)

		comutils.EchoWithTime("Payout stats for date %v has been created with id=%v", date, stats.ID)
	},
}

var transactionBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "System Balance",
}
var transactionBalanceGenSystemStartBalancesCmd = &cobra.Command{
	Use:   "gen_system_start_balances",
	Short: "Generate system start balances",
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = time.Now().Format(constants.DateFormatISO)
		}
		useFullMode, _ := cmd.Flags().GetBool("full")
		if useFullMode {
			crontasks.SystemTradingGenStartBalancesFull(date, 0)
		} else {
			crontasks.SystemTradingGenStartBalancesAccumulate(date)
		}
	},
}

func init() {
	rootCmd.AddCommand(transactionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transactionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transactionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	transactionCmd.AddCommand(transactionCronCmd)

	transactionServeCmd.Flags().String("host", "localhost", "Host to bind")
	transactionServeCmd.Flags().Int("port", 8080, "Port to bind")
	transactionCmd.AddCommand(transactionServeCmd)

	transactionPayoutGenerateStartCoinBalancesCmd.Flags().String("date", "", "Date to generate")
	transactionPayoutCmd.AddCommand(transactionPayoutGenerateStartCoinBalancesCmd)
	transactionPayoutGenerateStatsCmd.Flags().StringP("date", "d", "", "Date to generate")
	transactionPayoutCmd.AddCommand(transactionPayoutGenerateStatsCmd)
	transactionCmd.AddCommand(transactionPayoutCmd)

	transactionBalanceGenSystemStartBalancesCmd.Flags().StringP("date", "d", "", "Date to generate")
	transactionBalanceGenSystemStartBalancesCmd.Flags().BoolP("full", "f", false, "Full load mode (very heavy!)")
	transactionBalanceCmd.AddCommand(transactionBalanceGenSystemStartBalancesCmd)
	transactionCmd.AddCommand(transactionBalanceCmd)
}
