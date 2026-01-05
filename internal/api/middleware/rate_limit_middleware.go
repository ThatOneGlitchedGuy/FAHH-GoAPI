package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang-app/internal/errors"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requestCounts map[string]*clientQuota
	mu            sync.RWMutex
	maxRequests   int
	windowSize    time.Duration
	cleanupTicker *time.Ticker
}

type clientQuota struct {
	requests   []time.Time
	lastActive time.Time
}

func NewRateLimiter(maxRequests int, windowSize time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requestCounts: make(map[string]*clientQuota),
		maxRequests:   maxRequests,
		windowSize:    windowSize,
		cleanupTicker: time.NewTicker(5 * time.Minute),
	}

	go rl.cleanupOldEntries()
	return rl
}

func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	quota, exists := rl.requestCounts[clientID]

	if !exists {
		rl.requestCounts[clientID] = &clientQuota{
			requests:   []time.Time{now},
			lastActive: now,
		}
		return true
	}

	quota.lastActive = now
	cutoff := now.Add(-rl.windowSize)
	var validRequests []time.Time

	for _, reqTime := range quota.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	if len(validRequests) >= rl.maxRequests {
		return false
	}

	quota.requests = append(validRequests, now)
	return true
}

func (rl *RateLimiter) cleanupOldEntries() {
	for range rl.cleanupTicker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientID, quota := range rl.requestCounts {
			if now.Sub(quota.lastActive) > 1*time.Hour {
				delete(rl.requestCounts, clientID)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
}

var defaultRateLimiter = NewRateLimiter(100, 1*time.Minute)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !defaultRateLimiter.Allow(clientIP) {
			appErr := errors.RateLimitError()
			c.JSON(http.StatusTooManyRequests, appErr)
			c.Abort()
			return
		}
		c.Next()
	}
}
