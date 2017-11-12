package discollect

import (
	"context"
	"errors"
)

var (
	ErrRateLimitExceeded = errors.New("discollect: rate limit exceeded")
)

// A RateLimiter is used for per-site and per-config rate limits
// abstracted out into an interface so that distributed rate limiting
// is practical
type RateLimiter interface {
	// limit blocks until the rate limit is ok
	Limit(ctx context.Context, rl *RateLimit, url string) error
}
