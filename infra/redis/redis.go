package redis

import (
	"context"
	"net/http"
	"ratelimit/util"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisClientStruct struct {
	redis *redis.Client
	mutex sync.Mutex
}

var RedisClient *RedisClientStruct

func InitializeRedis() {

	if RedisClient == nil {
		client := redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		})

		RedisClient = &RedisClientStruct{
			redis: client,
		}
	}
}

func Incr(key string) int64 {
	RedisClient.mutex.Lock()
	defer RedisClient.mutex.Unlock()

	val, err := RedisClient.redis.Incr(key).Result()
	if err != nil {
		return 0
	}

	return val
}

func SetExpiry(ctx *context.Context, key string, expiry time.Duration) bool {

	RedisClient.mutex.Lock()
	defer RedisClient.mutex.Unlock()

	val, err := RedisClient.redis.Expire(key, expiry).Result()
	if err != nil {
		return false
	}
	return val
}

func GetGet(ctx *context.Context, ket string) map[string]string {

	RedisClient.mutex.Lock()
	defer RedisClient.mutex.Unlock()

	result, err := RedisClient.redis.HGetAll(ket).Result()
	util.ErrorIf(err != nil, "REDIS_FETCH_ERROR", http.StatusInternalServerError, "redis key not found")

	return result
}
