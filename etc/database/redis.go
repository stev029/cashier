package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	opt, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	RedisClient = redis.NewClient(opt)

	if err := RedisClient.Ping(context.TODO()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}

	log.Println("Connected to Redis")
}
