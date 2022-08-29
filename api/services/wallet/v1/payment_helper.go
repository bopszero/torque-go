package v1

import (
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func genUserExecuteOrderTokenCacheKey(uid meta.UID, token string) string {
	return fmt.Sprintf("order:token:%v:%v", uid, token)
}

func dumpOrder(order *models.Order, responseOrder *Order) error {
	if err := copier.Copy(responseOrder, order); err != nil {
		return err
	}

	extraData := meta.O{}
	for key, value := range order.ExtraData {
		if strings.HasPrefix(key, "_") {
			continue
		}

		extraData[key] = value
	}

	responseOrder.ExtraData = extraData

	return nil
}
