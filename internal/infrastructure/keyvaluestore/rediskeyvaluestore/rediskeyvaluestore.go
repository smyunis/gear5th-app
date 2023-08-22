package rediskeyvaluestore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/gear5th/gear5th-app/internal/application"
)

type RedisKeyValueStore struct {
	client *redis.Client
}

func NewRedisKeyValueStore(bootstrapper RedisBootstrapper) RedisKeyValueStore {
	return RedisKeyValueStore{
		bootstrapper.Client,
	}
}

func (s RedisKeyValueStore) Get(key string) (string, error) {
	r := s.client.Get(context.Background(), key)
	res, err := r.Result()
	if err != nil {
		return "", application.ErrEntityNotFound
	}
	return res, nil
}

func (s RedisKeyValueStore) Save(key string, value string, ttl time.Duration) error {
	r := s.client.Set(context.Background(), key, value, ttl)
	return r.Err()
}
