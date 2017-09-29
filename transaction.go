package redis_bank

import (
	"log"
	"strconv"
)

func GetLastTransaction(username string, accountName string) float64 {
	if userExists(username) && hasAccount(username, accountName) {
		str_value := r.LRange("transactions:"+username+":"+accountName, 0, 0).Val()[0]
		return strconv.ParseFloat(str_value, 64)
	}

	return 0
}

func UpdateBalance(username string, accountName string, amount float64) {
	if userExists(username) {
		if !hasAccount(username, accountName) {
			log.Println(username, "in", accountName+":", "transaction in non-existent account")
		} else {
			pipe := r.TxPipeline()
			prevAmount := GetBalance(username, accountName)
			newAmount := prevAmount + amount
			pipe.LPush("transactions:"+username+":"+accountName, amount)
			pipe.Set("account:"+username+":"+accountName,
				strconv.FormatFloat(newAmount, 'f', PRECISION, 64), 0)
			pipe.Exec()
		}
	} else {
		log.Println(username, accountName+":", "transaction for non-existent user")
	}
}

func RevertLastTransaction(usernames string, account string) {
	UpdateBalance(username, accountName, GetLastTransaction(username, account))
}

func updateBalancePercentage(username string, accountName string, perc float64) {
	prevAmount := GetBalance(username, accountName)
	UpdateBalance(username, accountName, prevAmount/100*perc)
}
