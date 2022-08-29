package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/services/sysinternal/crontasks"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/sysforwardingmod"
)

var internalCmd = &cobra.Command{
	Use:   "internal",
	Short: "Torque Internal Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(database.Init())
	},
}

// var internalServeCmd = &cobra.Command{
// 	Use:   "serve",
// 	Short: "Serve Internal HTTP service",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// host, _ := cmd.Flags().GetString("host")
// 		port, _ := cmd.Flags().GetInt("port")
// 		internal.StartServer("127.0.0.1", port)
// 	},
// }
// var internalCronCmd = &cobra.Command{
// 	Use:   "cron",
// 	Short: "Start Internal Cron service",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		internal.StartCron()
// 	},
// }

var internalForwardingCmd = &cobra.Command{
	Use:   "forwarding",
	Short: "System Forwarding",
}
var internalForwardingStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start System Forwarding Service",
	Run: func(cmd *cobra.Command, args []string) {
		currencyVal, _ := cmd.Flags().GetString("currency")
		networkVal, _ := cmd.Flags().GetString("network")

		var (
			currency = meta.NewCurrencyU(currencyVal)
			network  = meta.NewBlockchainNetwork(networkVal)
		)
		if !blockchainmod.IsSupportedIndex(currency, network) {
			comutils.EchoWithTime("Invalid input coin | currency=%v,network=%v", currency, network)
			return
		}
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = time.Now().Format(constants.DateFormatISO)
		}

		if _, err := comutils.TimeParse(constants.DateFormatISO, date); err != nil {
			comutils.EchoWithTime("Invalid input date | date=%v,err=%v", date, err.Error())
			return
		}

		cronSkip := cron.New(
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger),
			),
		)

		cronSkip.AddFunc("@every 5s", func() { crontasks.ForwardingTradingDeposits(currency, network, date) })

		cronSkip.Start()
		comutils.EchoWithTime("Started the %v Forwarding Service for date %v.", currency, date)

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, config.SignalInterrupt)
		signal.Notify(signalChannel, config.SignalTerminate)

		osSignal := <-signalChannel
		comutils.EchoWithTime("Received signal `%s`, stopping the service.", osSignal)
		<-cronSkip.Stop().Done()
		comutils.EchoWithTime("Stopped the %v Forwarding Service for date %v.", currency, date)
	},
}
var internalForwardingImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import System Forwarding addresses",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		fileIO, err := os.Open(filePath)
		comutils.PanicOnError(err)

		ctx := comcontext.NewContext()
		comutils.EchoWithTime("Importing addresses...")
		csvReader := csv.NewReader(fileIO)
		count, err := sysforwardingmod.ImportAddresses(ctx, csvReader)
		if err != nil {
			comutils.EchoWithTime("Import addresses failed | err=%s", err.Error())
		} else {
			comutils.EchoWithTime("Import %v addresses finished.", count)
		}
	},
}
var internalForwardingReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report System Forwarding tranactions",
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = time.Now().Format(constants.DateFormatISO)
		}
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			filePath = fmt.Sprintf("%v.csv", date)
		}

		comutils.EchoWithTime("Generating report...")
		err := sysforwardingmod.GenerateReport(date, filePath)
		if err != nil {
			comutils.EchoWithTime("Generate report failed | err=%s", err.Error())
		} else {
			comutils.EchoWithTime("Generate report finished. View at `%v`.", filePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(internalCmd)
	// internalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	internalForwardingStartCmd.Flags().StringP("currency", "c", "", "Currency to forward")
	internalForwardingStartCmd.MarkFlagRequired("currency")
	internalForwardingStartCmd.Flags().String("network", "", "Network to forward")
	internalForwardingStartCmd.Flags().String("date", "", "Date to forward in YYYY-MM-DD")
	internalForwardingCmd.AddCommand(internalForwardingStartCmd)

	internalForwardingImportCmd.Flags().StringP("file", "f", "", "File to import")
	internalForwardingImportCmd.MarkFlagRequired("file")
	internalForwardingCmd.AddCommand(internalForwardingImportCmd)

	internalForwardingReportCmd.Flags().String("date", "", "Date to forward in YYYY-MM-DD")
	internalForwardingReportCmd.Flags().StringP("file", "f", "", "File to export")
	internalForwardingCmd.AddCommand(internalForwardingReportCmd)

	internalCmd.AddCommand(internalForwardingCmd)
}
