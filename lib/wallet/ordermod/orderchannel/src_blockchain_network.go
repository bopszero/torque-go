package orderchannel

import (
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

func init() {
	comutils.PanicOnError(
		ordermod.ChannelRegister(&SrcBlockchainNetworkChannel{}),
	)
}

type SrcBlockchainNetworkChannel struct {
	baseChannel
}

func (this *SrcBlockchainNetworkChannel) GetType() meta.ChannelType {
	return constants.ChannelTypeSrcBlockchainNetwork
}

func (this *SrcBlockchainNetworkChannel) PreValidate(ctx comcontext.Context, order *models.Order) error {
	coin, err := blockchainmod.GetCoinNative(order.Currency)
	if err != nil {
		return err
	}
	blockchainAccount, err := coin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return err
	}
	balance, err := blockchainAccount.GetBalance()
	if err != nil {
		return err
	}

	if balance.LessThan(order.AmountTotal) {
		return constants.ErrorBalanceNotEnough
	}

	if order.Currency == constants.CurrencyTetherUSD {
		return this.validateUsdtBalance(ctx, order)
	}

	return nil
}

func (this *SrcBlockchainNetworkChannel) validateUsdtBalance(
	ctx comcontext.Context, order *models.Order,
) (err error) {
	usdtCoin := blockchainmod.GetCoinNativeF(constants.CurrencyTetherUSD)
	usdtAccount, err := usdtCoin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return
	}
	ethCoin := blockchainmod.GetCoinNativeF(constants.CurrencyEthereum)
	ethAccount, err := ethCoin.NewAccountSystem(ctx, order.UID)
	if err != nil {
		return
	}

	feeInfo, err := usdtAccount.GetFeeInfoToAddress("")
	if err != nil {
		return
	}
	ethBalance, err := ethAccount.GetBalance()
	if err != nil {
		return
	}

	ethMaxFee := feeInfo.GetBaseLimitMaxValue()
	if ethBalance.LessThan(ethMaxFee) {
		return constants.ErrorBlockchainBalanceNotEnoughForFee.WithData(meta.O{
			"threshold": ethMaxFee,
			"currency":  feeInfo.BaseCurrency,
		})
	}

	return nil
}
