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
	"gitlab.com/snap-clickstaff/torque-go/api/services/general"
	"gitlab.com/snap-clickstaff/torque-go/database"
)

// serveGeneralCmd represents the serve_general command
var serveGeneralCmd = &cobra.Command{
	Use:   "serve_general",
	Short: "Serve General HTTP service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		comutils.PanicOnError(database.Init())
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		general.StartServer(host, port)
	},
}

func init() {
	rootCmd.AddCommand(serveGeneralCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveGeneralCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveGeneralCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveGeneralCmd.Flags().String("host", "localhost", "Host to bind")
	serveGeneralCmd.Flags().Int("port", 8080, "Port to bind")
}
