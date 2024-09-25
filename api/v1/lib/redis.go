package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
    Rdb *redis.Client
    Ctx = context.Background()
)

func InitRedis() {
    Rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
}

func GetValue(key string) (string, error) {
    value, err := Rdb.Get(Ctx, key).Result()
    if err == redis.Nil {
        fmt.Printf("Key does not exist: %v", key)
        return "", nil
    } else if err != nil {
        fmt.Printf("Failed to get value from Redis: %v", err)
        return "", err
    }
    return value, nil
}

func SetValue(key string, value string, expiration time.Duration) error {
    err := Rdb.Set(Ctx, key, value, expiration).Err()
    if err != nil {
        fmt.Printf("Failed to set value in Redis: %v", err)
        return err
    }
    return nil
}
