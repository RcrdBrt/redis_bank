// Package redisbank provides an API for a banking system on top of Redis.
package redisbank

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	// PRECISION is number of digits after "." in a float64
	PRECISION int = 4
)

var r = redis.NewClient(&redis.Options{
	Addr:         "127.0.0.1:6379",
	Password:     "",
	DB:           2,
	PoolSize:     4,
	PoolTimeout:  30 * time.Second,
	DialTimeout:  10 * time.Second,
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
})
var m = sync.Mutex{}
