package helpers

import (
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func GenerateBlockchainKeyHolders(
	currency meta.Currency, network meta.BlockchainNetwork,
	quantity int,
) []blockchainmod.KeyHolder {
	var (
		coin     = blockchainmod.GetCoinF(currency, network)
		accounts = make([]blockchainmod.KeyHolder, 0, quantity)
	)
	for i := 0; i < quantity; i++ {
		keyHolder, err := coin.NewKey()
		comutils.PanicOnError(err)

		accounts = append(accounts, keyHolder)
	}

	return accounts
}
