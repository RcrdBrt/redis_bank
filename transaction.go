// Package redisbank provides an API for a banking system on top of Redis.
package redisbank

import (
	"log"
	"strconv"
)

// GetLastTransaction returns the float64 amount of the last transaction
// done by username in his accountName.
func GetLastTransaction(username string, accountName string) float64 {
	if userExists(username) && hasAccount(username, accountName) {
		list_name := "transactions:" + username + ":" + accountName
		if r.LLen(list_name).Val() > 0 {
			str_value := r.LRange(list_name, 0, 0).Val()[0]
			res, _ := strconv.ParseFloat(str_value, 64)

			return res
		} else {
			return 0
		}
	}

	return 0
}

// UpdateBalance creates a new transaction of the amount given in username's accountName.
func UpdateBalance(username string, accountName string, amount float64) {
	if userExists(username) {
		if !hasAccount(username, accountName) {
			log.Println(username, "in", accountName+":", "transaction in non-existent account")
		} else {
			prevAmount := GetBalance(username, accountName)
			newAmount := prevAmount + amount
			pipe := r.TxPipeline()
			pipe.LPush("transactions:"+username+":"+accountName, amount)
			pipe.Set("account:"+username+":"+accountName,
				strconv.FormatFloat(newAmount, 'f', PRECISION, 64), 0)
			pipe.Exec()
		}
	} else {
		log.Println(username, accountName+":", "transaction for non-existent user")
	}
}

// RevertLastTransaction undoes the last transaction given a username and an accountName.
// It does so by making a fully new transaction of the opposite amount of the reverting one.
func RevertLastTransaction(username string, accountName string) {
	UpdateBalance(username, accountName, -GetLastTransaction(username, accountName))
}

func updateBalancePercentage(username string, accountName string, perc float64) {
	prevAmount := GetBalance(username, accountName)
	UpdateBalance(username, accountName, prevAmount/100*perc)
}
