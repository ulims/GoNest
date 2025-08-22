package gonest

import (
	"sync"
)

// Module represents a NestJS-like module that can contain controllers, services, and other modules
type Module struct {
	Name        string
	Controllers []interface{}
	Services    []interface{}
	Modules     []*Module
	Providers   []interface{}
	Imports     []*Module
	Exports     []interface{}
}

// ModuleBuilder provides a fluent interface for building modules
type ModuleBuilder struct {
	module *Module
}

// NewModule creates a new module with the given name
func NewModule(name string) *ModuleBuilder {
	return &ModuleBuilder{
		module: &Module{
			Name:        name,
			Controllers: make([]interface{}, 0),
			Services:    make([]interface{}, 0),
			Modules:     make([]*Module, 0),
			Providers:   make([]interface{}, 0),
			Imports:     make([]*Module, 0),
			Exports:     make([]interface{}, 0),
		},
	}
}

// Controller adds a controller to the module
func (mb *ModuleBuilder) Controller(controller interface{}) *ModuleBuilder {
	mb.module.Controllers = append(mb.module.Controllers, controller)
	return mb
}

// Service adds a service to the module
func (mb *ModuleBuilder) Service(service interface{}) *ModuleBuilder {
	mb.module.Services = append(mb.module.Services, service)
	return mb
}

// Module adds a sub-module to the module
func (mb *ModuleBuilder) Module(subModule *Module) *ModuleBuilder {
	mb.module.Modules = append(mb.module.Modules, subModule)
	return mb
}

// Provider adds a provider to the module
func (mb *ModuleBuilder) Provider(provider interface{}) *ModuleBuilder {
	mb.module.Providers = append(mb.module.Providers, provider)
	return mb
}

// Import adds an import to the module
func (mb *ModuleBuilder) Import(importModule *Module) *ModuleBuilder {
	mb.module.Imports = append(mb.module.Imports, importModule)
	return mb
}

// Export adds an export to the module
func (mb *ModuleBuilder) Export(export interface{}) *ModuleBuilder {
	mb.module.Exports = append(mb.module.Exports, export)
	return mb
}

// Build returns the built module
func (mb *ModuleBuilder) Build() *Module {
	return mb.module
}

// ModuleRegistry manages all modules in the application
type ModuleRegistry struct {
	modules map[string]*Module
	mutex   sync.RWMutex
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]*Module),
	}
}

// Register registers a module in the registry
func (mr *ModuleRegistry) Register(module *Module) {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	mr.modules[module.Name] = module
}

// Get retrieves a module by name
func (mr *ModuleRegistry) Get(name string) (*Module, bool) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	module, exists := mr.modules[name]
	return module, exists
}

// GetAll returns all registered modules
func (mr *ModuleRegistry) GetAll() map[string]*Module {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	result := make(map[string]*Module)
	for k, v := range mr.modules {
		result[k] = v
	}
	return result
}

// GetControllers returns all controllers from all modules
func (mr *ModuleRegistry) GetControllers() []interface{} {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	var controllers []interface{}
	for _, module := range mr.modules {
		controllers = append(controllers, module.Controllers...)
	}
	return controllers
}

// GetServices returns all services from all modules
func (mr *ModuleRegistry) GetServices() []interface{} {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	var services []interface{}
	for _, module := range mr.modules {
		services = append(services, module.Services...)
	}
	return services
}

// GetProviders returns all providers from all modules
func (mr *ModuleRegistry) GetProviders() []interface{} {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	var providers []interface{}
	for _, module := range mr.modules {
		providers = append(providers, module.Providers...)
	}
	return providers
}
