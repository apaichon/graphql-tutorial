package cache

import (
	"log"
)

// CacheBackend represents the backend for the cache.
type CacheBackend int

const (
	RedisBackend CacheBackend = iota
	SQLiteBackend
)

func IntToCacheBackend(i int) CacheBackend {
	switch i {
	case 0:
		return RedisBackend
	case 1:
		return SQLiteBackend
	default:
		return RedisBackend // or return an error, depending on your use case
	}
}

// Cache represents a cache with support for different backends.
type Cache struct {
	backend CacheBackend
	db      CacheDB
}

// CacheDB represents the interface for interacting with the cache database.
type CacheDB interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Remove(key string) error
	Removes(key string)
	Close() error
}

// NewCache creates a new cache with the specified backend.
func NewCache(backend CacheBackend) *Cache {
	var db CacheDB
	switch backend {
	case RedisBackend:
		db,err := GetRedisInstance()
		if (err !=nil) {
			panic(err)
		}
		return &Cache{backend: backend, db: db}
	case SQLiteBackend:
		// db = &SQLiteInMemClient{}
		db,err := GetSqliteInMemInstance()
		if (err !=nil) {
			panic(err)
		}
		return &Cache{backend: backend, db: db}
	default:
		log.Fatalf("Unsupported cache backend: %v", backend)
	}
	return &Cache{backend: backend, db: db}
}

// Get retrieves the value associated with the given key from the cache.
func (c *Cache) Get(key string) (string, error) {
	return c.db.Get(key)
}

// Set sets the value associated with the given key in the cache.
func (c *Cache) Set(key, value string) error {
	return c.db.Set(key, value)
}

// Remove removes the specified key from the cache.
func (c *Cache) Remove(key string) error {
	return c.db.Remove(key)
}

// Remove removes the specified key from the cache.
func (c *Cache) Removes(key string) {
	c.db.Removes(key)
}