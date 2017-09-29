package redis_bank

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func userExists(username string) bool {
	return r.HExists("user_ids", username).Val()
}

func getUserId(username string) string {
	if !userExists(username) {
		return ""
	}

	return r.HGet("user_ids", username).Val()
}

func getUserHash(username string) string {
	return username + ":" + getUserId(username)
}

func NewUnsecureUser(username string) {
	NewUser(username, "default")
}

func NewUser(username string, passwd string) {
	if userExists(username) { // user already registered
		log.Println(username+":", "username already registered!")
	}
	// brand new user
	enc_passwd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	m.Lock() // mutex ON
	r.Incr("user_id")
	r.HSet("user_ids", username, r.Get("user_id").Val())
	r.Incr("user_tot")
	r.HMSet(getUserHash(username), map[string]interface{}{
		"password": string(enc_passwd),
		"enabled":  "1",
	})
	m.Unlock() // mutex OFF
}

func DeleteUser(username string) {
	if !userExists(username) {
		log.Println(username+":", "username is not registered!")
	}
	pipe := r.TxPipeline() // pipeline start
	accounts := GetUserAccounts(username)
	user_id := getUserId(username)
	for _, val := range accounts {
		pipe.Del("transactions:" + username + ":" + val)
		DeleteAccount(username, val)
	}
	pipe.Del(username + ":" + user_id)
	pipe.Decr("user_tot")
	pipe.Del(getUserHash(username))
	pipe.Del("accounts:" + username)

	pipe.HDel("user_ids", username)
	if _, err := pipe.Exec(); err != nil { // pipeline exec
		log.Println("DeleteUser():", err.Error())
	}
}

func AuthUser(username string, passwd string) bool {
	if !userExists(username) {
		log.Println("Username is not registered!")
		return false
	}
	stored_passwd := []byte(r.HGet(getUserHash(username), "password").Val())
	err := bcrypt.CompareHashAndPassword(stored_passwd, []byte(passwd))
	if err != nil {
		log.Println()
		return false
	}

	return true
}
