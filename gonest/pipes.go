package gonest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Pipe interface for data transformation and validation
type Pipe interface {
	Transform(value interface{}) (interface{}, error)
}

// PipeFunc is a function type that implements Pipe interface
type PipeFunc func(value interface{}) (interface{}, error)

// Transform implements Pipe interface for PipeFunc
func (f PipeFunc) Transform(value interface{}) (interface{}, error) {
	return f(value)
}

// PipeMetadata contains pipe information
type PipeMetadata struct {
	Pipe     Pipe
	Priority int
	Metadata map[string]interface{}
}

// PipeRegistry manages all pipes
type PipeRegistry struct {
	pipes map[string]*PipeMetadata
}

// NewPipeRegistry creates a new pipe registry
func NewPipeRegistry() *PipeRegistry {
	return &PipeRegistry{
		pipes: make(map[string]*PipeMetadata),
	}
}

// Register registers a pipe
func (pr *PipeRegistry) Register(name string, pipe Pipe, priority int) {
	pr.pipes[name] = &PipeMetadata{
		Pipe:     pipe,
		Priority: priority,
		Metadata: make(map[string]interface{}),
	}
}

// Get retrieves a pipe by name
func (pr *PipeRegistry) Get(name string) (Pipe, bool) {
	if metadata, exists := pr.pipes[name]; exists {
		return metadata.Pipe, true
	}
	return nil, false
}

// GetAll returns all registered pipes
func (pr *PipeRegistry) GetAll() map[string]*PipeMetadata {
	return pr.pipes
}

// Pipe decorators
type PipeDecorator struct {
	Pipes []string
}

// UsePipes decorator for applying pipes to parameters
func UsePipes(pipes ...string) PipeDecorator {
	return PipeDecorator{Pipes: pipes}
}

// ValidationPipe validates data using struct tags
type ValidationPipe struct {
	validator *validator.Validate
}

// NewValidationPipe creates a new validation pipe
func NewValidationPipe() *ValidationPipe {
	return &ValidationPipe{
		validator: validator.New(),
	}
}

// Transform validates the input data
func (vp *ValidationPipe) Transform(value interface{}) (interface{}, error) {
	if err := vp.validator.Struct(value); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}
	return value, nil
}

// ParseIntPipe parses string to integer
type ParseIntPipe struct{}

// NewParseIntPipe creates a new integer parsing pipe
func NewParseIntPipe() *ParseIntPipe {
	return &ParseIntPipe{}
}

// Transform parses string to integer
func (pip *ParseIntPipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return nil, fmt.Errorf("cannot parse %v to int", value)
	}
}

// ParseFloatPipe parses string to float
type ParseFloatPipe struct{}

// NewParseFloatPipe creates a new float parsing pipe
func NewParseFloatPipe() *ParseFloatPipe {
	return &ParseFloatPipe{}
}

// Transform parses string to float
func (pfp *ParseFloatPipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return nil, fmt.Errorf("cannot parse %v to float", value)
	}
}

// ParseBoolPipe parses string to boolean
type ParseBoolPipe struct{}

// NewParseBoolPipe creates a new boolean parsing pipe
func NewParseBoolPipe() *ParseBoolPipe {
	return &ParseBoolPipe{}
}

// Transform parses string to boolean
func (pbp *ParseBoolPipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseBool(strings.ToLower(v))
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	default:
		return nil, fmt.Errorf("cannot parse %v to bool", value)
	}
}

// TrimPipe trims whitespace from strings
type TrimPipe struct{}

// NewTrimPipe creates a new trim pipe
func NewTrimPipe() *TrimPipe {
	return &TrimPipe{}
}

// Transform trims whitespace from strings
func (tp *TrimPipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v), nil
	default:
		return value, nil
	}
}

// LowercasePipe converts strings to lowercase
type LowercasePipe struct{}

// NewLowercasePipe creates a new lowercase pipe
func NewLowercasePipe() *LowercasePipe {
	return &LowercasePipe{}
}

// Transform converts strings to lowercase
func (lp *LowercasePipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strings.ToLower(v), nil
	default:
		return value, nil
	}
}

// UppercasePipe converts strings to uppercase
type UppercasePipe struct{}

// NewUppercasePipe creates a new uppercase pipe
func NewUppercasePipe() *UppercasePipe {
	return &UppercasePipe{}
}

// Transform converts strings to uppercase
func (up *UppercasePipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return strings.ToUpper(v), nil
	default:
		return value, nil
	}
}

// DefaultValuePipe provides default values
type DefaultValuePipe struct {
	DefaultValue interface{}
}

// NewDefaultValuePipe creates a new default value pipe
func NewDefaultValuePipe(defaultValue interface{}) *DefaultValuePipe {
	return &DefaultValuePipe{DefaultValue: defaultValue}
}

// Transform provides default value if input is empty
func (dvp *DefaultValuePipe) Transform(value interface{}) (interface{}, error) {
	if value == nil || value == "" {
		return dvp.DefaultValue, nil
	}
	return value, nil
}

// JSONPipe parses JSON strings
type JSONPipe struct {
	TargetType reflect.Type
}

// NewJSONPipe creates a new JSON parsing pipe
func NewJSONPipe(targetType reflect.Type) *JSONPipe {
	return &JSONPipe{TargetType: targetType}
}

// Transform parses JSON strings to structs
func (jp *JSONPipe) Transform(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		target := reflect.New(jp.TargetType).Interface()
		if err := json.Unmarshal([]byte(v), target); err != nil {
			return nil, fmt.Errorf("JSON parsing failed: %v", err)
		}
		return target, nil
	default:
		return value, nil
	}
}

// CustomPipe allows custom transformation logic
type CustomPipe struct {
	TransformFunc func(interface{}) (interface{}, error)
}

// NewCustomPipe creates a new custom pipe
func NewCustomPipe(transformFunc func(interface{}) (interface{}, error)) *CustomPipe {
	return &CustomPipe{TransformFunc: transformFunc}
}

// Transform applies custom transformation
func (cp *CustomPipe) Transform(value interface{}) (interface{}, error) {
	return cp.TransformFunc(value)
}

// PipeChain chains multiple pipes together
type PipeChain struct {
	pipes []Pipe
}

// NewPipeChain creates a new pipe chain
func NewPipeChain(pipes ...Pipe) *PipeChain {
	return &PipeChain{pipes: pipes}
}

// Transform applies all pipes in sequence
func (pc *PipeChain) Transform(value interface{}) (interface{}, error) {
	result := value
	for _, pipe := range pc.pipes {
		transformed, err := pipe.Transform(result)
		if err != nil {
			return nil, err
		}
		result = transformed
	}
	return result, nil
}

// PipeMiddleware creates middleware from pipes
func PipeMiddleware(pipes ...Pipe) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Apply pipes to request parameters
			// Implementation would apply pipes to specific parameters
			// based on decorators or configuration
			_ = pipes // Suppress unused variable warning
			return next(c)
		}
	}
}

// ApplyPipes applies pipes to a value
func ApplyPipes(value interface{}, pipes ...Pipe) (interface{}, error) {
	chain := NewPipeChain(pipes...)
	return chain.Transform(value)
}

// Built-in pipe instances
var (
	ParseIntPipeInstance   = NewParseIntPipe()
	ParseFloatPipeInstance = NewParseFloatPipe()
	ParseBoolPipeInstance  = NewParseBoolPipe()
	TrimPipeInstance       = NewTrimPipe()
	LowercasePipeInstance  = NewLowercasePipe()
	UppercasePipeInstance  = NewUppercasePipe()
	ValidationPipeInstance = NewValidationPipe()
)
