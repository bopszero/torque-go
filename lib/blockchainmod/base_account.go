package blockchainmod

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type blockchainGuestAccount struct {
	coin    Coin
	client  Client
	address string
}

func (this *blockchainGuestAccount) GetCurrency() meta.Currency {
	return this.coin.GetCurrency()
}

func (this *blockchainGuestAccount) GetAddress() string {
	return this.address
}

func (this *blockchainGuestAccount) GetClient() Client {
	return this.client
}

func (this *blockchainGuestAccount) GetBalance() (decimal.Decimal, error) {
	return this.client.GetBalance(this.GetAddress())
}

func (this *blockchainGuestAccount) GetTxn(hash string) (Transaction, error) {
	txn, err := this.client.GetTxn(hash)
	if err != nil {
		return nil, err
	}

	txn.SetOwnerAddress(this.GetAddress())
	return txn, nil
}

func (this *blockchainGuestAccount) GetTxns(paging meta.Paging) ([]Transaction, error) {
	txns, err := this.client.GetTxns(this.GetAddress(), paging)
	if err != nil {
		return nil, err
	}

	for _, txn := range txns {
		txn.SetOwnerAddress(this.GetAddress())
	}
	return txns, nil
}

func (this *blockchainGuestAccount) PushTxn(data []byte) error {
	return this.client.PushTxnRaw(data)
}

func (this *blockchainGuestAccount) getAddressUtxOutputCountFast(address string) (
	utxOutputCount uint32, err error,
) {
	utxOutputCountCacheKey := fmt.Sprintf(
		"blockchain:account:%v:utxo_count:%v",
		this.GetCurrency(), address,
	)

	err = comcache.GetOrCreate(
		comcache.GetRemoteCache(),
		utxOutputCountCacheKey,
		15*time.Second,
		&utxOutputCount,
		func() (interface{}, error) {
			utxOutputs, err := this.client.GetUtxOutputs(
				this.GetAddress(), decimal.NewFromInt(MaxBalance))
			if err != nil {
				return nil, err
			}

			return len(utxOutputs), nil
		},
	)
	return
}

type blockchainOwnerAccount struct {
	GuestAccount
	coin      Coin
	keyHolder KeyHolder
}

func (this *blockchainOwnerAccount) GetPublicKey() string {
	return this.keyHolder.GetPublicKey()
}

func (this *blockchainOwnerAccount) GetPrivateKey() string {
	return this.keyHolder.GetPrivateKey()
}

func (this *blockchainOwnerAccount) GenerateSignedTxn(
	address string, amount decimal.Decimal, offerFeeInfo *FeeInfo,
) (_ SingleTxnSigner, err error) {
	txnSigner, err := this.coin.NewTxnSignerSingle(nil, offerFeeInfo)
	if err != nil {
		return
	}

	err = txnSigner.SetSrc(
		this.keyHolder.GetPrivateKey(),
		this.keyHolder.GetAddress())
	if err != nil {
		return
	}
	if err = txnSigner.SetDst(address, amount); err != nil {
		return
	}
	if err = txnSigner.Sign(false); err != nil {
		return
	}

	return txnSigner, nil
}

type blockchainSystemAccount struct {
	OwnerAccount
	addressInfo models.UserAddress
}

func (this *blockchainSystemAccount) GetUID() meta.UID {
	return this.addressInfo.UID
}
