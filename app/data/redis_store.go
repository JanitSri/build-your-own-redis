package data

import (
	"sync"
	"time"
)

type DataStore interface {
	Get(key any) (any, bool)
	Set(key, value any)
}

type RedisStore struct {
	cmap sync.Map
}

func NewRedisStore() *RedisStore {
	return &RedisStore{}
}

func (rs *RedisStore) Get(key any) (any, bool) {
	return rs.cmap.Load(key)
}

func (rs *RedisStore) Set(key, value any) {
	rs.cmap.Store(key, value)
}

type RedisValue struct {
	value  string
	expiry time.Time
}

func NewRedisValue(value string, expiry time.Time) *RedisValue {
	return &RedisValue{
		value,
		expiry,
	}
}

func (rv RedisValue) IsExpired() bool {
	return !rv.expiry.IsZero() && time.Now().After(rv.expiry)
}

func (rv *RedisValue) SetExpiry(t time.Time) {
	rv.expiry = t
}

func (rv RedisValue) Value() string {
	return rv.value
}
