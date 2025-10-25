package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Rdb *redis.Client

func New() error {
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
	_, err := Rdb.Ping(context.TODO()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	zap.L().Debug("Redis connected successfully")
	return nil
}
