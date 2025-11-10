// Package cache : implement in-memory cache using dragonflyDB and redis-client
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	once     sync.Once
	instance *Cache
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
	prefix string
}

// InitCache : init dragonflyDB
func InitCache(prefix string) (*Cache, error) {
	var initErr error
	once.Do(func() {
		// loading dragonfly options
		drgonOpts := getDragonFlyOptions()

		// new client instance
		client := redis.NewClient(drgonOpts)

		// Ping to check connection
		ctx := context.Background()

		if _, err := client.Ping(ctx).Result(); err != nil {
			initErr = fmt.Errorf("failed to connect to DragonflyDB: %v", err)
			return
		}

		log.Println("Connected to DragonflyDB successfully...")

		// return client
		instance = &Cache{
			client: client,
			ctx:    ctx,
			prefix: prefix,
		}
	})
	return instance, initErr
}

// GetInstance returns the singleton cache instance
func GetInstance() *Cache {
	if instance == nil {
		log.Fatal("Cache not initialized. Call InitCache() first.")
	}
	return instance
}

// getOptions : getting the db options
func getDragonFlyOptions() *redis.Options {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("unable to load config : %v", err)
	}

	address := fmt.Sprintf("%s:%s", cfg.DragonflyDB.Host, cfg.DragonflyDB.Port)
	return &redis.Options{
		Addr:     address,
		Password: cfg.DragonflyDB.Password,
		DB:       cfg.DragonflyDB.DB,
	}
}

// HealthCheck checks if cache is healthy
func (c *Cache) HealthCheck() error {
	_, err := c.client.Ping(c.ctx).Result()
	return err
}

// --- AI GEN --- //

// buildKey adds prefix to the key
func (c *Cache) buildKey(key string) string {
	if c.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Set stores a value with expiration
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return c.client.Set(c.ctx, c.buildKey(key), data, expiration).Err()
}

// Get retrieves a value
func (c *Cache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, c.buildKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// --- AI GEN --- //
