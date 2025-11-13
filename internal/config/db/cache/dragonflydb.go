// Package cache : implement in-memory cache using dragonflyDB and redis-client
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
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

		utils.LogInfo(utils.DragonflyIsConnected.Type, utils.DragonflyIsConnected.Msg)

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
		err := errors.New("Cache not initialized. Call InitCache() first")
		utils.LogErrorWithLevel("fatal", utils.DragonflyFailedToInit.Type, utils.DragonflyFailedToInit.Code, utils.DragonflyFailedToInit.Msg, err)
	}
	return instance
}

// getOptions : getting the db options
func getDragonFlyOptions() *redis.Options {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.LogErrorWithLevel("fatal", utils.DragonflyFailedToLoadOptions.Type, utils.DragonflyFailedToLoadOptions.Code, utils.DragonflyFailedToLoadOptions.Msg, err)
	}

	address := fmt.Sprintf("%s:%s", cfg.DragonflyDB.Host, cfg.DragonflyDB.Port)
	return &redis.Options{
		Addr:     address,
		Password: cfg.DragonflyDB.Password,
		DB:       cfg.DragonflyDB.DB,

		// disable MAINT_NOTIFICATIONS which is a warning system designed to send a "heads-up" message
		// to your application before the database has to restart for a maintenance event.
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
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
