package testdoubles

import (
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

var cache map[string]string = make(map[string]string)

type KVStoreMock struct{}

func NewKVStoreMock() KVStoreMock {
	return KVStoreMock{}
}

func (KVStoreMock) Get(key string) (string, error) {
	val, ok := cache[key]
	if !ok {
		return "", shared.NewEntityNotFoundError(key, "in memory store")
	}
	return val, nil
}

func (KVStoreMock) Save(key string, value string, ttl time.Duration) error {
	cache[key] = value
	return nil
}
