package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/cmd/helpers"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

var blockchainCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain Service",
}

var blockchainAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account",
}
var blockchainAccountGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate accounts address and keys.",
	Run: func(cmd *cobra.Command, args []string) {
		currencyStr, _ := cmd.Flags().GetString("currency")
		networkStr, _ := cmd.Flags().GetString("network")
		output, _ := cmd.Flags().GetString("output")

		quantity, err := cmd.Flags().GetInt("quantity")
		comutils.PanicOnError(err)

		var (
			currency = meta.NewCurrencyU(currencyStr)
			network  = meta.NewBlockchainNetwork(networkStr)
		)

		comutils.EchoWithTime("generate blockchain accounts | %s | started.", currency)

		keyHolders := helpers.GenerateBlockchainKeyHolders(currency, network, quantity)
		if output == "" {
			for _, keyHolder := range keyHolders {
				fmt.Printf("%16s: %s\n", "Address", keyHolder.GetAddress())
				fmt.Printf("%16s: %s\n", "Public Key", keyHolder.GetPublicKey())
				fmt.Printf("%16s: %s\n", "Private Key", keyHolder.GetPrivateKey())
				fmt.Printf("--------------------------------\n")
			}
		} else {
			csvFile, err := os.Create(output)
			comutils.PanicOnError(err)

			csvWriter := csv.NewWriter(bufio.NewWriter(csvFile))
			csvWriter.Write([]string{"address", "public key", "private key"})
			for _, keyHolder := range keyHolders {
				err := csvWriter.Write([]string{
					keyHolder.GetAddress(),
					keyHolder.GetPublicKey(),
					keyHolder.GetPrivateKey(),
				})
				comutils.PanicOnError(err)
			}

			csvWriter.Flush()
			comutils.PanicOnError(csvFile.Close())
		}

		comutils.EchoWithTime("generate blockchain accounts | %s | finished.", currency)
	},
}

var blockchainRippleCmd = &cobra.Command{
	Use:   "ripple",
	Short: "Ripple related commands.",
}
var blockchainRippleGenTagCmd = &cobra.Command{
	Use:   "gen-tag",
	Short: "Generate tagged X-Addresses from a legacy Ripple address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		isTestnet, err := cmd.Flags().GetBool("testnet")
		comutils.PanicOnError(err)
		firstTag, err := cmd.Flags().GetUint32("first-tag")
		comutils.PanicOnError(err)
		quantity, err := cmd.Flags().GetUint32("quantity")
		comutils.PanicOnError(err)
		outputPath, _ := cmd.Flags().GetString("output")

		lastTag := firstTag + quantity - 1
		if lastTag > math.MaxUint32 {
			comutils.EchoWithTime("the last tag is too big `%d`", firstTag+quantity)
			return
		}
		xAddress, err := blockchainmod.RippleParseXAddress(address, !isTestnet)
		if err != nil {
			comutils.EchoWithTime("cannot parse Ripple address `%s`", address)
			return
		}
		if hasTag, _ := xAddress.GetTag(); hasTag {
			comutils.EchoWithTime(
				"Cannot use the tagged Ripple address `%s` to generate tag",
				xAddress.GetRootTagAddress())
			return
		}

		comutils.EchoWithTime("generate tagged X-Addresses for address `%s`...", address)
		taggedXAddresses := make([]blockchainmod.RippleXAddress, 0, quantity)

		var (
			rootAddress = xAddress.GetRootAddress()
			iterTag     = firstTag
		)
		for ; iterTag <= lastTag; iterTag++ {
			newXAddress, err := blockchainmod.NewRippleXAddress(rootAddress, &iterTag, !isTestnet)
			comutils.PanicOnError(err)
			taggedXAddresses = append(taggedXAddresses, newXAddress)
		}

		if outputPath == "" {
			fmt.Printf("Base address: %s\n", xAddress.GetRootAddress())
			fmt.Println("--------------------------------")
			for _, xAddress := range taggedXAddresses {
				_, tag := xAddress.GetTag()
				fmt.Printf("%10s: %10d\n", "Tag", tag)
				fmt.Printf("%10s: %32s\n", "X-Address", xAddress.GetAddress())
				fmt.Println("--------------------------------")
			}
		} else {
			csvFile, err := os.Create(outputPath)
			comutils.PanicOnError(err)
			defer func() { comutils.PanicOnError(csvFile.Close()) }()

			csvWriter := csv.NewWriter(bufio.NewWriter(csvFile))
			csvWriter.Write([]string{"root_address", "tag", "x_address"})
			for _, xAddress := range taggedXAddresses {
				_, tag := xAddress.GetTag()
				err := csvWriter.Write([]string{
					xAddress.GetRootAddress(),
					comutils.Stringify(tag),
					xAddress.GetAddress(),
				})
				comutils.PanicOnError(err)
			}

			csvWriter.Flush()
		}

		comutils.EchoWithTime("generate tagged X-Addresses for address `%s` finished", address)
	},
}

func init() {
	rootCmd.AddCommand(blockchainCmd)

	blockchainAccountGenerateCmd.Flags().StringP("currency", "c", "", "Currency (required)")
	blockchainAccountGenerateCmd.MarkFlagRequired("currency")
	blockchainAccountGenerateCmd.Flags().String("network", "", "Network")
	blockchainAccountGenerateCmd.Flags().IntP("quantity", "n", 1, "Account quantity")
	blockchainAccountGenerateCmd.Flags().StringP("output", "o", "", "Output file. Leave empty for StdIn.")
	blockchainAccountCmd.AddCommand(blockchainAccountGenerateCmd)
	blockchainCmd.AddCommand(blockchainAccountCmd)

	blockchainRippleGenTagCmd.Flags().BoolP("testnet", "t", false, "Use testnet mode")
	blockchainRippleGenTagCmd.Flags().Uint32P("first-tag", "f", 0, "First tag")
	blockchainRippleGenTagCmd.Flags().Uint32P("quantity", "n", 1, "Tag quantity")
	blockchainRippleGenTagCmd.Flags().StringP("output", "o", "", "Output file. Leave empty for StdIn.")
	blockchainRippleCmd.AddCommand(blockchainRippleGenTagCmd)
	blockchainCmd.AddCommand(blockchainRippleCmd)
}
