package cmd

import (
	"gitlab.com/snap-clickstaff/go-common/comlocale"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/services/auth"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/database"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Torque Auth Service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(
			comlocale.Init(config.LanguageCode, config.SignalReload))
		comutils.PanicOnError(database.Init())
	},
}

var authServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Auth HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		auth.StartServer(host, port)
	},
}
var authCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start Auth Cron service",
	Run: func(cmd *cobra.Command, args []string) {
		auth.StartCron()
	},
}

var authGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Start Auth generate",
}

var kycDOBNotMatchCsvCmd = &cobra.Command{
	Use:   "dob_not_match_csv",
	Short: "Generate file csv has kyc dob not match between user enter and jumio scan result",
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		auth.KycDOBNotMatchCsv(output)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	authCmd.AddCommand(authCronCmd)

	authCmd.AddCommand(authGenerateCmd)
	authGenerateCmd.AddCommand(kycDOBNotMatchCsvCmd)
	kycDOBNotMatchCsvCmd.Flags().StringP("output", "o", "", "Output file. Leave empty for StdIn.")

	authServeCmd.Flags().String("host", "localhost", "Host to bind")
	authServeCmd.Flags().Int("port", 8080, "Port to bind")
	authCmd.AddCommand(authServeCmd)
}
