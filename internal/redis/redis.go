package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Client *redis.Client
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Client{
		Client: rdb,
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *Client) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := c.Client.Ping(ctx).Err()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = err.Error()
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	return stats
}
