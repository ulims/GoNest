package gonest

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Event represents a system event
type Event struct {
	Name      string                 `json:"name"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	ID        string                 `json:"id"`
}

// NewEvent creates a new event
func NewEvent(name string, data interface{}) *Event {
	return &Event{
		Name:      name,
		Data:      data,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
		ID:        generateEventID(),
	}
}

// WithMetadata adds metadata to the event
func (e *Event) WithMetadata(key string, value interface{}) *Event {
	e.Metadata[key] = value
	return e
}

// WithSource sets the event source
func (e *Event) WithSource(source string) *Event {
	e.Source = source
	return e
}

// EventListener represents an event listener
type EventListener func(ctx context.Context, event *Event) error

// EventListenerConfig holds configuration for event listeners
type EventListenerConfig struct {
	Priority int
	Async    bool
	Filter   func(*Event) bool
	Retry    RetryConfig
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     func(attempt int, delay time.Duration) time.Duration
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
		Backoff: func(attempt int, delay time.Duration) time.Duration {
			return delay * time.Duration(attempt*2)
		},
	}
}

// listenerEntry represents a registered listener
type listenerEntry struct {
	listener EventListener
	config   EventListenerConfig
}

// EventEmitter provides event emission and listening capabilities
type EventEmitter struct {
	listeners map[string][]*listenerEntry
	logger    *logrus.Logger
	mutex     sync.RWMutex
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter(logger *logrus.Logger) *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]*listenerEntry),
		logger:    logger,
	}
}

// On registers an event listener
func (ee *EventEmitter) On(eventName string, listener EventListener, config ...EventListenerConfig) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	var listenerConfig EventListenerConfig
	if len(config) > 0 {
		listenerConfig = config[0]
	} else {
		listenerConfig = EventListenerConfig{
			Priority: 0,
			Async:    false,
			Retry:    DefaultRetryConfig(),
		}
	}

	entry := &listenerEntry{
		listener: listener,
		config:   listenerConfig,
	}

	ee.listeners[eventName] = append(ee.listeners[eventName], entry)
	ee.logger.Infof("Registered listener for event: %s", eventName)
}

// Once registers an event listener that will be called only once
func (ee *EventEmitter) Once(eventName string, listener EventListener, config ...EventListenerConfig) {
	var called bool
	var mutex sync.Mutex

	wrappedListener := func(ctx context.Context, event *Event) error {
		mutex.Lock()
		defer mutex.Unlock()

		if called {
			return nil
		}
		called = true

		// Remove this listener after execution
		ee.Off(eventName, listener)
		return listener(ctx, event)
	}

	ee.On(eventName, wrappedListener, config...)
}

// Off removes an event listener
func (ee *EventEmitter) Off(eventName string, listener EventListener) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	listeners := ee.listeners[eventName]
	for i, entry := range listeners {
		if reflect.ValueOf(entry.listener).Pointer() == reflect.ValueOf(listener).Pointer() {
			ee.listeners[eventName] = append(listeners[:i], listeners[i+1:]...)
			ee.logger.Infof("Removed listener for event: %s", eventName)
			break
		}
	}
}

// RemoveAllListeners removes all listeners for an event
func (ee *EventEmitter) RemoveAllListeners(eventName string) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	delete(ee.listeners, eventName)
	ee.logger.Infof("Removed all listeners for event: %s", eventName)
}

// Emit emits an event to all registered listeners
func (ee *EventEmitter) Emit(ctx context.Context, event *Event) error {
	ee.mutex.RLock()
	listeners := make([]*listenerEntry, len(ee.listeners[event.Name]))
	copy(listeners, ee.listeners[event.Name])
	ee.mutex.RUnlock()

	if len(listeners) == 0 {
		// No listeners for event: %s
		return nil
	}

	// Emitting event %s to %d listeners

	// Sort listeners by priority (higher priority first)
	for i := 0; i < len(listeners)-1; i++ {
		for j := i + 1; j < len(listeners); j++ {
			if listeners[i].config.Priority < listeners[j].config.Priority {
				listeners[i], listeners[j] = listeners[j], listeners[i]
			}
		}
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(listeners))

	for _, entry := range listeners {
		// Check filter condition
		if entry.config.Filter != nil && !entry.config.Filter(event) {
			continue
		}

		if entry.config.Async {
			wg.Add(1)
			go func(e *listenerEntry) {
				defer wg.Done()
				if err := ee.executeListener(ctx, e, event); err != nil {
					errors <- err
				}
			}(entry)
		} else {
			if err := ee.executeListener(ctx, entry, event); err != nil {
				return err
			}
		}
	}

	// Wait for async listeners if any
	if len(listeners) > 0 {
		go func() {
			wg.Wait()
			close(errors)
		}()
	}

	// Collect async errors
	for err := range errors {
		if err != nil {
			ee.logger.WithError(err).Error("Async listener error")
		}
	}

	return nil
}

// executeListener executes a listener with retry logic
func (ee *EventEmitter) executeListener(ctx context.Context, entry *listenerEntry, event *Event) error {
	var lastErr error
	delay := entry.config.Retry.Delay

	for attempt := 1; attempt <= entry.config.Retry.MaxAttempts; attempt++ {
		err := entry.listener(ctx, event)
		if err == nil {
			return nil
		}

		lastErr = err
		ee.logger.WithError(err).Warnf("Listener failed on attempt %d for event %s", attempt, event.Name)

		if attempt < entry.config.Retry.MaxAttempts {
			time.Sleep(delay)
			if entry.config.Retry.Backoff != nil {
				delay = entry.config.Retry.Backoff(attempt, delay)
			}
		}
	}

	return fmt.Errorf("listener failed after %d attempts: %w", entry.config.Retry.MaxAttempts, lastErr)
}

// EventBus provides global event management
type EventBus struct {
	emitters map[string]*EventEmitter
	logger   *logrus.Logger
	mutex    sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus(logger *logrus.Logger) *EventBus {
	return &EventBus{
		emitters: make(map[string]*EventEmitter),
		logger:   logger,
	}
}

// GetEmitter returns an event emitter for a namespace
func (eb *EventBus) GetEmitter(namespace string) *EventEmitter {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if emitter, exists := eb.emitters[namespace]; exists {
		return emitter
	}

	emitter := NewEventEmitter(eb.logger)
	eb.emitters[namespace] = emitter
	eb.logger.Infof("Created event emitter for namespace: %s", namespace)
	return emitter
}

// Emit emits an event to a specific namespace
func (eb *EventBus) Emit(ctx context.Context, namespace string, event *Event) error {
	emitter := eb.GetEmitter(namespace)
	return emitter.Emit(ctx, event)
}

// On registers a listener in a specific namespace
func (eb *EventBus) On(namespace, eventName string, listener EventListener, config ...EventListenerConfig) {
	emitter := eb.GetEmitter(namespace)
	emitter.On(eventName, listener, config...)
}

// EmitGlobal emits an event to all namespaces
func (eb *EventBus) EmitGlobal(ctx context.Context, event *Event) error {
	eb.mutex.RLock()
	emitters := make([]*EventEmitter, 0, len(eb.emitters))
	for _, emitter := range eb.emitters {
		emitters = append(emitters, emitter)
	}
	eb.mutex.RUnlock()

	for _, emitter := range emitters {
		if err := emitter.Emit(ctx, event); err != nil {
			eb.logger.WithError(err).Error("Failed to emit global event")
		}
	}

	return nil
}

// EventHandler decorator for automatic event listener registration
type EventHandler struct {
	EventName string
	Config    EventListenerConfig
}

// EventService provides high-level event operations
type EventService struct {
	eventBus  *EventBus
	namespace string
	logger    *logrus.Logger
}

// NewEventService creates a new event service
func NewEventService(eventBus *EventBus, namespace string, logger *logrus.Logger) *EventService {
	return &EventService{
		eventBus:  eventBus,
		namespace: namespace,
		logger:    logger,
	}
}

// Emit emits an event
func (es *EventService) Emit(ctx context.Context, eventName string, data interface{}) error {
	event := NewEvent(eventName, data).WithSource(es.namespace)
	return es.eventBus.Emit(ctx, es.namespace, event)
}

// EmitAsync emits an event asynchronously
func (es *EventService) EmitAsync(ctx context.Context, eventName string, data interface{}) {
	go func() {
		if err := es.Emit(ctx, eventName, data); err != nil {
			es.logger.WithError(err).Error("Failed to emit async event")
		}
	}()
}

// On registers an event listener
func (es *EventService) On(eventName string, listener EventListener, config ...EventListenerConfig) {
	es.eventBus.On(es.namespace, eventName, listener, config...)
}

// Once registers a one-time event listener
func (es *EventService) Once(eventName string, listener EventListener, config ...EventListenerConfig) {
	emitter := es.eventBus.GetEmitter(es.namespace)
	emitter.Once(eventName, listener, config...)
}

// Pre-defined event types for common scenarios
const (
	// Application events
	EventApplicationError = "application.error"

	// User events
	EventUserCreated         = "user.created"
	EventUserUpdated         = "user.updated"
	EventUserDeleted         = "user.deleted"
	EventUserLogin           = "user.login"
	EventUserLogout          = "user.logout"
	EventUserPasswordChanged = "user.password.changed"

	// Database events
	EventDatabaseConnected    = "database.connected"
	EventDatabaseDisconnected = "database.disconnected"
	EventDatabaseError        = "database.error"

	// HTTP events
	EventHTTPRequest  = "http.request"
	EventHTTPResponse = "http.response"
	EventHTTPError    = "http.error"

	// Cache events
	EventCacheHit   = "cache.hit"
	EventCacheMiss  = "cache.miss"
	EventCacheEvict = "cache.evict"

	// WebSocket events
	EventWebSocketConnect    = "websocket.connect"
	EventWebSocketDisconnect = "websocket.disconnect"
	EventWebSocketMessage    = "websocket.message"
	EventWebSocketError      = "websocket.error"
)

// EventData structures for common events
type UserEventData struct {
	UserID   string                 `json:"user_id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Action   string                 `json:"action"`
	Metadata map[string]interface{} `json:"metadata"`
}

type HTTPEventData struct {
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	StatusCode int               `json:"status_code"`
	Duration   time.Duration     `json:"duration"`
	UserAgent  string            `json:"user_agent"`
	IP         string            `json:"ip"`
	Headers    map[string]string `json:"headers"`
}

type DatabaseEventData struct {
	Operation string        `json:"operation"`
	Table     string        `json:"table"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
}

// EventMiddleware creates middleware that emits HTTP events
func EventMiddleware(eventService *EventService) func(next interface{}) interface{} {
	return func(next interface{}) interface{} {
		return func(ctx context.Context, data interface{}) error {
			// Emit request event
			eventService.EmitAsync(ctx, EventHTTPRequest, data)

			// Execute next handler
			start := time.Now()
			var err error
			if handler, ok := next.(func(context.Context, interface{}) error); ok {
				err = handler(ctx, data)
			}
			duration := time.Since(start)

			// Emit response/error event
			if err != nil {
				errorEvent := NewEvent(EventHTTPError, map[string]interface{}{
					"error":    err.Error(),
					"duration": duration,
					"data":     data,
				})
				eventService.EmitAsync(ctx, EventHTTPError, errorEvent.Data)
			} else {
				responseEvent := NewEvent(EventHTTPResponse, map[string]interface{}{
					"duration": duration,
					"data":     data,
				})
				eventService.EmitAsync(ctx, EventHTTPResponse, responseEvent.Data)
			}

			return err
		}
	}
}

// Utility functions

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}

// EventPattern represents an event pattern for pattern matching
type EventPattern struct {
	Pattern string
	Handler EventListener
}

// PatternEventEmitter extends EventEmitter with pattern matching
type PatternEventEmitter struct {
	*EventEmitter
	patterns []EventPattern
}

// NewPatternEventEmitter creates a new pattern event emitter
func NewPatternEventEmitter(logger *logrus.Logger) *PatternEventEmitter {
	return &PatternEventEmitter{
		EventEmitter: NewEventEmitter(logger),
		patterns:     make([]EventPattern, 0),
	}
}

// OnPattern registers a pattern-based event listener
func (pee *PatternEventEmitter) OnPattern(pattern string, handler EventListener) {
	pee.patterns = append(pee.patterns, EventPattern{
		Pattern: pattern,
		Handler: handler,
	})
}

// Emit emits an event and checks pattern matches
func (pee *PatternEventEmitter) Emit(ctx context.Context, event *Event) error {
	// First emit to exact listeners
	if err := pee.EventEmitter.Emit(ctx, event); err != nil {
		return err
	}

	// Then check pattern matches
	for _, pattern := range pee.patterns {
		if pee.matchesPattern(event.Name, pattern.Pattern) {
			if err := pattern.Handler(ctx, event); err != nil {
				pee.logger.WithError(err).Errorf("Pattern handler failed for event %s", event.Name)
			}
		}
	}

	return nil
}

// matchesPattern checks if an event name matches a pattern
func (pee *PatternEventEmitter) matchesPattern(eventName, pattern string) bool {
	// Simple wildcard matching (*, **)
	// This is a simplified implementation
	// A full implementation would use proper glob matching
	if pattern == "*" {
		return true
	}

	// For now, just do exact matching
	return eventName == pattern
}
