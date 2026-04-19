package cache

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

type Redis struct {
	Client *redis.Client
}

func NewRedis() *Redis {

	redisURL := os.Getenv("REDIS_URL")

	var client *redis.Client

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatal("Invalid REDIS_URL:", err)
		}
		client = redis.NewClient(opt)

	} else {
		client = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		})
	}
	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {

		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected successfully")

	return &Redis{Client: client}
}
