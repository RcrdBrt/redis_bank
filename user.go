package redis_bank

import (
	"errors"
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

	return r.HGet("user_ids", username).String()
}

func getUserHash(username string) string {
	return username + ":" + getUserId(username)
}

func newUser(username string, passwd string) error {
	if userExists(username) { // user already registered
		err := "Username already registered!"
		return errors.New(err)
	}
	// brand new user
	enc_passwd, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	m.Lock()               // mutex ON
	pipe := r.TxPipeline() // pipeline start
	pipe.Incr("user_id")
	pipe.Incr("user_tot")
	pipe.HMSet(getUserHash(username), map[string]interface{}{
		"password": string(enc_passwd),
		"enabled":  "1",
	})
	pipe.HSet("user_ids", username, r.Get("user_id"))
	if _, err := pipe.Exec(); err != nil { // pipeline exec
		return err
	}
	m.Unlock() // mutex OFF

	return nil
}

func deleteUser(username string) error {
	if !userExists(username) {
		err := "Username is not registered!"
		return errors.New(err)
	}
	m.Lock()               // mutex ON
	pipe := r.TxPipeline() // pipeline start
	pipe.Decr("user_tot")
	pipe.HDel(getUserHash(username), "password", "enabled")
	pipe.HDel("user_ids", username)
	if _, err := pipe.Exec(); err != nil { // pipeline exec
		return err
	}
	m.Unlock() // mutex OFF

	return nil
}

func authUser(username string, passwd string) bool {
	if !userExists(username) {
		log.Println("Username is not registered!")
		return false
	}
	stored_passwd := []byte(r.HGet(getUserHash(username), "password").String())
	err := bcrypt.CompareHashAndPassword(stored_passwd, []byte(passwd))
	if err != nil {
		log.Println()
		return false
	}

	return true
}
