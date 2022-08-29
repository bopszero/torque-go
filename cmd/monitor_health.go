/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	HEALTH_URL             = "https://torquebot.net/health/"
	HEALTH_REQUEST_TIMEOUT = "2s"
)

// monitorHealthCmd represents the monitor_health command
var monitorHealthCmd = &cobra.Command{
	Use:   "monitor_health",
	Short: "Monitor 'torquebot.net' site health",
	Long: `Monitor 'torquebot.net' site health
by sending '/health/' requests in a specific interval.`,
	Run: func(cmd *cobra.Command, args []string) {
		interval, err := cmd.Flags().GetString("interval")
		comutils.PanicOnError(err)

		intervalDuration, err := time.ParseDuration(interval)
		comutils.PanicOnError(err)

		fmt.Printf("Start monitoring health with interval %s\n", interval)
		execute(intervalDuration)
	},
}

func init() {
	rootCmd.AddCommand(monitorHealthCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorHealthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorHealthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	monitorHealthCmd.Flags().String("interval", "30s", "Interval to check")
}

func execute(interval time.Duration) {
	logger := comlogging.GetLogger()

	reqTimeout, _ := time.ParseDuration(HEALTH_REQUEST_TIMEOUT)
	client := resty.New().SetTimeout(reqTimeout)

	for {
		var notifyErr error

		respone, err := client.R().Get(HEALTH_URL)
		if err != nil {
			notifyErr = notifyError(err)

		} else if respone.StatusCode() != http.StatusOK {
			notifyErr = notifyUnexpectedResponse(respone)
		}

		if notifyErr != nil {
			logger.WithError(err).Error("Cannot send Telegram notification.")
		}

		time.Sleep(interval)
	}
}

func notifyError(err error) error {
	return notifyViaTeleBot(fmt.Sprintf("[ERROR] %s", err))
}

func notifyUnexpectedResponse(response *resty.Response) error {
	return notifyViaTeleBot(fmt.Sprintf("[RESPONSE] %s", response.Status()))
}

type myRecipient struct{}

func (r *myRecipient) Recipient() string {
	return viper.GetString(config.KeyTeleBotHealthRecipientID)
}

func notifyViaTeleBot(statusDesc string) error {
	now := time.Now()
	message := fmt.Sprintf(
		"%s | Torque health issue: %s",
		now.Format(time.RFC3339), statusDesc,
	)

	if config.Debug {
		fmt.Println(message)
		return nil
	}

	teleBot, err := telebot.NewBot(telebot.Settings{
		Token: viper.GetString(config.KeyTeleBotHealthToken),
	})
	if err != nil {
		return err
	}

	_, err = teleBot.Send(&myRecipient{}, message)
	return err
}
