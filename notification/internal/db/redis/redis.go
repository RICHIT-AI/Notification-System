package redisDB

import "github.com/go-redis/redis/v9"

func InitRedis() *redis.Client {
	// TODO: Get variables from environment
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
