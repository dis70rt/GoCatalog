package limiter

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type LimiterManager struct {
	Buckets map[string]*TokenBucket
	Rate int
	Capacity int
	mu sync.Mutex	
}

var manager LimiterManager
var once sync.Once

func GetLimiterManager(rate int, capacity int) *LimiterManager {
	once.Do(func() {
		manager = LimiterManager{
			Buckets: make(map[string]*TokenBucket),
			Rate: rate,
			Capacity: capacity,
		}
	})
	return &manager
}

func (lm *LimiterManager) GetBucket(clientIP string) *TokenBucket {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	bucket, exists := lm.Buckets[clientIP]
	if !exists {
		bucket = NewTokenBucket(lm.Rate, lm.Capacity)
		lm.Buckets[clientIP] = bucket
	}
	return bucket
}

func (lm *LimiterManager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		bucket := lm.GetBucket(clientIP)
		if !bucket.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate Limit Exceed"})
			return
		}
		c.Next()
	}
}