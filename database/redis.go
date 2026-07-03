package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	address := os.Getenv("REDIS_ADDR")
	if address == "" {
		address = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr: address,
		DB:   0,
	})
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Redis connected succesfully")
}
