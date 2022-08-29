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
	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/cmd/helpers"
)

// commonCmd represents the common command
var commonCmd = &cobra.Command{
	Use:   "common",
	Short: "Common commands",
}
var commonDownloadTranslationsCmd = &cobra.Command{
	Use:   "download_translations",
	Short: "Download Wallet Translation files",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			authFilePath, _ = cmd.Flags().GetString("auth-file")
			sheetID         = "1_kpKTW1q7zGiOGf8dqSwNhlywO-pr05O7Thd34fyOBo"
			readRange       = "torque-wallet!A:Z"
		)
		comutils.PanicOnError(
			helpers.CommonDownloadTranslations("./locale/", authFilePath, sheetID, readRange),
		)
	},
}

func init() {
	rootCmd.AddCommand(commonCmd)

	commonDownloadTranslationsCmd.Flags().StringP("auth-file", "f", "", "JSON auth file path")
	commonDownloadTranslationsCmd.MarkFlagRequired("auth-file")
	commonCmd.AddCommand(commonDownloadTranslationsCmd)
}
