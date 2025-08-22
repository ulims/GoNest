package gonest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ConfigProvider interface for different configuration sources
type ConfigProvider interface {
	Load() (map[string]interface{}, error)
	Watch(callback func(map[string]interface{})) error
	GetName() string
}

// EnvironmentConfigProvider loads configuration from environment variables
type EnvironmentConfigProvider struct {
	prefix string
}

// NewEnvironmentConfigProvider creates a new environment config provider
func NewEnvironmentConfigProvider(prefix string) *EnvironmentConfigProvider {
	return &EnvironmentConfigProvider{prefix: prefix}
}

// Load loads configuration from environment variables
func (ecp *EnvironmentConfigProvider) Load() (map[string]interface{}, error) {
	config := make(map[string]interface{})

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key, value := pair[0], pair[1]

		// Filter by prefix if provided
		if ecp.prefix != "" && !strings.HasPrefix(key, ecp.prefix+"_") {
			continue
		}

		// Remove prefix
		if ecp.prefix != "" {
			key = strings.TrimPrefix(key, ecp.prefix+"_")
		}

		// Convert to nested structure
		ecp.setNestedValue(config, strings.ToLower(key), ecp.parseValue(value))
	}

	return config, nil
}

// Watch watches for environment variable changes
func (ecp *EnvironmentConfigProvider) Watch(callback func(map[string]interface{})) error {
	// Environment variables typically don't change during runtime
	// This is a placeholder for potential OS-specific implementations
	return nil
}

// GetName returns the provider name
func (ecp *EnvironmentConfigProvider) GetName() string {
	return "environment"
}

// setNestedValue sets a nested value in the config map
func (ecp *EnvironmentConfigProvider) setNestedValue(config map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, "_")
	current := config

	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
		} else {
			if _, exists := current[k]; !exists {
				current[k] = make(map[string]interface{})
			}
			if next, ok := current[k].(map[string]interface{}); ok {
				current = next
			}
		}
	}
}

// parseValue attempts to parse string values to appropriate types
func (ecp *EnvironmentConfigProvider) parseValue(value string) interface{} {
	// Try boolean
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	// Try integer
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}

	// Try float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Try duration
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}

	// Return as string
	return value
}

// FileConfigProvider loads configuration from files
type FileConfigProvider struct {
	filePath string
	format   string
}

// NewFileConfigProvider creates a new file config provider
func NewFileConfigProvider(filePath string) *FileConfigProvider {
	ext := strings.ToLower(filepath.Ext(filePath))
	format := "json"

	switch ext {
	case ".yaml", ".yml":
		format = "yaml"
	case ".json":
		format = "json"
	}

	return &FileConfigProvider{
		filePath: filePath,
		format:   format,
	}
}

// Load loads configuration from file
func (fcp *FileConfigProvider) Load() (map[string]interface{}, error) {
	data, err := os.ReadFile(fcp.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}

	switch fcp.format {
	case "json":
		err = json.Unmarshal(data, &config)
	case "yaml":
		err = yaml.Unmarshal(data, &config)
	default:
		err = fmt.Errorf("unsupported config format: %s", fcp.format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// Watch watches for file changes
func (fcp *FileConfigProvider) Watch(callback func(map[string]interface{})) error {
	// File watching would require a file system watcher
	// This is a placeholder for potential implementations
	return nil
}

// GetName returns the provider name
func (fcp *FileConfigProvider) GetName() string {
	return fmt.Sprintf("file(%s)", fcp.filePath)
}

// ConfigService manages application configuration
type ConfigService struct {
	providers []ConfigProvider
	config    map[string]interface{}
	logger    *logrus.Logger
	mutex     sync.RWMutex
	watchers  []func(map[string]interface{})
}

// NewConfigService creates a new configuration service
func NewConfigService(logger *logrus.Logger) *ConfigService {
	return &ConfigService{
		providers: make([]ConfigProvider, 0),
		config:    make(map[string]interface{}),
		logger:    logger,
		watchers:  make([]func(map[string]interface{}), 0),
	}
}

// AddProvider adds a configuration provider
func (cs *ConfigService) AddProvider(provider ConfigProvider) {
	cs.providers = append(cs.providers, provider)
	cs.logger.Infof("Added config provider: %s", provider.GetName())
}

// Load loads configuration from all providers
func (cs *ConfigService) Load() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	mergedConfig := make(map[string]interface{})

	for _, provider := range cs.providers {
		config, err := provider.Load()
		if err != nil {
			cs.logger.WithError(err).Warnf("Failed to load config from provider: %s", provider.GetName())
			continue
		}

		// Merge configurations (later providers override earlier ones)
		cs.mergeConfig(mergedConfig, config)
		// Loaded config from provider: %s
	}

	cs.config = mergedConfig
	cs.notifyWatchers()

	return nil
}

// mergeConfig merges two configuration maps
func (cs *ConfigService) mergeConfig(dest, src map[string]interface{}) {
	for key, value := range src {
		if destValue, exists := dest[key]; exists {
			if destMap, ok := destValue.(map[string]interface{}); ok {
				if srcMap, ok := value.(map[string]interface{}); ok {
					cs.mergeConfig(destMap, srcMap)
					continue
				}
			}
		}
		dest[key] = value
	}
}

// Get retrieves a configuration value
func (cs *ConfigService) Get(key string) interface{} {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	return cs.getNestedValue(cs.config, key)
}

// GetString retrieves a string configuration value
func (cs *ConfigService) GetString(key string, defaultValue ...string) string {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", value)
}

// GetInt retrieves an integer configuration value
func (cs *ConfigService) GetInt(key string, defaultValue ...int) int {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetFloat retrieves a float configuration value
func (cs *ConfigService) GetFloat(key string, defaultValue ...float64) float64 {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0.0
}

// GetBool retrieves a boolean configuration value
func (cs *ConfigService) GetBool(key string, defaultValue ...bool) bool {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetDuration retrieves a duration configuration value
func (cs *ConfigService) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	value := cs.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case time.Duration:
		return v
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	case int64:
		return time.Duration(v)
	case int:
		return time.Duration(v)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetStruct retrieves and unmarshals a configuration section into a struct
func (cs *ConfigService) GetStruct(key string, dest interface{}) error {
	value := cs.Get(key)
	if value == nil {
		return errors.New("configuration key not found")
	}

	// Convert to JSON and back to properly unmarshal
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal config value: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Set sets a configuration value
func (cs *ConfigService) Set(key string, value interface{}) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.setNestedValue(cs.config, key, value)
	cs.notifyWatchers()
}

// getNestedValue retrieves a nested value from the config map
func (cs *ConfigService) getNestedValue(config map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	current := config

	for i, k := range keys {
		if i == len(keys)-1 {
			return current[k]
		}

		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// setNestedValue sets a nested value in the config map
func (cs *ConfigService) setNestedValue(config map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := config

	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
		} else {
			if _, exists := current[k]; !exists {
				current[k] = make(map[string]interface{})
			}
			if next, ok := current[k].(map[string]interface{}); ok {
				current = next
			}
		}
	}
}

// OnChange registers a callback for configuration changes
func (cs *ConfigService) OnChange(callback func(map[string]interface{})) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.watchers = append(cs.watchers, callback)
}

// notifyWatchers notifies all watchers of configuration changes
func (cs *ConfigService) notifyWatchers() {
	configCopy := make(map[string]interface{})
	for k, v := range cs.config {
		configCopy[k] = v
	}

	for _, watcher := range cs.watchers {
		go watcher(configCopy)
	}
}

// GetAll returns a copy of the entire configuration
func (cs *ConfigService) GetAll() map[string]interface{} {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	result := make(map[string]interface{})
	for k, v := range cs.config {
		result[k] = v
	}
	return result
}

// ConfigValidator interface for configuration validation
type ConfigValidator interface {
	Validate(config map[string]interface{}) error
}

// StructConfigValidator validates configuration against a struct
type StructConfigValidator struct {
	structType reflect.Type
}

// NewStructConfigValidator creates a new struct config validator
func NewStructConfigValidator(example interface{}) *StructConfigValidator {
	return &StructConfigValidator{
		structType: reflect.TypeOf(example),
	}
}

// Validate validates configuration against the struct
func (scv *StructConfigValidator) Validate(config map[string]interface{}) error {
	// This is a simplified validation
	// A full implementation would validate all required fields and types
	return nil
}

// ConfigModule provides configuration dependency injection
type ConfigModule struct {
	service *ConfigService
}

// NewConfigModule creates a new config module
func NewConfigModule(service *ConfigService) *ConfigModule {
	return &ConfigModule{service: service}
}

// Configure configures the module with providers
func (cm *ConfigModule) Configure(providers ...ConfigProvider) error {
	for _, provider := range providers {
		cm.service.AddProvider(provider)
	}

	return cm.service.Load()
}

// GetService returns the config service
func (cm *ConfigModule) GetService() *ConfigService {
	return cm.service
}

// Configuration decorators

// ConfigProperty decorator for automatic config injection
type ConfigProperty struct {
	Key          string
	DefaultValue interface{}
	Required     bool
}

// InjectConfig injects configuration values into a struct
func InjectConfig(configService *ConfigService, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	structValue := targetValue.Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Check for config tag
		configTag := fieldType.Tag.Get("config")
		if configTag == "" {
			continue
		}

		// Parse tag options
		tagParts := strings.Split(configTag, ",")
		key := tagParts[0]

		var defaultValue interface{}
		required := false

		for _, part := range tagParts[1:] {
			if strings.HasPrefix(part, "default=") {
				defaultValue = strings.TrimPrefix(part, "default=")
			} else if part == "required" {
				required = true
			}
		}

		// Get configuration value
		configValue := configService.Get(key)
		if configValue == nil {
			if required {
				return fmt.Errorf("required configuration key not found: %s", key)
			}
			if defaultValue != nil {
				configValue = defaultValue
			} else {
				continue
			}
		}

		// Set field value
		if err := setFieldValue(field, configValue); err != nil {
			return fmt.Errorf("failed to set config field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue sets a struct field value with type conversion
func setFieldValue(field reflect.Value, value interface{}) error {
	if !field.CanSet() {
		return errors.New("field cannot be set")
	}

	valueReflect := reflect.ValueOf(value)
	fieldType := field.Type()

	// Direct assignment if types match
	if valueReflect.Type().AssignableTo(fieldType) {
		field.Set(valueReflect)
		return nil
	}

	// Type conversion
	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, ok := convertToInt(value); ok {
			field.SetInt(i)
		} else {
			return fmt.Errorf("cannot convert %v to int", value)
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := convertToFloat(value); ok {
			field.SetFloat(f)
		} else {
			return fmt.Errorf("cannot convert %v to float", value)
		}
	case reflect.Bool:
		if b, ok := convertToBool(value); ok {
			field.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %v to bool", value)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return nil
}

// Helper functions for type conversion
func convertToInt(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

func convertToFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func convertToBool(value interface{}) (bool, bool) {
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b, true
		}
	}
	return false, false
}

// Environment-specific configuration loading
func LoadEnvironmentConfig(env string, logger *logrus.Logger) (*ConfigService, error) {
	configService := NewConfigService(logger)

	// Add environment variables provider
	configService.AddProvider(NewEnvironmentConfigProvider("GONEST"))

	// Add base configuration file
	baseConfigFile := "config.yaml"
	if _, err := os.Stat(baseConfigFile); err == nil {
		configService.AddProvider(NewFileConfigProvider(baseConfigFile))
	}

	// Add environment-specific configuration file
	envConfigFile := fmt.Sprintf("config.%s.yaml", env)
	if _, err := os.Stat(envConfigFile); err == nil {
		configService.AddProvider(NewFileConfigProvider(envConfigFile))
	}

	// Load all configurations
	if err := configService.Load(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return configService, nil
}
