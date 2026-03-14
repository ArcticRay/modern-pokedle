package pokemon

import (
	"context"
	"fmt"
)

type Service struct {
	client *Client
	cache  *Cache
}

func NewService(client *Client, cache *Cache) *Service {
	return &Service{
		client: client,
		cache:  cache,
	}
}

func (s *Service) GetPokemon(ctx context.Context, nameOrID string) (*Pokemon, error) {
	if cached, err := s.cache.Get(ctx, nameOrID); err == nil && cached != nil {
		return cached, nil
	}

	pokemon, err := s.client.GetPokemon(ctx, nameOrID)
	if err != nil {
		return nil, fmt.Errorf("get pokemon: %w", err)
	}

	if err := s.cache.Set(ctx, nameOrID, pokemon); err != nil {
		return nil, fmt.Errorf("cache pokemon: %w", err)
	}

	return pokemon, nil
}
