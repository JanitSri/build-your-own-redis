package data

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/JanitSri/codecrafters-build-your-own-redis/customerror"
)

type DataStore interface {
	Get(key any) (any, bool)
	Set(key, value any)
	Keys() []any
	GetConfig(string) string
}

type RedisStore struct {
	cmap   sync.Map
	config RedisConfig
}

func NewRedisStore(rc RedisConfig) *RedisStore {
	return &RedisStore{
		config: rc,
	}
}

func (rs *RedisStore) Get(key any) (any, bool) {
	return rs.cmap.Load(key)
}

func (rs *RedisStore) Set(key, value any) {
	rs.cmap.Store(key, value)
}

func (rs *RedisStore) Keys() []any {
	var keys []any
	rs.cmap.Range(func(k, _ any) bool {
		keys = append(keys, k)
		return true
	})

	return keys
}

type RedisValue struct {
	value  any
	expiry time.Time
}

func NewRedisValue(value any, expiry time.Time) *RedisValue {
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

func (rv RedisValue) Value() any {
	return rv.value
}

type RedisConfig struct {
	dir        string
	dbFileName string
}

func NewRedisConfig(dir, dbFileName string) *RedisConfig {
	return &RedisConfig{
		dir,
		dbFileName,
	}
}

func (rs *RedisStore) GetConfig(name string) string {
	var c string
	switch strings.ToUpper(name) {
	case "DIR":
		c = rs.config.dir
	case "DBFILENAME":
		c = rs.config.dbFileName
	default:
		log.Fatal(customerror.InvalidServerConfigError{Name: name})
	}

	return c
}
