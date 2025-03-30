package data

import "sync"

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
