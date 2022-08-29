package orderchannel

import (
	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/blockchainmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/wallet/ordermod"
)

type BlockchainNetworkOrderExtraData struct {
	FromAddress string              `json:"from_address"`
	Fee         meta.CurrencyAmount `json:"fee"`
	ExplorerURL string              `json:"explorer_url"`
}

func DumpBlockchainTxnToOrder(
	uid meta.UID,
	txn blockchainmod.Transaction, order *models.Order,
	dumpDetails bool,
) (err error) {
	txnAmount, err := txn.GetAmount()
	if err != nil {
		return err
	}

	order.UID = uid
	order.Direction = txn.GetDirection()
	order.Currency = txn.GetCurrency()
	order.Status = blockchainmod.GetBlockchainOrderStatus(txn.GetLocalStatus())

	order.AmountSubTotal = txnAmount
	order.AmountFee = decimal.Zero
	order.AmountDiscount = decimal.Zero
	order.AmountTotal = txnAmount.Add(order.AmountFee)

	order.CreateTime = txn.GetTimeUnix()
	order.UpdateTime = order.CreateTime

	order.SrcChannelType = constants.ChannelTypeSrcBlockchainNetwork
	order.SrcChannelAmount = txnAmount
	order.DstChannelType = constants.ChannelTypeDstBlockchainNetwork
	order.DstChannelAmount = txnAmount

	if order.Direction == constants.DirectionTypeReceive {
		order.SrcChannelRef = txn.GetHash()
	} else {
		order.DstChannelRef = txn.GetHash()

		dstChannelMeta := DstBlockchainNetworkMeta{}
		if dstChannelMeta.InputDataHex, err = txn.GetInputDataHex(); err != nil {
			return
		}
		if dstChannelMeta.ToAddress, _ = txn.GetToAddress(); dstChannelMeta.ToAddress == "" {
			dstChannelMeta.ToAddress = "<unknown>"
		}

		err = ordermod.SetOrderChannelMetaData(
			order,
			constants.ChannelTypeDstBlockchainNetwork, &dstChannelMeta)
		if err != nil {
			return
		}
	}

	return DumpBlockchainInfoToOrder(order, txn)
}

func DumpBlockchainInfoToOrder(order *models.Order, txn blockchainmod.Transaction) error {
	if _, ok := order.ExtraData[constants.OrderExtraDataBlockchainInfo]; ok {
		return nil
	}

	var txnHash string
	if order.Direction == constants.DirectionTypeSend {
		txnHash = order.DstChannelRef
	} else {
		txnHash = order.SrcChannelRef
	}

	var (
		coin      = blockchainmod.GetCoinNativeF(order.Currency)
		extraData = BlockchainNetworkOrderExtraData{
			ExplorerURL: coin.GenTxnExplorerURL(txnHash),
		}
	)
	if txn != nil {
		extraData.Fee = txn.GetFee()

		if fromAddress, err := txn.GetFromAddress(); err == nil {
			extraData.FromAddress = fromAddress
		} else {
			extraData.FromAddress = "<unknown>"
		}
	}
	order.ExtraData[constants.OrderExtraDataBlockchainInfo] = extraData

	return nil
}
