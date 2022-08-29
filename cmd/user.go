package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/cmd/helpers"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use: "user",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(database.Init())
	},
}

// generateDateTierChangesCmd represents the generateDateTierChanges command
var generateDateTierChangesCmd = &cobra.Command{
	Use:   "gen_date_tier_changes",
	Short: "Generate User tier changes of date",
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		doClearData, _ := cmd.Flags().GetBool("clear-data")
		if date == "" {
			date = time.Now().Format(constants.DateFormatISO)
		}

		helpers.GenerateDateUserTierChanges(date, doClearData)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateDateTierChangesCmd.Flags().String("date", "", "Date to fetch User balances")
	generateDateTierChangesCmd.Flags().Bool("clear-data", false, "Clear date data before generate")
	userCmd.AddCommand(generateDateTierChangesCmd)
}
