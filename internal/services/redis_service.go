package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(redisURL string) *RedisService {
	// Parse Redis URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Failed to parse Redis URL, using defaults: %v", err)
		opts = &redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}
	}

	// Configure for high performance
	opts.PoolSize = 100
	opts.MinIdleConns = 10
	opts.MaxIdleConns = 50
	opts.ConnMaxIdleTime = 5 * time.Minute
	opts.ConnMaxLifetime = 10 * time.Minute

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️  Redis connection failed: %v", err)
		log.Println("📝 Redis caching will be disabled")
		return &RedisService{client: nil}
	}

	log.Println("✅ Redis connected successfully")
	return &RedisService{client: client}
}

// Set stores a key-value pair with expiration
func (r *RedisService) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("redis client not available")
	}

	return r.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisService) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("redis client not available")
	}

	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// CheckRateLimit implements rate limiting using sliding window
func (r *RedisService) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) bool {
	if r.client == nil {
		return true // Allow if Redis unavailable
	}

	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// Use pipeline for atomic operations
	pipe := r.client.Pipeline()

	// Remove expired entries
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count current requests in window
	countCmd := pipe.ZCard(ctx, key)

	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})

	// Set expiration
	pipe.Expire(ctx, key, window)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Rate limit check failed: %v", err)
		return true // Allow on error
	}

	// Check if limit exceeded
	count := countCmd.Val()
	return count < int64(limit)
}

// SetJSON stores a JSON object
func (r *RedisService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Set(ctx, key, value, expiration).Err()
}

// GetJSON retrieves a JSON object
func (r *RedisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Get(ctx, key).Scan(dest)
}

// Increment atomically increments a counter
func (r *RedisService) Increment(ctx context.Context, key string) (int64, error) {
	if r.client == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	return r.client.Incr(ctx, key).Result()
}

// SetExpiration sets expiration on existing key
func (r *RedisService) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Expire(ctx, key, expiration).Err()
}

// GetTTL gets time to live for a key
func (r *RedisService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	if r.client == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	return r.client.TTL(ctx, key).Result()
}

// LPush pushes to the left of a list
func (r *RedisService) LPush(ctx context.Context, key string, values ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.LPush(ctx, key, values...).Err()
}

// RPop pops from the right of a list
func (r *RedisService) RPop(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("redis client not available")
	}

	return r.client.RPop(ctx, key).Result()
}

// LLen gets list length
func (r *RedisService) LLen(ctx context.Context, key string) (int64, error) {
	if r.client == nil {
		return 0, fmt.Errorf("redis client not available")
	}

	return r.client.LLen(ctx, key).Result()
}

// SetNX sets a key only if it doesn't exist (for locks)
func (r *RedisService) SetNX(ctx context.Context, key, value string, expiration time.Duration) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("redis client not available")
	}

	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Lock creates a distributed lock
func (r *RedisService) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	return r.SetNX(ctx, lockKey, "locked", expiration)
}

// Unlock releases a distributed lock
func (r *RedisService) Unlock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return r.Delete(ctx, lockKey)
}

// FlushAll clears all Redis data (use with caution)
func (r *RedisService) FlushAll(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.FlushAll(ctx).Err()
}

// GetStats returns Redis connection stats
func (r *RedisService) GetStats() *redis.PoolStats {
	if r.client == nil {
		return nil
	}

	return r.client.PoolStats()
}

// Ping tests Redis connectivity
func (r *RedisService) Ping(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	if r.client == nil {
		return nil
	}

	return r.client.Close()
}

// CacheFlowExecution caches flow execution state
func (r *RedisService) CacheFlowExecution(ctx context.Context, executionID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("flow_execution:%s", executionID)
	return r.SetJSON(ctx, key, data, expiration)
}

// GetCachedFlowExecution retrieves cached flow execution state
func (r *RedisService) GetCachedFlowExecution(ctx context.Context, executionID string, dest interface{}) error {
	key := fmt.Sprintf("flow_execution:%s", executionID)
	return r.GetJSON(ctx, key, dest)
}

// CacheUserSession caches user session data
func (r *RedisService) CacheUserSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.Set(ctx, key, userID, expiration)
}

// GetCachedUserSession retrieves cached user session
func (r *RedisService) GetCachedUserSession(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return r.Get(ctx, key)
}

// IncrementMessageCount increments message count for analytics
func (r *RedisService) IncrementMessageCount(ctx context.Context, deviceID string) error {
	key := fmt.Sprintf("message_count:%s", deviceID)
	_, err := r.Increment(ctx, key)
	if err != nil {
		return err
	}

	// Set expiration if this is a new counter
	return r.SetExpiration(ctx, key, 24*time.Hour)
}

// GetMessageCount gets message count for a device
func (r *RedisService) GetMessageCount(ctx context.Context, deviceID string) (int64, error) {
	key := fmt.Sprintf("message_count:%s", deviceID)
	result, err := r.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	// Convert string to int64
	count := int64(0)
	fmt.Sscanf(result, "%d", &count)
	return count, nil
}
