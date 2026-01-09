package ratelimit

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// Limiter provides rate limiting functionality using a token bucket algorithm.
// It wraps golang.org/x/time/rate.Limiter with convenience methods.
type Limiter struct {
	limiter *rate.Limiter
	enabled bool
	rps     float64 // Requests per second
}

// New creates a new rate limiter with the specified requests per second.
// If rps is 0 or negative, rate limiting is disabled (unlimited).
//
// The rate limiter uses a token bucket algorithm:
//   - Tokens are added to the bucket at the specified rate (rps)
//   - Each operation consumes one token
//   - If no tokens available, the operation waits until a token is available
//
// Parameters:
//   - rps: Maximum requests per second (0 or negative = unlimited)
//
// Returns a configured Limiter instance.
func New(rps float64) *Limiter {
	if rps <= 0 {
		// Rate limiting disabled
		return &Limiter{
			enabled: false,
			rps:     0,
		}
	}

	// Create token bucket rate limiter
	// Burst allows small bursts above the rate (set to 1 for strict rate limiting)
	burst := 1
	if rps < 1 {
		// For fractional rates (e.g., 0.5 rps = 1 req per 2 seconds),
		// allow burst of 1 to permit the first request immediately
		burst = 1
	}

	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		enabled: true,
		rps:     rps,
	}
}

// Wait blocks until the rate limiter permits an operation to proceed.
// If rate limiting is disabled, this returns immediately.
// If the context is canceled before a token is available, returns an error.
//
// This is the primary method to use before performing a rate-limited operation.
//
// Example:
//
//	limiter := ratelimit.New(10.0) // 10 requests per second
//	if err := limiter.Wait(ctx); err != nil {
//	    return fmt.Errorf("rate limit wait failed: %w", err)
//	}
//	// Proceed with operation
func (l *Limiter) Wait(ctx context.Context) error {
	if !l.enabled {
		return nil // No rate limiting
	}

	return l.limiter.Wait(ctx)
}

// Allow checks if an operation can proceed immediately without waiting.
// Returns true if a token is available, false otherwise.
// This is useful for non-blocking checks.
//
// Example:
//
//	if limiter.Allow() {
//	    // Proceed immediately
//	} else {
//	    // Rate limit exceeded, handle accordingly
//	}
func (l *Limiter) Allow() bool {
	if !l.enabled {
		return true // No rate limiting
	}

	return l.limiter.Allow()
}

// Reserve reserves a token for an operation and returns when it will be available.
// This allows checking how long the caller must wait without actually waiting.
//
// Returns a Reservation that can be queried for delay time.
// For disabled limiters, returns nil to indicate unlimited rate.
func (l *Limiter) Reserve() *rate.Reservation {
	if !l.enabled {
		// Return nil for disabled limiters
		return nil
	}

	return l.limiter.Reserve()
}

// Enabled returns true if rate limiting is active, false if disabled.
func (l *Limiter) Enabled() bool {
	return l.enabled
}

// RPS returns the configured requests per second rate.
// Returns 0 if rate limiting is disabled.
func (l *Limiter) RPS() float64 {
	return l.rps
}

// String returns a human-readable representation of the rate limiter configuration.
func (l *Limiter) String() string {
	if !l.enabled {
		return "rate limiting disabled (unlimited)"
	}
	if l.rps < 1 {
		// Display as "1 request per N seconds" for fractional rates
		interval := time.Duration(1/l.rps) * time.Second
		return fmt.Sprintf("%.2f rps (1 request per %v)", l.rps, interval)
	}
	return fmt.Sprintf("%.2f rps", l.rps)
}
