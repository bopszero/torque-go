package blockchainmod

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

type EthereumAccount struct {
	blockchainGuestAccount
}

func (this *EthereumAccount) GetFeeInfoToAddress(toAddress string) (feeInfo FeeInfo, err error) {
	isContractAddress := false
	if toAddress != "" {
		isContractCacheKey := fmt.Sprintf("blockchain:account:eth:is_contract:%v", toAddress)

		err = comcache.GetOrCreate(
			comcache.GetRemoteCache(),
			isContractCacheKey,
			5*time.Minute,
			&isContractAddress,
			func() (interface{}, error) {
				getEthClient, err := NewEthereumEtherscanClientWithSystemKey()
				if err != nil {
					return nil, err
				}
				contractCode, err := getEthClient.GetCode(toAddress)
				if err != nil {
					return nil, err
				}

				return len(comutils.HexTrim(contractCode)) > 0, nil
			},
		)
		if err != nil {
			return
		}
	}

	cacheKey := fmt.Sprintf("blockchain:account:fee:%v-%v", this.GetCurrency(), isContractAddress)

	err = comcache.GetOrCreate(
		comcache.GetRemoteCache(),
		cacheKey,
		30*time.Second,
		&feeInfo,
		func() (interface{}, error) {
			feeInfo, err := this.coin.GetDefaultFeeInfo()
			if err != nil {
				return nil, err
			}
			if isContractAddress {
				feeInfo.LimitMaxQuantity = GetConfig().EthereumContractDefaultGasLimit
			}
			txnMaxFee := feeInfo.Price.
				Mul(decimal.NewFromInt(int64(feeInfo.LimitMaxQuantity)))
			feeInfo.SetLimitMaxValue(txnMaxFee, constants.CurrencySubEthereumWei)

			return feeInfo, nil
		},
	)
	if err != nil {
		return
	}

	return feeInfo, nil
}
