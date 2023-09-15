package fastcachekeyvaluestore

import (
	"encoding/json"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"gitlab.com/gear5th/gear5th-app/internal/application"
)

const maxCacheSizeBytes = 1 * 1024 * 1024 * 1024 // 1GB

var store *fastcache.Cache

func init() {
	store = fastcache.New(maxCacheSizeBytes)
}

type cachePayload struct {
	Value  string    `json:"value"`
	Expiry time.Time `json:"expiry"`
}

type FastCacheKeyValueStore struct{}

func NewFastCacheKeyValueStore() FastCacheKeyValueStore {
	return FastCacheKeyValueStore{}
}

func (s FastCacheKeyValueStore) Get(key string) (string, error) {
	payload := make([]byte, 0)
	res := store.GetBig(payload, []byte(key))
	p := &cachePayload{}
	err := json.Unmarshal(res, p)
	if err != nil {
		return "", application.ErrEntityNotFound
	}
	if p.Expiry.Before(time.Now()) {
		// Cache has expired
		store.Del([]byte(key))
		return "", application.ErrEntityNotFound
	}
	return p.Value, nil
}

func (s FastCacheKeyValueStore) Save(key string, value string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = 8760 * time.Hour // 1 year 
	}
	p := &cachePayload{
		Value:  value,
		Expiry: time.Now().Add(ttl),
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	store.SetBig([]byte(key), payload)

	return nil
}
