package redis

import (
	"context"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := r.NewClient(&r.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

func (rc *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return rc.Client.Publish(ctx, channel, message).Err()
}

func (rc *RedisClient) Subscribe(ctx context.Context, channel string) *r.PubSub {
	return rc.Client.Subscribe(ctx, channel)
}
