package blockchainmod

import (
	"fmt"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comcache"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/dbquery"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/lockmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gorm.io/gorm"
)

func UpdateCurrencyLatestBlockHeight(coin Coin, blockHeight uint64) error {
	return database.GetDbF(database.AliasWalletMaster).
		Model(&models.BlockchainNetworkInfo{}).
		Where(&models.BlockchainNetworkInfo{
			Network: coin.GetNetwork(),
		}).
		Where(dbquery.Lt(models.BlockchainNetworkInfoColLatestBlockHeight, blockHeight)).
		Updates(&models.BlockchainNetworkInfo{
			LatestBlockHeight: blockHeight,
			UpdateTime:        time.Now().Unix(),
		}).
		Error
}

func CreateUserAddress(ctx comcontext.Context, uid meta.UID, coin Coin) (
	userAddress models.UserAddress, err error,
) {
	network := coin.GetNetwork()
	lock, err := lockmod.LockSimple("blockchain:user:create_address:%v-%v", uid, network)
	if err != nil {
		return
	}
	defer lock.Unlock()

	keyHolder, err := coin.NewKey()
	if err != nil {
		return
	}

	err = database.Atomic(ctx, database.AliasWalletMaster, func(dbTxn *gorm.DB) (err error) {
		err = dbTxn.
			First(&userAddress, &models.UserAddress{UID: uid, Network: network}).
			Error
		if database.IsDbError(err) {
			return utils.WrapError(err)
		}
		if userAddress.ID != 0 {
			return nil
		}

		userAddress.UID = uid
		userAddress.Network = network
		userAddress.Address = keyHolder.GetAddress()
		userAddress.Key = models.NewUserAddressKeyEncryptedField(keyHolder.GetPrivateKey())
		userAddress.CreateTime = time.Now().Unix()

		if err = dbTxn.Create(&userAddress).Error; err != nil {
			return utils.WrapError(err)
		}

		return nil
	})
	return
}

func GetOrCreateUserAddress(
	ctx comcontext.Context, uid meta.UID, coin Coin,
) (userAddress models.UserAddress, err error) {
	err = database.
		GetDbF(database.AliasWalletSlave).
		First(&userAddress, &models.UserAddress{UID: uid, Network: coin.GetNetwork()}).
		Error
	if database.IsDbError(err) {
		err = utils.WrapError(err)
		return
	}
	if userAddress.ID > 0 {
		return
	}

	userAddress, err = CreateUserAddress(ctx, uid, coin)
	return
}

func GetUserAddress(ctx comcontext.Context, uid meta.UID, coin Coin) (models.UserAddress, error) {
	return GetOrCreateUserAddress(ctx, uid, coin)
}

func GetUserAddressFast(
	ctx comcontext.Context, uid meta.UID, coin Coin,
) (userAddress models.UserAddress, err error) {
	cacheKey := fmt.Sprintf("blockchain:address:%v:%v", uid, coin.GetIndexCode())

	err = comcache.GetOrCreate(
		comcache.GetMemoryCache(),
		cacheKey,
		5*time.Minute,
		&userAddress,
		func() (interface{}, error) {
			return GetUserAddress(ctx, uid, coin)
		},
	)
	return
}

func GetUserAddressFastF(ctx comcontext.Context, uid meta.UID, coin Coin) models.UserAddress {
	address, err := GetUserAddressFast(ctx, uid, coin)
	comutils.PanicOnError(err)

	return address
}
