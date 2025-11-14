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

// GetFromCacheOrFetchDB implements cache-aside pattern using provided ctx
func (c *Cache) GetFromCacheOrFetchDB(ctx context.Context, key string, dest interface{}, fetchFromDB func() (interface{}, error), expDate time.Duration) error {
	// try cache with provided ctx
	if err := c.GetWithContext(ctx, key, dest); err == nil {
		// cache hit
		return nil
	} else {
		// If GetWithContext returned an error other than redis.Nil, log it and continue to fetch from DB
		if err != redis.Nil {
			utils.LogErrorWithLevel("warn",
				utils.DragonflyFailedToWriteCache.Type,
				utils.DragonflyFailedToWriteCache.Code,
				"cache read failed, falling back to DB",
				err,
			)
		}
	}

	// cache miss -> fetch from DB
	data, err := fetchFromDB()
	if err != nil {
		return fmt.Errorf("failed to fetch data from DB: %w", err)
	}

	// attempt to write to cache (best-effort)
	if err := c.Set(key, data, expDate); err != nil {
		utils.LogErrorWithLevel("warn",
			utils.DragonflyFailedToWriteCache.Type,
			utils.DragonflyFailedToWriteCache.Code,
			utils.DragonflyFailedToWriteCache.Msg,
			err,
		)
	}

	// marshal/unmarshal into dest
	mData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal db result: %w", err)
	}
	if err := json.Unmarshal(mData, dest); err != nil {
		return fmt.Errorf("failed to unmarshal into dest: %w", err)
	}

	return nil
}

// buildKey adds prefix to the key
func (c *Cache) buildKey(key string) string {
	if c.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Set : stores a value with expiration
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return c.client.Set(c.ctx, c.buildKey(key), data, expiration).Err()
}

// Get : retrieves a value
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

// GetWithContext : retrieves a value with context
func (c *Cache) GetWithContext(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.buildKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a single key from cache
func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, c.buildKey(key)).Err()
}

// Invalidate removes multiple keys at once
func (c *Cache) Invalidate(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	builtKeys := make([]string, len(keys))
	for i, key := range keys {
		builtKeys[i] = c.buildKey(key)
	}
	return c.client.Del(c.ctx, builtKeys...).Err()
}

// Flush : clears all keys
func (c *Cache) Flush() error {
	if c.prefix == "" {
		return c.client.FlushDB(c.ctx).Err()
	}
	// atomic deletion using lua script
	script := `
        local keys = redis.call('KEYS', ARGV[1])
        if #keys > 0 then
            return redis.call('DEL', unpack(keys))
        end
        return 0
    `
	return c.client.Eval(c.ctx, script, []string{}, fmt.Sprintf("%s:*", c.prefix)).Err()
}
