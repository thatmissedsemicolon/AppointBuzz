package lib

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetValue(key string) (string, error) {
	value, err := Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key does not exist: %v", key)
		return "", nil
	} else if err != nil {
		log.Printf("Failed to get value from Redis: %v", err)
		return "", err
	}
	return value, nil
}

func SetValue(key string, value string, expiration time.Duration) error {
	err := Rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set value in Redis: %v", err)
		return err
	}
	return nil
}

func DeleteValue(key string) error {
    err := Rdb.Del(ctx, key).Err()
    if err != nil {
        log.Printf("Failed to delete value from Redis: %v", err)
        return err
    }
    return nil
}

func GetTTL(key string) (int64, error) {
	ttl, err := Rdb.TTL(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to get TTL from Redis: %v", err)
		return 0, err
	}
	return int64(ttl.Seconds()), nil
}
