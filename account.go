package redis_bank

import (
	"errors"
	"strconv"
)

func hasAccount(username string, accountName string) bool {
	return r.SIsMember("accounts:"+username, accountName).Val()
}

func getBalance(username string, accountName string) (float64, error) {
	if !hasAccount(username, accountName) {
		err := username + "does not have a \"" + accountName + "\" account"
		return 0, errors.New(err)
	}
	amount, err := strconv.ParseFloat(r.Get("account:"+username+":"+accountName).Val(), 64)

	return amount, err
}

func getUserAccounts(username string) []string {
	return r.SMembers("accounts:" + username).Val()
}

func NewAccount(username string, accountName string) error {
	if hasAccount(username, accountName) {
		err := "Account already existent"
		return errors.New(err)
	}
	m.Lock()
	pipe := r.TxPipeline()
	pipe.Set("account:"+username+":"+accountName, "0", 0)
	pipe.SAdd("accounts:"+username, accountName)
	pipe.Exec()
	m.Unlock()

	return nil
}

func DeleteAccount(username string, accountName string) error {
	if hasAccount(username, accountName) {
		err := "Account not existent"
		return errors.New(err)
	}
	m.Lock()
	pipe := r.TxPipeline()
	pipe.Del("account:" + username + ":" + accountName)
	pipe.SRem("accounts:"+username, accountName)
	pipe.Exec()
	m.Unlock()

	return nil
}
