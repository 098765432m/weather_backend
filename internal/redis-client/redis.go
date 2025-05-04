package redisClient

import (
	"context"

	"github.com/redis/go-redis/v9"
)
func InitRedisClient(address string, password string, db int) (*redis.Client ,error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: password,
		DB: db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}