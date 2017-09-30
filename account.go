// Package redisbank provides an API for a banking system on top of Redis.
package redisbank

import (
	"log"
	"strconv"
)

func hasAccount(username string, accountName string) bool {
	return r.SIsMember("accounts:"+username, accountName).Val()
}

// GetBalance returns the float64 amount of money
// the username has on his accountName.
func GetBalance(username string, accountName string) float64 {
	if !hasAccount(username, accountName) {
		log.Println(username + " does not have a \"" + accountName + "\" account")
		return 0
	}
	amount, _ := strconv.ParseFloat(r.Get("account:"+username+":"+accountName).Val(), 64)

	return amount
}

// GetUserAccounts returns the total list of username's accounts
// as a slice of strings.
func GetUserAccounts(username string) []string {
	return r.SMembers("accounts:" + username).Val()
}

// NewAccount creates a new account given the username and a unique accountName.
func NewAccount(username string, accountName string) {
	if userExists(username) {
		if hasAccount(username, accountName) {
			log.Println(accountName, "for user", username+":", "account already existent")
		}
		pipe := r.TxPipeline()
		pipe.Set("account:"+username+":"+accountName, "0", 0)
		pipe.SAdd("accounts:"+username, accountName)
		pipe.Exec()
	} else {
		log.Println(accountName, "for user", username+":", "user non-existent")
	}
}

// DeleteAccount deletes a username's account corresponding to the name accountName.
// As per the REST APIs directives, it's idempotent.
func DeleteAccount(username string, accountName string) {
	pipe := r.TxPipeline()
	pipe.Del("account:" + username + ":" + accountName)
	pipe.Del("transactions:" + username + ":" + accountName)
	pipe.SRem("accounts:"+username, accountName)
	pipe.Exec()
}
