package redis

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CheckCache checks if the given key exists in the Redis cache.
func CheckCache(key string, c *gin.Context) bool {
	val, err := Rdb.Get(c.Request.Context(), key).Bytes()
	if err != nil || len(val) == 0 {
		zap.L().Debug("Cache miss", zap.String("key", key), zap.Any("val", val))
		return false
	}

	zap.L().Debug("Cache hit", zap.String("key", key))

	var data any
	if err := json.Unmarshal(val, &data); err != nil {
		return false
	}

	c.JSON(http.StatusOK, data)
	return true
}

func InvalidateCache(key string) {
	if err := Rdb.Del(context.Background(), key).Err(); err != nil {
		zap.L().Error("Failed to invalidate cache", zap.String("key", key), zap.Error(err))
	}

	zap.L().Debug("Cache invalidated", zap.String("key", key))
}
