package redis_bank

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	TIMEOUT   int = 2
	PRECISION int = 4 // number of digits after "." in the amount
)

var r *redis.Client = redis.NewClient(&redis.Options{
	Addr:         "127.0.0.1:6379",
	Password:     "",
	DB:           2,
	PoolSize:     4,
	PoolTimeout:  30 * time.Second,
	DialTimeout:  10 * time.Second,
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
})
var m sync.Mutex = sync.Mutex{}
