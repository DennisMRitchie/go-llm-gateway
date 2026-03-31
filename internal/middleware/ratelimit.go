package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter stores per-IP rate limiters.
type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rps      rate.Limit
	burst    int
}

var globalLimiter = &ipLimiter{
	limiters: make(map[string]*rate.Limiter),
	rps:      100,
	burst:    20,
}

func (il *ipLimiter) get(ip string) *rate.Limiter {
	il.mu.Lock()
	defer il.mu.Unlock()
	if l, ok := il.limiters[ip]; ok {
		return l
	}
	l := rate.NewLimiter(il.rps, il.burst)
	il.limiters[ip] = l
	return l
}

// RateLimit returns a Gin middleware that enforces per-IP rate limiting.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := globalLimiter.get(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded — try again shortly",
			})
			return
		}
		c.Next()
	}
}
