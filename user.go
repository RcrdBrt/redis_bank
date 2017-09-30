// Package redisbank provides an API for a banking system on top of Redis.
package redisbank

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func userExists(username string) bool {
	return r.HExists("user_ids", username).Val()
}

// GetUserID returns the username unique (incremental) ID as a string.
func GetUserID(username string) string {
	if !userExists(username) {
		return ""
	}

	return r.HGet("user_ids", username).Val()
}

// UserIsEnabled checks whether user is enabled or not.
// The logic in case he isn't is missing.
func UserIsEnabled(username string) string {
	return r.HGet(username+":"+GetUserID(username), "enabled").Val()
}

func getUserHash(username string) string {
	return username + ":" + GetUserID(username)
}

// NewUnsecureUser creates a new user given a username with
// the hardcoded "default" password
func NewUnsecureUser(username string) {
	NewUser(username, "default")
}

// NewUser creates a new user given a username and a password.
// Username must not be in the system already or it fails.
// There are no general restrictions for the password.
func NewUser(username string, passwd string) {
	if userExists(username) { // user already registered
		log.Println(username+":", "username already registered!")
	} else {
		// brand new user
		encPasswd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		r.Incr("user_id")
		r.HSet("user_ids", username, r.Get("user_id").Val())
		r.Incr("user_tot")
		r.HMSet(getUserHash(username), map[string]interface{}{
			"password": string(encPasswd),
			"enabled":  "1",
		})
	}
}

// DeleteUser deletes the user username from the system and
// removes all his accounts and transactions logs.
func DeleteUser(username string) {
	if !userExists(username) {
		log.Println(username+":", "username is not registered!")
	}
	pipe := r.TxPipeline() // pipeline start
	accounts := GetUserAccounts(username)
	userID := GetUserID(username)
	for _, val := range accounts {
		pipe.Del("transactions:" + username + ":" + val)
		DeleteAccount(username, val)
	}
	pipe.Del(username + ":" + userID)
	pipe.Decr("user_tot")
	pipe.Del(getUserHash(username))
	pipe.Del("accounts:" + username)

	pipe.HDel("user_ids", username)
	if _, err := pipe.Exec(); err != nil { // pipeline exec
		log.Println("DeleteUser():", err.Error())
	}
}

// AuthUser authenticates username and passwd against the database
// and returns "true" or "false" in case it's successful or not respectively.
// Usage of this function is still missing in this system.
func AuthUser(username string, passwd string) bool {
	if !userExists(username) {
		log.Println("Username is not registered!")
		return false
	}
	storedPasswd := []byte(r.HGet(getUserHash(username), "password").Val())
	err := bcrypt.CompareHashAndPassword(storedPasswd, []byte(passwd))
	if err != nil {
		log.Println()
		return false
	}

	return true
}
