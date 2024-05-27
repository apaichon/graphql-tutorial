package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	redisOnce     sync.Once
	redisInstance *RedisClient
)

// RedisClient represents a simple Redis client.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient() (*RedisClient, error) {

	addr := viper.GetString("CACHE_CON_STR")
	password := viper.GetString("CACHE_PASSWORD")
	db := viper.GetInt("CACHE_INDEX")
	fmt.Printf("Address:%s Index:%v", addr, db)

	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

// GetInstance returns the singleton instance of the Redis client.
func GetRedisInstance() (*RedisClient, error) {
	redisOnce.Do(func() {
		var err error
		redisInstance, err = NewRedisClient()
		if err != nil {
			log.Fatalf("Error creating Redis client: %v", err)
			// fmt.Printf("Error creating Redis client: %v", err)
		}
	})
	return redisInstance, nil
}

// Close closes the Redis connection.
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

// Get retrieves the value associated with the given key from Redis.
func (rc *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' not found", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get value for key '%s': %v", key, err)
	}
	return val, nil
}

// Set sets the value associated with the given key in Redis.
func (rc *RedisClient) Set(key, value string) error {
	ctx := context.Background()
	cacheAge := viper.GetInt("CACHE_AGE")
	err := rc.client.Set(ctx, key, value, time.Duration(cacheAge)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set value for key '%s': %v", key, err)
	}
	return nil
}

// Remove removes the specified key from Redis.
func (rc *RedisClient) Remove(key string) error {
	ctx := context.Background()
	deleted, err := rc.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to remove key '%s' from Redis: %v", key, err)
	}
	if deleted == 0 {
		return fmt.Errorf("key '%s' does not exist in Redis", key)
	}
	return nil
}

// Remove removes the specified key from Redis.
func (rc *RedisClient) Removes(key string) {
	ctx := context.Background()
	// Use Lua script to delete keys by pattern
	script := `
	 local keys = redis.call('KEYS', ARGV[1])
	 for i=1,#keys do
		 redis.call('DEL', keys[i])
	 end
	 return keys
 `
	// Execute Lua script
	result, err := rc.client.Eval(ctx, script, []string{}, key+"*").Result()
	if err != nil {
		panic(err)
	}

	// Print deleted keys
	deletedKeys, _ := result.([]interface{})
	log.Printf("deletedKeys%v", deletedKeys...)
	for _, key := range deletedKeys {
		// deleted, err :=
		rc.client.Del(ctx, key.(string)) //.Result()
	}
}
