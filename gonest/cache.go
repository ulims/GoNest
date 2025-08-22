package gonest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CacheProvider interface for different cache implementations
type CacheProvider interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Exists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}

// CacheItem represents a cached item
type CacheItem struct {
	Value     []byte
	ExpiresAt time.Time
	CreatedAt time.Time
}

// IsExpired checks if the cache item has expired
func (ci *CacheItem) IsExpired() bool {
	return !ci.ExpiresAt.IsZero() && time.Now().After(ci.ExpiresAt)
}

// MemoryCache implements in-memory caching
type MemoryCache struct {
	items  map[string]*CacheItem
	mutex  sync.RWMutex
	logger *logrus.Logger
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(logger *logrus.Logger) *MemoryCache {
	mc := &MemoryCache{
		items:  make(map[string]*CacheItem),
		logger: logger,
	}

	// Start cleanup goroutine
	go mc.cleanup()

	return mc
}

// Get retrieves a value from cache
func (mc *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	item, exists := mc.items[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	if item.IsExpired() {
		go mc.Delete(context.Background(), key) // Async cleanup
		return nil, errors.New("key expired")
	}

	return item.Value, nil
}

// Set stores a value in cache
func (mc *MemoryCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	item := &CacheItem{
		Value:     value,
		CreatedAt: time.Now(),
	}

	if expiration > 0 {
		item.ExpiresAt = time.Now().Add(expiration)
	}

	mc.items[key] = item
	return nil
}

// Delete removes a value from cache
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	delete(mc.items, key)
	return nil
}

// Clear removes all items from cache
func (mc *MemoryCache) Clear(ctx context.Context) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.items = make(map[string]*CacheItem)
	return nil
}

// Exists checks if a key exists in cache
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	item, exists := mc.items[key]
	if !exists {
		return false, nil
	}

	if item.IsExpired() {
		go mc.Delete(context.Background(), key) // Async cleanup
		return false, nil
	}

	return true, nil
}

// TTL returns the time-to-live for a key
func (mc *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	item, exists := mc.items[key]
	if !exists {
		return 0, errors.New("key not found")
	}

	if item.ExpiresAt.IsZero() {
		return -1, nil // No expiration
	}

	if item.IsExpired() {
		return 0, nil
	}

	return time.Until(item.ExpiresAt), nil
}

// Keys returns all keys matching a pattern
func (mc *MemoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	var keys []string
	for key := range mc.items {
		// Simple pattern matching (for now, just return all keys if pattern is "*")
		if pattern == "*" || key == pattern {
			if item := mc.items[key]; !item.IsExpired() {
				keys = append(keys, key)
			}
		}
	}

	return keys, nil
}

// cleanup removes expired items periodically
func (mc *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mc.mutex.Lock()
		for key, item := range mc.items {
			if item.IsExpired() {
				delete(mc.items, key)
			}
		}
		mc.mutex.Unlock()
	}
}

// CacheService provides high-level caching operations
type CacheService struct {
	provider  CacheProvider
	logger    *logrus.Logger
	keyPrefix string
}

// NewCacheService creates a new cache service
func NewCacheService(provider CacheProvider, logger *logrus.Logger) *CacheService {
	return &CacheService{
		provider:  provider,
		logger:    logger,
		keyPrefix: "gonest:",
	}
}

// SetKeyPrefix sets the key prefix for all cache operations
func (cs *CacheService) SetKeyPrefix(prefix string) {
	cs.keyPrefix = prefix
}

// generateKey generates a cache key with prefix
func (cs *CacheService) generateKey(key string) string {
	return cs.keyPrefix + key
}

// Get retrieves and unmarshals a value from cache
func (cs *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := cs.provider.Get(ctx, cs.generateKey(key))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Set marshals and stores a value in cache
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cs.provider.Set(ctx, cs.generateKey(key), data, expiration)
}

// GetOrSet retrieves a value or sets it if not found
func (cs *CacheService) GetOrSet(ctx context.Context, key string, dest interface{}, provider func() (interface{}, error), expiration time.Duration) error {
	err := cs.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss, get from provider
	value, err := provider()
	if err != nil {
		return err
	}

	// Store in cache
	if err := cs.Set(ctx, key, value, expiration); err != nil {
		cs.logger.WithError(err).Warn("Failed to cache value")
	}

	// Set the destination value
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(valueBytes, dest)
}

// Delete removes a value from cache
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	return cs.provider.Delete(ctx, cs.generateKey(key))
}

// Clear removes all cached values
func (cs *CacheService) Clear(ctx context.Context) error {
	return cs.provider.Clear(ctx)
}

// Exists checks if a key exists in cache
func (cs *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	return cs.provider.Exists(ctx, cs.generateKey(key))
}

// Cache decorator configuration
type CacheConfig struct {
	Key        string
	TTL        time.Duration
	Condition  func(echo.Context) bool
	KeyBuilder func(echo.Context) string
}

// CacheInterceptor provides caching functionality for HTTP requests
type CacheInterceptor struct {
	cacheService *CacheService
	logger       *logrus.Logger
}

// NewCacheInterceptor creates a new cache interceptor
func NewCacheInterceptor(cacheService *CacheService, logger *logrus.Logger) *CacheInterceptor {
	return &CacheInterceptor{
		cacheService: cacheService,
		logger:       logger,
	}
}

// Middleware returns cache middleware
func (ci *CacheInterceptor) Middleware(config *CacheConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &CacheConfig{
			TTL: 5 * time.Minute,
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check condition if provided
			if config.Condition != nil && !config.Condition(c) {
				return next(c)
			}

			// Only cache GET requests
			if c.Request().Method != "GET" {
				return next(c)
			}

			// Generate cache key
			var cacheKey string
			if config.KeyBuilder != nil {
				cacheKey = config.KeyBuilder(c)
			} else if config.Key != "" {
				cacheKey = config.Key
			} else {
				cacheKey = fmt.Sprintf("http:%s:%s", c.Request().Method, c.Request().URL.Path)
				if c.Request().URL.RawQuery != "" {
					cacheKey += "?" + c.Request().URL.RawQuery
				}
			}

			// Try to get from cache
			var cachedResponse CachedResponse
			err := ci.cacheService.Get(c.Request().Context(), cacheKey, &cachedResponse)
			if err == nil {
				// Cache hit
				// Cache hit for key: %s
				c.Response().Header().Set("X-Cache", "HIT")

				// Set headers
				for key, value := range cachedResponse.Headers {
					c.Response().Header().Set(key, value)
				}

				return c.Blob(http.StatusOK, cachedResponse.ContentType, cachedResponse.Body)
			}

			// Cache miss, execute handler and cache response
			// Cache miss for key: %s
			c.Response().Header().Set("X-Cache", "MISS")

			// Create a custom response writer to capture the response
			rec := &ResponseRecorder{
				ResponseWriter: c.Response().Writer,
				statusCode:     200,
				headers:        make(map[string]string),
			}
			c.Response().Writer = rec

			// Execute the handler
			err = next(c)
			if err != nil {
				return err
			}

			// Cache the response if successful
			if rec.statusCode >= 200 && rec.statusCode < 300 {
				cachedResponse := CachedResponse{
					StatusCode:  rec.statusCode,
					ContentType: rec.headers["Content-Type"],
					Headers:     rec.headers,
					Body:        rec.body,
					CachedAt:    time.Now(),
				}

				if cacheErr := ci.cacheService.Set(c.Request().Context(), cacheKey, cachedResponse, config.TTL); cacheErr != nil {
					ci.logger.WithError(cacheErr).Warn("Failed to cache response")
				} else {
					// Cached response for key: %s
				}
			}

			return nil
		}
	}
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode  int               `json:"status_code"`
	ContentType string            `json:"content_type"`
	Headers     map[string]string `json:"headers"`
	Body        []byte            `json:"body"`
	CachedAt    time.Time         `json:"cached_at"`
}

// ResponseRecorder captures HTTP response data
type ResponseRecorder struct {
	ResponseWriter interface{}
	statusCode     int
	headers        map[string]string
	body           []byte
}

// Write captures response body
func (rr *ResponseRecorder) Write(data []byte) (int, error) {
	rr.body = append(rr.body, data...)
	if writer, ok := rr.ResponseWriter.(interface{ Write([]byte) (int, error) }); ok {
		return writer.Write(data)
	}
	return len(data), nil
}

// Header captures response headers
func (rr *ResponseRecorder) Header() http.Header {
	if writer, ok := rr.ResponseWriter.(interface{ Header() http.Header }); ok {
		return writer.Header()
	}
	return make(http.Header)
}

// WriteHeader captures status code and headers
func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode

	// Capture headers
	header := rr.Header()
	rr.headers["Content-Type"] = header.Get("Content-Type")

	if writer, ok := rr.ResponseWriter.(interface{ WriteHeader(int) }); ok {
		writer.WriteHeader(statusCode)
	}
}

// CacheManager provides cache management utilities
type CacheManager struct {
	services map[string]*CacheService
	logger   *logrus.Logger
	mutex    sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager(logger *logrus.Logger) *CacheManager {
	return &CacheManager{
		services: make(map[string]*CacheService),
		logger:   logger,
	}
}

// RegisterCache registers a cache service with a name
func (cm *CacheManager) RegisterCache(name string, service *CacheService) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.services[name] = service
	cm.logger.Infof("Registered cache service: %s", name)
}

// GetCache retrieves a cache service by name
func (cm *CacheManager) GetCache(name string) (*CacheService, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	service, exists := cm.services[name]
	if !exists {
		return nil, fmt.Errorf("cache service '%s' not found", name)
	}

	return service, nil
}

// ClearAll clears all registered caches
func (cm *CacheManager) ClearAll(ctx context.Context) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for name, service := range cm.services {
		if err := service.Clear(ctx); err != nil {
			cm.logger.WithError(err).Errorf("Failed to clear cache: %s", name)
			return err
		}
	}

	return nil
}

// Cacheable decorator function
func Cacheable(key string, ttl time.Duration) func(interface{}) interface{} {
	return func(fn interface{}) interface{} {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		if fnType.Kind() != reflect.Func {
			panic("Cacheable can only be applied to functions")
		}

		return reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
			// This is a simplified implementation
			// In a real implementation, you would:
			// 1. Generate cache key from function name and args
			// 2. Check cache for existing result
			// 3. If miss, call original function and cache result
			// 4. Return cached or computed result

			// For now, just call the original function
			return fnValue.Call(args)
		}).Interface()
	}
}

// CacheEvict decorator for cache invalidation
func CacheEvict(key string, allEntries bool) func(interface{}) interface{} {
	return func(fn interface{}) interface{} {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		if fnType.Kind() != reflect.Func {
			panic("CacheEvict can only be applied to functions")
		}

		return reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
			// Execute the original function
			results := fnValue.Call(args)

			// After execution, evict cache entries
			// This would require access to the cache service
			// Implementation would go here

			return results
		}).Interface()
	}
}

// CachePut decorator for updating cache
func CachePut(key string, condition func(interface{}) bool) func(interface{}) interface{} {
	return func(fn interface{}) interface{} {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		if fnType.Kind() != reflect.Func {
			panic("CachePut can only be applied to functions")
		}

		return reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
			// Execute the original function
			results := fnValue.Call(args)

			// After execution, update cache with result
			// This would require access to the cache service
			// Implementation would go here

			return results
		}).Interface()
	}
}
