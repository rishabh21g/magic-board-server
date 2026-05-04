package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter creates a new instance of RateLimiter with the provided Redis client, request limit, and time window for rate limiting
func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// Allow checks if the user has exceeded the rate limit and returns true if the request is allowed, false otherwise
func (r *RateLimiter) Allow(ctx context.Context, userID string) bool {
	key := "rate:" + userID
	count, _ := r.client.Incr(ctx, key).Result()
	if count == 1 {
		r.client.Expire(ctx, key, r.window)
	}
	if count > int64(r.limit) {
		return false
	}
	return true
}
