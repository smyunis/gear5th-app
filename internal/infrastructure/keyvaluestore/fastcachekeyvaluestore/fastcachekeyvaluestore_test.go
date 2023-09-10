package fastcachekeyvaluestore_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/keyvaluestore/fastcachekeyvaluestore"
)

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

var cache fastcachekeyvaluestore.FastCacheKeyValueStore

func setup() {
	cache = fastcachekeyvaluestore.NewFastCacheKeyValueStore()
}

func TestSaveACache(t *testing.T) {
	key := "mykey"
	val := "myval"

	cache.Save(key, val, 2*time.Hour)

}

func TestCanGetSavedCache(t *testing.T) {
	key := "mykey"
	val := "myval"

	cache.Save(key, val, 2*time.Hour)

	v, err := cache.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if val != v {
		t.FailNow()
	}
}

func TestCanNotGetExpiredCache(t *testing.T) {
	key := "mykey"
	val := "myval"

	cache.Save(key, val, 2*time.Second)
	time.Sleep(3 * time.Second)

	v, err := cache.Get(key)
	if err == nil || v != "" {
		t.FailNow()
	}

	if !errors.Is(err, application.ErrEntityNotFound) {
		t.FailNow()
	}
}
