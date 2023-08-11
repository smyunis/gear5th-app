package rediskeyvaluestore

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
)

var client *redis.Client

type RedisBootstrapper struct {
	Client *redis.Client
}

func NewRedisBootstrapper(config infrastructure.ConfigurationProvider) RedisBootstrapper {

	if client == nil {
		redisURL := config.Get("REDIS_URL", "redis://localhost:6379/0")
		url, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(fmt.Errorf("could not connect to redis database: %w", err))
		}
		client = redis.NewClient(url)
	}
	return RedisBootstrapper{
		client,
	}
}
