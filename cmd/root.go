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
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

var (
	configFilePath string
	logFilePath    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "torque-go",
	Short: "Torque Extension Service - written in Go",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	defer config.CmdExecuteRootDefers()
	defer func() {
		if config.Debug {
			return
		}
		if err := recover(); err != nil {
			sentry.CaptureException(comutils.ToError(err))
			panic(err)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(func() {
		if logFilePath == "" {
			logFilePath = viper.GetString(config.KeyLogFile)
		}
		comutils.PanicOnError(
			comlogging.Init(logFilePath, viper.GetBool(config.KeyDebug)),
		)
	})
	cobra.OnInitialize(func() {
		comutils.PanicOnError(
			comcache.Init(config.KeyCacheMap, &comcache.MsgPackEncoder{}),
		)
	})
	cobra.OnInitialize(func() {
		config.CmdRegisterRootDefer(ants.Release)
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
		&configFilePath, "config", "",
		"config file (default is ./config.yaml)",
	)
	rootCmd.PersistentFlags().StringVar(
		&logFilePath, "log", "",
		"log file (default is in config file)",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	config.InitWithPath(configFilePath)
	config.InitSentry()
}
