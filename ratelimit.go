package gonest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// RateLimitStrategy defines different rate limiting strategies
type RateLimitStrategy string

const (
	FixedWindowStrategy   RateLimitStrategy = "fixed_window"
	SlidingWindowStrategy RateLimitStrategy = "sliding_window"
	TokenBucketStrategy   RateLimitStrategy = "token_bucket"
	LeakyBucketStrategy   RateLimitStrategy = "leaky_bucket"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Strategy       RateLimitStrategy
	MaxRequests    int64
	WindowDuration time.Duration
	KeyGenerator   func(echo.Context) string
	ErrorMessage   string
	Headers        bool
	SkipFunc       func(echo.Context) bool
	OnLimitReached func(echo.Context) error
	Store          RateLimitStore
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Strategy:       FixedWindowStrategy,
		MaxRequests:    100,
		WindowDuration: time.Hour,
		ErrorMessage:   "Rate limit exceeded",
		Headers:        true,
		KeyGenerator: func(c echo.Context) string {
			return getClientIP(c)
		},
		Store: NewMemoryRateLimitStore(),
	}
}

// RateLimitStore interface for storing rate limit data
type RateLimitStore interface {
	Get(ctx context.Context, key string) (*RateLimitEntry, error)
	Set(ctx context.Context, key string, entry *RateLimitEntry, expiration time.Duration) error
	Increment(ctx context.Context, key string, expiration time.Duration) (*RateLimitEntry, error)
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// RateLimitEntry represents a rate limit entry
type RateLimitEntry struct {
	Count      int64     `json:"count"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	ExpiresAt  time.Time `json:"expires_at"`
	Tokens     float64   `json:"tokens,omitempty"`      // For token bucket
	LastRefill time.Time `json:"last_refill,omitempty"` // For token bucket
}

// MemoryRateLimitStore implements in-memory rate limiting
type MemoryRateLimitStore struct {
	entries map[string]*RateLimitEntry
	mutex   sync.RWMutex
}

// NewMemoryRateLimitStore creates a new memory rate limit store
func NewMemoryRateLimitStore() *MemoryRateLimitStore {
	store := &MemoryRateLimitStore{
		entries: make(map[string]*RateLimitEntry),
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// Get retrieves a rate limit entry
func (mrs *MemoryRateLimitStore) Get(ctx context.Context, key string) (*RateLimitEntry, error) {
	mrs.mutex.RLock()
	defer mrs.mutex.RUnlock()

	entry, exists := mrs.entries[key]
	if !exists {
		return nil, fmt.Errorf("entry not found")
	}

	if time.Now().After(entry.ExpiresAt) {
		go mrs.Delete(context.Background(), key) // Async cleanup
		return nil, fmt.Errorf("entry expired")
	}

	return entry, nil
}

// Set stores a rate limit entry
func (mrs *MemoryRateLimitStore) Set(ctx context.Context, key string, entry *RateLimitEntry, expiration time.Duration) error {
	mrs.mutex.Lock()
	defer mrs.mutex.Unlock()

	if expiration > 0 {
		entry.ExpiresAt = time.Now().Add(expiration)
	}

	mrs.entries[key] = entry
	return nil
}

// Increment increments a rate limit entry
func (mrs *MemoryRateLimitStore) Increment(ctx context.Context, key string, expiration time.Duration) (*RateLimitEntry, error) {
	mrs.mutex.Lock()
	defer mrs.mutex.Unlock()

	now := time.Now()
	entry, exists := mrs.entries[key]

	if !exists || now.After(entry.ExpiresAt) {
		entry = &RateLimitEntry{
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
			ExpiresAt: now.Add(expiration),
		}
	} else {
		entry.Count++
		entry.LastSeen = now
	}

	mrs.entries[key] = entry
	return entry, nil
}

// Delete removes a rate limit entry
func (mrs *MemoryRateLimitStore) Delete(ctx context.Context, key string) error {
	mrs.mutex.Lock()
	defer mrs.mutex.Unlock()

	delete(mrs.entries, key)
	return nil
}

// Clear removes all entries
func (mrs *MemoryRateLimitStore) Clear(ctx context.Context) error {
	mrs.mutex.Lock()
	defer mrs.mutex.Unlock()

	mrs.entries = make(map[string]*RateLimitEntry)
	return nil
}

// cleanup removes expired entries periodically
func (mrs *MemoryRateLimitStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mrs.mutex.Lock()
		now := time.Now()
		for key, entry := range mrs.entries {
			if now.After(entry.ExpiresAt) {
				delete(mrs.entries, key)
			}
		}
		mrs.mutex.Unlock()
	}
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	config *RateLimitConfig
	logger *logrus.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config *RateLimitConfig, logger *logrus.Logger) *RateLimiter {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	return &RateLimiter{
		config: config,
		logger: logger,
	}
}

// Middleware returns rate limiting middleware
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rl.config.SkipFunc != nil && rl.config.SkipFunc(c) {
				return next(c)
			}

			key := rl.config.KeyGenerator(c)
			ctx := c.Request().Context()

			allowed, info, err := rl.checkLimit(ctx, key)
			if err != nil {
				rl.logger.WithError(err).Error("Rate limit check failed")
				return next(c) // Continue on error
			}

			// Set rate limit headers
			if rl.config.Headers {
				c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.MaxRequests))
				c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, rl.config.MaxRequests-info.Count)))
				c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", info.ExpiresAt.Unix()))

				if !allowed {
					retryAfter := time.Until(info.ExpiresAt).Seconds()
					c.Response().Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter))
				}
			}

			if !allowed {
				if rl.config.OnLimitReached != nil {
					return rl.config.OnLimitReached(c)
				}

				return echo.NewHTTPError(http.StatusTooManyRequests, rl.config.ErrorMessage)
			}

			return next(c)
		}
	}
}

// checkLimit checks if a request is within the rate limit
func (rl *RateLimiter) checkLimit(ctx context.Context, key string) (bool, *RateLimitEntry, error) {
	switch rl.config.Strategy {
	case FixedWindowStrategy:
		return rl.checkFixedWindow(ctx, key)
	case SlidingWindowStrategy:
		return rl.checkSlidingWindow(ctx, key)
	case TokenBucketStrategy:
		return rl.checkTokenBucket(ctx, key)
	case LeakyBucketStrategy:
		return rl.checkLeakyBucket(ctx, key)
	default:
		return rl.checkFixedWindow(ctx, key)
	}
}

// checkFixedWindow implements fixed window rate limiting
func (rl *RateLimiter) checkFixedWindow(ctx context.Context, key string) (bool, *RateLimitEntry, error) {
	entry, err := rl.config.Store.Increment(ctx, key, rl.config.WindowDuration)
	if err != nil {
		return false, nil, err
	}

	allowed := entry.Count <= rl.config.MaxRequests
	return allowed, entry, nil
}

// checkSlidingWindow implements sliding window rate limiting
func (rl *RateLimiter) checkSlidingWindow(ctx context.Context, key string) (bool, *RateLimitEntry, error) {
	// This is a simplified sliding window implementation
	// A full implementation would track individual request timestamps

	now := time.Now()
	entry, err := rl.config.Store.Get(ctx, key)
	if err != nil || now.After(entry.ExpiresAt) {
		// Create new entry
		entry = &RateLimitEntry{
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
			ExpiresAt: now.Add(rl.config.WindowDuration),
		}
		rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
		return true, entry, nil
	}

	// Calculate rate based on sliding window
	windowStart := now.Add(-rl.config.WindowDuration)
	if entry.FirstSeen.Before(windowStart) {
		// Adjust count based on sliding window
		elapsed := now.Sub(entry.FirstSeen)
		if elapsed >= rl.config.WindowDuration {
			entry.Count = 1
			entry.FirstSeen = now
		} else {
			// Linear approximation for sliding window
			ratio := float64(elapsed) / float64(rl.config.WindowDuration)
			entry.Count = int64(float64(entry.Count)*(1-ratio)) + 1
		}
	} else {
		entry.Count++
	}

	entry.LastSeen = now
	rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)

	allowed := entry.Count <= rl.config.MaxRequests
	return allowed, entry, nil
}

// checkTokenBucket implements token bucket rate limiting
func (rl *RateLimiter) checkTokenBucket(ctx context.Context, key string) (bool, *RateLimitEntry, error) {
	now := time.Now()
	entry, err := rl.config.Store.Get(ctx, key)

	if err != nil || now.After(entry.ExpiresAt) {
		// Create new entry with full bucket
		entry = &RateLimitEntry{
			Count:      1,
			FirstSeen:  now,
			LastSeen:   now,
			ExpiresAt:  now.Add(rl.config.WindowDuration),
			Tokens:     float64(rl.config.MaxRequests - 1),
			LastRefill: now,
		}
		rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
		return true, entry, nil
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(entry.LastRefill)
	refillRate := float64(rl.config.MaxRequests) / float64(rl.config.WindowDuration)
	tokensToAdd := refillRate * float64(elapsed)

	entry.Tokens = min(float64(rl.config.MaxRequests), entry.Tokens+tokensToAdd)
	entry.LastRefill = now

	// Check if token is available
	if entry.Tokens >= 1.0 {
		entry.Tokens--
		entry.Count++
		entry.LastSeen = now
		rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
		return true, entry, nil
	}

	// No tokens available
	rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
	return false, entry, nil
}

// checkLeakyBucket implements leaky bucket rate limiting
func (rl *RateLimiter) checkLeakyBucket(ctx context.Context, key string) (bool, *RateLimitEntry, error) {
	now := time.Now()
	entry, err := rl.config.Store.Get(ctx, key)

	if err != nil || now.After(entry.ExpiresAt) {
		// Create new entry
		entry = &RateLimitEntry{
			Count:      1,
			FirstSeen:  now,
			LastSeen:   now,
			ExpiresAt:  now.Add(rl.config.WindowDuration),
			LastRefill: now,
		}
		rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
		return true, entry, nil
	}

	// Leak tokens based on elapsed time
	elapsed := now.Sub(entry.LastRefill)
	leakRate := float64(rl.config.MaxRequests) / float64(rl.config.WindowDuration)
	tokensToLeak := leakRate * float64(elapsed)

	entry.Count = max(0, int64(float64(entry.Count)-tokensToLeak))
	entry.LastRefill = now

	// Check if bucket has space
	if entry.Count < rl.config.MaxRequests {
		entry.Count++
		entry.LastSeen = now
		rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
		return true, entry, nil
	}

	// Bucket is full
	rl.config.Store.Set(ctx, key, entry, rl.config.WindowDuration)
	return false, entry, nil
}

// Rate limiting decorators and utilities

// RateLimit decorator configuration
type RateLimitDecorator struct {
	MaxRequests    int64
	WindowDuration time.Duration
	Strategy       RateLimitStrategy
	KeyGenerator   func(echo.Context) string
}

// NewRateLimit creates a rate limit decorator
func NewRateLimit(maxRequests int64, windowDuration time.Duration) *RateLimitDecorator {
	return &RateLimitDecorator{
		MaxRequests:    maxRequests,
		WindowDuration: windowDuration,
		Strategy:       FixedWindowStrategy,
		KeyGenerator: func(c echo.Context) string {
			return getClientIP(c)
		},
	}
}

// WithStrategy sets the rate limiting strategy
func (rld *RateLimitDecorator) WithStrategy(strategy RateLimitStrategy) *RateLimitDecorator {
	rld.Strategy = strategy
	return rld
}

// WithKeyGenerator sets a custom key generator
func (rld *RateLimitDecorator) WithKeyGenerator(generator func(echo.Context) string) *RateLimitDecorator {
	rld.KeyGenerator = generator
	return rld
}

// Middleware returns the rate limiting middleware
func (rld *RateLimitDecorator) Middleware(logger *logrus.Logger) echo.MiddlewareFunc {
	config := &RateLimitConfig{
		Strategy:       rld.Strategy,
		MaxRequests:    rld.MaxRequests,
		WindowDuration: rld.WindowDuration,
		KeyGenerator:   rld.KeyGenerator,
		ErrorMessage:   "Rate limit exceeded",
		Headers:        true,
		Store:          NewMemoryRateLimitStore(),
	}

	rateLimiter := NewRateLimiter(config, logger)
	return rateLimiter.Middleware()
}

// Predefined rate limiting configurations

// PerMinute creates a rate limiter for requests per minute
func PerMinute(requests int64) *RateLimitDecorator {
	return NewRateLimit(requests, time.Minute)
}

// PerHour creates a rate limiter for requests per hour
func PerHour(requests int64) *RateLimitDecorator {
	return NewRateLimit(requests, time.Hour)
}

// PerDay creates a rate limiter for requests per day
func PerDay(requests int64) *RateLimitDecorator {
	return NewRateLimit(requests, 24*time.Hour)
}

// PerSecond creates a rate limiter for requests per second
func PerSecond(requests int64) *RateLimitDecorator {
	return NewRateLimit(requests, time.Second)
}

// Key generators

// IPKeyGenerator generates keys based on client IP
func IPKeyGenerator(c echo.Context) string {
	return "ip:" + getClientIP(c)
}

// UserKeyGenerator generates keys based on authenticated user
func UserKeyGenerator(c echo.Context) string {
	user, err := GetCurrentUser(c)
	if err != nil {
		return IPKeyGenerator(c) // Fallback to IP
	}
	return "user:" + user.ID
}

// RouteKeyGenerator generates keys based on route
func RouteKeyGenerator(c echo.Context) string {
	return "route:" + c.Path()
}

// CombinedKeyGenerator generates keys based on multiple factors
func CombinedKeyGenerator(c echo.Context) string {
	ip := getClientIP(c)
	route := c.Path()
	return fmt.Sprintf("combined:%s:%s", ip, route)
}

// getClientIP extracts client IP from request
func getClientIP(c echo.Context) string {
	// Try X-Forwarded-For header first
	forwarded := c.Request().Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try X-Real-IP header
	realIP := c.Request().Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return c.Request().RemoteAddr
	}

	return ip
}

// Advanced rate limiting features

// BurstRateLimiter allows burst requests within limits
type BurstRateLimiter struct {
	normalLimiter *RateLimiter
	burstLimiter  *RateLimiter
	logger        *logrus.Logger
}

// NewBurstRateLimiter creates a burst rate limiter
func NewBurstRateLimiter(normalConfig, burstConfig *RateLimitConfig, logger *logrus.Logger) *BurstRateLimiter {
	return &BurstRateLimiter{
		normalLimiter: NewRateLimiter(normalConfig, logger),
		burstLimiter:  NewRateLimiter(burstConfig, logger),
		logger:        logger,
	}
}

// Middleware returns burst rate limiting middleware
func (brl *BurstRateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check burst limit first
			burstAllowed, _, err := brl.burstLimiter.checkLimit(c.Request().Context(), brl.burstLimiter.config.KeyGenerator(c))
			if err != nil {
				brl.logger.WithError(err).Error("Burst rate limit check failed")
				return next(c) // Continue on error
			}

			if !burstAllowed {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Burst rate limit exceeded")
			}

			// Then check normal limit
			return brl.normalLimiter.Middleware()(next)(c)
		}
	}
}

// Utility functions
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
