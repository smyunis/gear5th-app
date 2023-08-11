package application

import "time"

type KeyValueStore interface {
	Get(key string) (string, error)
	Save(key string, value string, ttl time.Duration) error
}
