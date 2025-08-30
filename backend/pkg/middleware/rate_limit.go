package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

type RateLimiterConfig struct {
	RequestsPerSecond int
	Burst             int
	CleanupInterval   time.Duration
	TTL               time.Duration
}

func getVisitor(ip string, rps int, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(rps), burst)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors(ttl time.Duration, interval time.Duration) {
	for {
		time.Sleep(interval)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > ttl {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimiterMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	if config.CleanupInterval == 0 {
		config.CleanupInterval = time.Minute
	}
	if config.TTL == 0 {
		config.TTL = 3 * time.Minute
	}

	go cleanupVisitors(config.TTL, config.CleanupInterval)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip, config.RequestsPerSecond, config.Burst)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
