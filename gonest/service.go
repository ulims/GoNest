package gonest

import (
	"reflect"
	"sync"
)

// Service represents a NestJS-like service
type Service struct {
	Name      string
	Instance  interface{}
	Type      reflect.Type
	Value     reflect.Value
	Singleton bool
	Lazy      bool
}

// ServiceRegistry manages all services and their dependencies
type ServiceRegistry struct {
	services map[string]*Service
	mutex    sync.RWMutex
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*Service),
	}
}

// Register registers a service in the registry
func (sr *ServiceRegistry) Register(name string, service interface{}) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	sr.services[name] = &Service{
		Name:      name,
		Instance:  service,
		Type:      reflect.TypeOf(service),
		Value:     reflect.ValueOf(service),
		Singleton: true,
		Lazy:      false,
	}
}

// RegisterLazy registers a lazy service (created on first use)
func (sr *ServiceRegistry) RegisterLazy(name string, serviceType reflect.Type) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	sr.services[name] = &Service{
		Name:      name,
		Type:      serviceType,
		Singleton: true,
		Lazy:      true,
	}
}

// Get retrieves a service by name
func (sr *ServiceRegistry) Get(name string) (interface{}, bool) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	service, exists := sr.services[name]
	if !exists {
		return nil, false
	}

	// If lazy service, create instance
	if service.Lazy && service.Instance == nil {
		sr.mutex.RUnlock()
		sr.mutex.Lock()
		defer sr.mutex.RLock()

		// Double-check after acquiring write lock
		if service.Instance == nil {
			service.Instance = reflect.New(service.Type.Elem()).Interface()
		}
	}

	return service.Instance, true
}

// GetByType retrieves a service by type
func (sr *ServiceRegistry) GetByType(serviceType reflect.Type) (interface{}, bool) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	for _, service := range sr.services {
		if service.Type == serviceType {
			if service.Lazy && service.Instance == nil {
				sr.mutex.RUnlock()
				sr.mutex.Lock()
				defer sr.mutex.RLock()

				if service.Instance == nil {
					service.Instance = reflect.New(service.Type.Elem()).Interface()
				}
			}
			return service.Instance, true
		}
	}
	return nil, false
}

// GetAll returns all registered services
func (sr *ServiceRegistry) GetAll() map[string]*Service {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	result := make(map[string]*Service)
	for k, v := range sr.services {
		result[k] = v
	}
	return result
}

// Inject injects dependencies into a service
func (sr *ServiceRegistry) Inject(service interface{}) error {
	value := reflect.ValueOf(service)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Inject field dependencies
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := field.Type()

		// Check if field has inject tag
		if field.CanSet() {
			// Try to get service by type
			if dependency, exists := sr.GetByType(fieldType); exists {
				field.Set(reflect.ValueOf(dependency))
			}
		}
	}

	return nil
}

// ServiceBuilder provides a fluent interface for building services
type ServiceBuilder struct {
	service *Service
}

// NewService creates a new service
func NewService(name string) *ServiceBuilder {
	return &ServiceBuilder{
		service: &Service{
			Name:      name,
			Singleton: true,
			Lazy:      false,
		},
	}
}

// Instance sets the service instance
func (sb *ServiceBuilder) Instance(instance interface{}) *ServiceBuilder {
	sb.service.Instance = instance
	sb.service.Type = reflect.TypeOf(instance)
	sb.service.Value = reflect.ValueOf(instance)
	return sb
}

// Type sets the service type for lazy loading
func (sb *ServiceBuilder) Type(serviceType reflect.Type) *ServiceBuilder {
	sb.service.Type = serviceType
	sb.service.Lazy = true
	return sb
}

// Singleton sets whether the service is a singleton
func (sb *ServiceBuilder) Singleton(singleton bool) *ServiceBuilder {
	sb.service.Singleton = singleton
	return sb
}

// Lazy sets whether the service should be lazy loaded
func (sb *ServiceBuilder) Lazy(lazy bool) *ServiceBuilder {
	sb.service.Lazy = lazy
	return sb
}

// Build returns the built service
func (sb *ServiceBuilder) Build() *Service {
	return sb.service
}

// Injectable decorator for marking services as injectable
type Injectable struct{}

// InjectableDecoratorFunc creates an injectable decorator
func InjectableDecoratorFunc() Injectable {
	return Injectable{}
}

// InjectDecorator represents dependency injection
type InjectDecorator struct {
	Name string
}

// InjectDecoratorFunc creates an inject decorator
func InjectDecoratorFunc(name string) InjectDecorator {
	return InjectDecorator{Name: name}
}

// Service decorator for marking services
type ServiceDecorator struct{}

// ServiceDecoratorFunc creates a service decorator
func ServiceDecoratorFunc() ServiceDecorator {
	return ServiceDecorator{}
}

// Provider interface for custom service providers
type Provider interface {
	Provide() interface{}
}

// FactoryProvider provides services through a factory function
type FactoryProvider struct {
	Factory func() interface{}
}

// Provide implements the Provider interface
func (fp *FactoryProvider) Provide() interface{} {
	return fp.Factory()
}

// NewFactoryProvider creates a new factory provider
func NewFactoryProvider(factory func() interface{}) *FactoryProvider {
	return &FactoryProvider{Factory: factory}
}
