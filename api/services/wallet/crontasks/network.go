package crontasks

import (
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func NetworkUpdateAllLatestBlockHeights() {
	defer apiutils.CronRunWithRecovery("NetworkUpdateAllLatestBlockHeights", nil)

	logCurrencyErr := func(network meta.BlockchainNetwork, err error) {
		comlogging.GetLogger().
			WithError(err).
			WithField("network", network).
			Error("update network latest block height failed")
	}

	var networks []models.BlockchainNetworkInfo
	comutils.PanicOnError(
		database.GetDbF(database.AliasWalletSlave).
			Find(&networks).
			Error,
	)
	for _, net := range networks {
		coin, err := blockchainmod.GetNetworkMainCoin(net.Network)
		if err != nil {
			logCurrencyErr(net.Network, err)
			continue
		}
		client, err := coin.NewClientDefault()
		if err != nil {
			logCurrencyErr(net.Network, err)
			continue
		}
		block, err := client.GetLatestBlock()
		if err != nil {
			logCurrencyErr(net.Network, err)
			continue
		}
		if err := blockchainmod.UpdateCurrencyLatestBlockHeight(coin, block.GetHeight()); err != nil {
			logCurrencyErr(net.Network, err)
			continue
		}
	}
}
