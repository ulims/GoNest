package gonest

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

// LifecycleHook interface for application lifecycle events
type LifecycleHook interface {
	OnApplicationStart(ctx context.Context) error
	OnApplicationStop(ctx context.Context) error
	OnModuleInit(ctx context.Context) error
	OnModuleDestroy(ctx context.Context) error
}

// LifecycleHookFunc is a function type that implements LifecycleHook interface
type LifecycleHookFunc func(ctx context.Context) error

// OnApplicationStart implements LifecycleHook interface
func (f LifecycleHookFunc) OnApplicationStart(ctx context.Context) error {
	return f(ctx)
}

// OnApplicationStop implements LifecycleHook interface
func (f LifecycleHookFunc) OnApplicationStop(ctx context.Context) error {
	return f(ctx)
}

// OnModuleInit implements LifecycleHook interface
func (f LifecycleHookFunc) OnModuleInit(ctx context.Context) error {
	return f(ctx)
}

// OnModuleDestroy implements LifecycleHook interface
func (f LifecycleHookFunc) OnModuleDestroy(ctx context.Context) error {
	return f(ctx)
}

// LifecycleEvent represents a lifecycle event
type LifecycleEvent string

const (
	EventApplicationStart   LifecycleEvent = "application_start"
	EventApplicationStop    LifecycleEvent = "application_stop"
	EventModuleInit         LifecycleEvent = "module_init"
	EventModuleDestroy      LifecycleEvent = "module_destroy"
	EventServiceStart       LifecycleEvent = "service_start"
	EventServiceStop        LifecycleEvent = "service_stop"
	EventDatabaseConnect    LifecycleEvent = "database_connect"
	EventDatabaseDisconnect LifecycleEvent = "database_disconnect"
)

// LifecycleHookMetadata contains hook information
type LifecycleHookMetadata struct {
	Hook     LifecycleHook
	Priority int
	Event    LifecycleEvent
	Metadata map[string]interface{}
}

// LifecycleManager manages application lifecycle hooks
type LifecycleManager struct {
	hooks  map[LifecycleEvent][]*LifecycleHookMetadata
	logger *logrus.Logger
	mu     sync.RWMutex
}

// NewLifecycleManager creates a new lifecycle manager
func NewLifecycleManager(logger *logrus.Logger) *LifecycleManager {
	return &LifecycleManager{
		hooks:  make(map[LifecycleEvent][]*LifecycleHookMetadata),
		logger: logger,
	}
}

// RegisterHook registers a lifecycle hook
func (lm *LifecycleManager) RegisterHook(event LifecycleEvent, hook LifecycleHook, priority int) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	metadata := &LifecycleHookMetadata{
		Hook:     hook,
		Priority: priority,
		Event:    event,
		Metadata: make(map[string]interface{}),
	}

	lm.hooks[event] = append(lm.hooks[event], metadata)
	lm.logger.Infof("Registered lifecycle hook for event: %s with priority: %d", event, priority)
}

// TriggerEvent triggers a lifecycle event
func (lm *LifecycleManager) TriggerEvent(ctx context.Context, event LifecycleEvent) error {
	lm.mu.RLock()
	hooks := lm.hooks[event]
	lm.mu.RUnlock()

	if len(hooks) == 0 {
		// No hooks registered for event: %s
		return nil
	}

	lm.logger.Infof("Triggering lifecycle event: %s with %d hooks", event, len(hooks))

	for _, hook := range hooks {
		var err error
		switch event {
		case EventApplicationStart:
			err = hook.Hook.OnApplicationStart(ctx)
		case EventApplicationStop:
			err = hook.Hook.OnApplicationStop(ctx)
		case EventModuleInit:
			err = hook.Hook.OnModuleInit(ctx)
		case EventModuleDestroy:
			err = hook.Hook.OnModuleDestroy(ctx)
		default:
			lm.logger.Warnf("Unknown lifecycle event: %s", event)
			continue
		}

		if err != nil {
			lm.logger.Errorf("Lifecycle hook failed for event %s: %v", event, err)
			return err
		}
	}

	lm.logger.Infof("Successfully triggered lifecycle event: %s", event)
	return nil
}

// DatabaseLifecycleHook provides database lifecycle management
type DatabaseLifecycleHook struct {
	database DatabaseInterface
	logger   *logrus.Logger
}

// NewDatabaseLifecycleHook creates a new database lifecycle hook
func NewDatabaseLifecycleHook(database DatabaseInterface, logger *logrus.Logger) *DatabaseLifecycleHook {
	return &DatabaseLifecycleHook{
		database: database,
		logger:   logger,
	}
}

// OnApplicationStart connects to the database
func (dlh *DatabaseLifecycleHook) OnApplicationStart(ctx context.Context) error {
	dlh.logger.Info("Connecting to database...")
	return dlh.database.Connect()
}

// OnApplicationStop disconnects from the database
func (dlh *DatabaseLifecycleHook) OnApplicationStop(ctx context.Context) error {
	dlh.logger.Info("Disconnecting from database...")
	return dlh.database.Disconnect()
}

// OnModuleInit initializes database module
func (dlh *DatabaseLifecycleHook) OnModuleInit(ctx context.Context) error {
	dlh.logger.Info("Initializing database module...")
	return nil
}

// OnModuleDestroy destroys database module
func (dlh *DatabaseLifecycleHook) OnModuleDestroy(ctx context.Context) error {
	dlh.logger.Info("Destroying database module...")
	return nil
}

// CacheLifecycleHook provides cache lifecycle management
type CacheLifecycleHook struct {
	cache  CacheService
	logger *logrus.Logger
}

// NewCacheLifecycleHook creates a new cache lifecycle hook
func NewCacheLifecycleHook(cache CacheService, logger *logrus.Logger) *CacheLifecycleHook {
	return &CacheLifecycleHook{
		cache:  cache,
		logger: logger,
	}
}

// OnApplicationStart initializes cache
func (clh *CacheLifecycleHook) OnApplicationStart(ctx context.Context) error {
	clh.logger.Info("Initializing cache...")
	return nil
}

// OnApplicationStop cleans up cache
func (clh *CacheLifecycleHook) OnApplicationStop(ctx context.Context) error {
	clh.logger.Info("Cleaning up cache...")
	return nil
}

// OnModuleInit initializes cache module
func (clh *CacheLifecycleHook) OnModuleInit(ctx context.Context) error {
	clh.logger.Info("Initializing cache module...")
	return nil
}

// OnModuleDestroy destroys cache module
func (clh *CacheLifecycleHook) OnModuleDestroy(ctx context.Context) error {
	clh.logger.Info("Destroying cache module...")
	return nil
}

// MetricsLifecycleHook provides metrics lifecycle management
type MetricsLifecycleHook struct {
	metrics MetricsService
	logger  *logrus.Logger
}

// NewMetricsLifecycleHook creates a new metrics lifecycle hook
func NewMetricsLifecycleHook(metrics MetricsService, logger *logrus.Logger) *MetricsLifecycleHook {
	return &MetricsLifecycleHook{
		metrics: metrics,
		logger:  logger,
	}
}

// OnApplicationStart starts metrics collection
func (mlh *MetricsLifecycleHook) OnApplicationStart(ctx context.Context) error {
	mlh.logger.Info("Starting metrics collection...")
	return nil
}

// OnApplicationStop stops metrics collection
func (mlh *MetricsLifecycleHook) OnApplicationStop(ctx context.Context) error {
	mlh.logger.Info("Stopping metrics collection...")
	return nil
}

// OnModuleInit initializes metrics module
func (mlh *MetricsLifecycleHook) OnModuleInit(ctx context.Context) error {
	mlh.logger.Info("Initializing metrics module...")
	return nil
}

// OnModuleDestroy destroys metrics module
func (mlh *MetricsLifecycleHook) OnModuleDestroy(ctx context.Context) error {
	mlh.logger.Info("Destroying metrics module...")
	return nil
}

// MongoDBLifecycleHook provides MongoDB lifecycle management
type MongoDBLifecycleHook struct {
	mongodb *MongoDBService
	logger  *logrus.Logger
}

// NewMongoDBLifecycleHook creates a new MongoDB lifecycle hook
func NewMongoDBLifecycleHook(mongodb *MongoDBService, logger *logrus.Logger) *MongoDBLifecycleHook {
	return &MongoDBLifecycleHook{
		mongodb: mongodb,
		logger:  logger,
	}
}

// OnApplicationStart connects to MongoDB
func (mlh *MongoDBLifecycleHook) OnApplicationStart(ctx context.Context) error {
	mlh.logger.Info("Connecting to MongoDB...")
	return mlh.mongodb.Connect()
}

// OnApplicationStop disconnects from MongoDB
func (mlh *MongoDBLifecycleHook) OnApplicationStop(ctx context.Context) error {
	mlh.logger.Info("Disconnecting from MongoDB...")
	return mlh.mongodb.Disconnect()
}

// OnModuleInit initializes MongoDB module
func (mlh *MongoDBLifecycleHook) OnModuleInit(ctx context.Context) error {
	mlh.logger.Info("Initializing MongoDB module...")
	return nil
}

// OnModuleDestroy destroys MongoDB module
func (mlh *MongoDBLifecycleHook) OnModuleDestroy(ctx context.Context) error {
	mlh.logger.Info("Destroying MongoDB module...")
	return nil
}

// WebSocketLifecycleHook provides WebSocket lifecycle management
type WebSocketLifecycleHook struct {
	// hub    *WebSocketHub // Commented out - use websocket.go implementation
	logger *logrus.Logger
}

// WebSocket lifecycle functionality would be implemented here

// OnApplicationStart starts WebSocket hub
func (wlh *WebSocketLifecycleHook) OnApplicationStart(ctx context.Context) error {
	wlh.logger.Info("WebSocket hub lifecycle managed in websocket.go")
	return nil
}

// OnApplicationStop stops WebSocket hub
func (wlh *WebSocketLifecycleHook) OnApplicationStop(ctx context.Context) error {
	wlh.logger.Info("Stopping WebSocket hub...")
	// Implementation would gracefully shutdown WebSocket connections
	return nil
}

// OnModuleInit initializes WebSocket module
func (wlh *WebSocketLifecycleHook) OnModuleInit(ctx context.Context) error {
	wlh.logger.Info("Initializing WebSocket module...")
	return nil
}

// OnModuleDestroy destroys WebSocket module
func (wlh *WebSocketLifecycleHook) OnModuleDestroy(ctx context.Context) error {
	wlh.logger.Info("Destroying WebSocket module...")
	return nil
}

// Lifecycle decorators
type LifecycleDecorator struct {
	Events []LifecycleEvent
}

// Lifecycle decorator for marking lifecycle hooks
func Lifecycle(events ...LifecycleEvent) LifecycleDecorator {
	return LifecycleDecorator{Events: events}
}

// OnStart decorator for application start hook
func OnStart() LifecycleDecorator {
	return Lifecycle(EventApplicationStart)
}

// OnStop decorator for application stop hook
func OnStop() LifecycleDecorator {
	return Lifecycle(EventApplicationStop)
}

// OnModuleInit decorator for module init hook
func OnModuleInit() LifecycleDecorator {
	return Lifecycle(EventModuleInit)
}

// OnModuleDestroy decorator for module destroy hook
func OnModuleDestroy() LifecycleDecorator {
	return Lifecycle(EventModuleDestroy)
}

// Built-in lifecycle hook priorities
const (
	PriorityHigh   = 100
	PriorityNormal = 50
	PriorityLow    = 10
)
