package cache

import (
	"github.com/redis/go-redis/v9"
	"context"
)

var Ctx = context.Background()

type Redis struct{
	Client *redis.Client
}

func NewRedis()*Redis{
	client:=redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:0,
	})
	return &Redis{Client: client}
}