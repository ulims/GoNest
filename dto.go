package gonest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// DTO represents a Data Transfer Object
type DTO struct {
	Type   reflect.Type
	Value  reflect.Value
	Fields map[string]*DTOField
	Tags   map[string]string
}

// DTOField represents a field in a DTO
type DTOField struct {
	Name     string
	Type     reflect.Type
	Value    reflect.Value
	Tags     map[string]string
	Required bool
	Min      int
	Max      int
	Pattern  string
}

// DTOBuilder provides a fluent interface for building DTOs
type DTOBuilder struct {
	dto *DTO
}

// NewDTO creates a new DTO
func NewDTO() *DTOBuilder {
	return &DTOBuilder{
		dto: &DTO{
			Fields: make(map[string]*DTOField),
			Tags:   make(map[string]string),
		},
	}
}

// Field adds a field to the DTO
func (db *DTOBuilder) Field(name string, fieldType reflect.Type, tags map[string]string) *DTOBuilder {
	field := &DTOField{
		Name:  name,
		Type:  fieldType,
		Tags:  tags,
		Value: reflect.New(fieldType).Elem(),
	}

	// Parse validation tags
	if required, exists := tags["required"]; exists && required == "true" {
		field.Required = true
	}

	if min, exists := tags["min"]; exists {
		if minInt, err := parseInt(min); err == nil {
			field.Min = minInt
		}
	}

	if max, exists := tags["max"]; exists {
		if maxInt, err := parseInt(max); err == nil {
			field.Max = maxInt
		}
	}

	if pattern, exists := tags["pattern"]; exists {
		field.Pattern = pattern
	}

	db.dto.Fields[name] = field
	return db
}

// Tag adds a tag to the DTO
func (db *DTOBuilder) Tag(key, value string) *DTOBuilder {
	db.dto.Tags[key] = value
	return db
}

// Build returns the built DTO
func (db *DTOBuilder) Build() *DTO {
	return db.dto
}

// DTOValidator validates DTOs
type DTOValidator struct {
	validator *validator.Validate
}

// NewDTOValidator creates a new DTO validator
func NewDTOValidator() *DTOValidator {
	return &DTOValidator{
		validator: validator.New(),
	}
}

// Validate validates a DTO instance
func (dv *DTOValidator) Validate(dto interface{}) error {
	return dv.validator.Struct(dto)
}

// ValidateField validates a specific field
func (dv *DTOValidator) ValidateField(dto interface{}, field string) error {
	return dv.validator.Var(reflect.ValueOf(dto).FieldByName(field).Interface(), field)
}

// CreateDTO creates a DTO from a struct
func CreateDTO(structType reflect.Type) *DTO {
	dto := &DTO{
		Type:   structType,
		Fields: make(map[string]*DTOField),
		Tags:   make(map[string]string),
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := field.Type

		// Parse struct tags
		tags := parseStructTags(field.Tag)

		dtoField := &DTOField{
			Name:  field.Name,
			Type:  fieldType,
			Tags:  tags,
			Value: reflect.New(fieldType).Elem(),
		}

		// Parse validation tags
		if required, exists := tags["required"]; exists && required == "true" {
			dtoField.Required = true
		}

		if min, exists := tags["min"]; exists {
			if minInt, err := parseInt(min); err == nil {
				dtoField.Min = minInt
			}
		}

		if max, exists := tags["max"]; exists {
			if maxInt, err := parseInt(max); err == nil {
				dtoField.Max = maxInt
			}
		}

		if pattern, exists := tags["pattern"]; exists {
			dtoField.Pattern = pattern
		}

		dto.Fields[field.Name] = dtoField
	}

	return dto
}

// parseStructTags parses struct tags into a map
func parseStructTags(tag reflect.StructTag) map[string]string {
	tags := make(map[string]string)

	// Parse validation tags
	if validateTag := tag.Get("validate"); validateTag != "" {
		parts := strings.Split(validateTag, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				tags[kv[0]] = kv[1]
			} else {
				tags[part] = "true"
			}
		}
	}

	// Parse JSON tags
	if jsonTag := tag.Get("json"); jsonTag != "" {
		tags["json"] = jsonTag
	}

	// Parse binding tags
	if bindTag := tag.Get("bind"); bindTag != "" {
		tags["bind"] = bindTag
	}

	return tags
}

// parseInt parses an integer from a string
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// DTO decorators for validation
type Validate struct {
	Rules string
}

// ValidateDecoratorFunc creates a validation decorator
func ValidateDecoratorFunc(rules string) Validate {
	return Validate{Rules: rules}
}

// IsString decorator for string validation
func IsString() Validate {
	return Validate{Rules: "string"}
}

// IsEmail decorator for email validation
func IsEmail() Validate {
	return Validate{Rules: "email"}
}

// IsURL decorator for URL validation
func IsURL() Validate {
	return Validate{Rules: "url"}
}

// Min decorator for minimum value validation
func Min(value int) Validate {
	return Validate{Rules: fmt.Sprintf("min=%d", value)}
}

// Max decorator for maximum value validation
func Max(value int) Validate {
	return Validate{Rules: fmt.Sprintf("max=%d", value)}
}

// Required decorator for required field validation
func Required() Validate {
	return Validate{Rules: "required"}
}

// Pattern decorator for pattern validation
func Pattern(pattern string) Validate {
	return Validate{Rules: fmt.Sprintf("pattern=%s", pattern)}
}

// Transform decorator for data transformation
type Transform struct {
	Function func(interface{}) interface{}
}

// TransformDecoratorFunc creates a transform decorator
func TransformDecoratorFunc(fn func(interface{}) interface{}) Transform {
	return Transform{Function: fn}
}

// SerializeDecorator represents serialization
type SerializeDecorator struct {
	Format string
}

// SerializeDecoratorFunc creates a serialize decorator
func SerializeDecoratorFunc(format string) SerializeDecorator {
	return SerializeDecorator{Format: format}
}
