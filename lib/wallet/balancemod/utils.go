package balancemod

import (
	"fmt"
	"math/rand"
	"time"
)

func genTxnCode(prefix string) string {
	todayStr := time.Now().Format("060102")
	randInt := rand.Uint64() % 1000000
	return fmt.Sprintf("%s-%s-%06d", prefix, todayStr, randInt)
}

func GenInvestmentWithdrawCode() string {
	return genTxnCode("C")
}

func GenProfitWithdrawCode() string {
	return genTxnCode("T")
}

func GenProfitReinvestCode() string {
	txnCode := genTxnCode("R")
	randInt := rand.Uint64() % 100000000
	return fmt.Sprintf("%s-%08d", txnCode, randInt)
}
