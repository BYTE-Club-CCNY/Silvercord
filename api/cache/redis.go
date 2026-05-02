package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		log.Println("REDIS_URL not set, caching disabled")
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis unavailable (%v), caching disabled\n", err)
		return nil
	}

	log.Println("Redis connected at", addr)
	return rdb
}

func GetJSON[T any](ctx context.Context, rdb *redis.Client, key string) (T, bool) {
	var zero T
	if rdb == nil {
		return zero, false
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return zero, false
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return zero, false
	}

	return result, true
}

func SetJSON(ctx context.Context, rdb *redis.Client, key string, value any, ttl time.Duration) error {
	if rdb == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, ttl).Err()
}

func Delete(ctx context.Context, rdb *redis.Client, keys ...string) {
	if rdb == nil {
		return
	}
	rdb.Del(ctx, keys...)
}
