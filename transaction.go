package redis_bank

import (
	"errors"
	"strconv"
)

func UpdateBalance(username string, accountName string, amount float64) error {
	if !hasAccount(username, accountName) {
		err := "Transaction in non-existent account"
		return errors.New(err)
	}
	m.Lock()
	pipe := r.TxPipeline()
	prevAmount, err := getBalance(username, accountName)
	if err != nil {
		return err
	}
	newAmount := prevAmount + amount
	pipe.RPush("transactions:"+username+":"+accountName, newAmount)
	pipe.Set("account:"+username+":"+accountName,
		strconv.FormatFloat(newAmount, 'f', PRECISION, 64), 0)
	pipe.Exec()
	m.Unlock()

	return nil
}

func updateBalancePercentage(username string, accountName string, perc float64) error {
	prevAmount, err := getBalance(username, accountName)
	if err != nil {
		return err
	}
	return updateBalance(username, accountName, prevAmount/100*perc)
}
