package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(redisURL string) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)

	return &Cache{
		client: client,
		ttl:    24 * time.Hour,
	}, nil
}

func (c *Cache) Get(ctx context.Context, nameOrID string) (*Pokemon, error) {
	key := fmt.Sprintf("pokemon:%s", nameOrID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var pokemon Pokemon
	if err := json.Unmarshal([]byte(val), &pokemon); err != nil {
		return nil, fmt.Errorf("unmarshal pokemon: %w", err)
	}

	return &pokemon, nil
}

func (c *Cache) Set(ctx context.Context, nameOrID string, pokemon *Pokemon) error {
	key := fmt.Sprintf("pokemon:%s", nameOrID)

	data, err := json.Marshal(pokemon)
	if err != nil {
		return fmt.Errorf("marshal pokemon: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}
