package redis_bank

import (
	"github.com/go-redis/redis"
	"sync"
)

const (
	TIMEOUT   int = 2
	PRECISION int = 4 // number of digits after "." in the amount
)

var r *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       2,
})
var m sync.Mutex = sync.Mutex{}
