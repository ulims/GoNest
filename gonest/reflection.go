package gonest

import (
	"reflect"
	"strings"
)

// Reflection utilities for dependency injection and metadata handling

// GetStructMetadata extracts metadata from a struct
func GetStructMetadata(structType reflect.Type) map[string]interface{} {
	metadata := make(map[string]interface{})

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldMetadata := make(map[string]interface{})

		// Extract tags
		tags := parseStructTags(field.Tag)
		fieldMetadata["tags"] = tags

		// Extract field type
		fieldMetadata["type"] = field.Type.String()
		fieldMetadata["kind"] = field.Type.Kind().String()

		// Check if field is exported
		fieldMetadata["exported"] = field.IsExported()

		// Check if field is embedded
		fieldMetadata["embedded"] = field.Anonymous

		metadata[field.Name] = fieldMetadata
	}

	return metadata
}

// GetMethodMetadata extracts metadata from methods
func GetMethodMetadata(structType reflect.Type) map[string]interface{} {
	metadata := make(map[string]interface{})

	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		methodMetadata := make(map[string]interface{})

		// Extract method signature
		methodMetadata["name"] = method.Name
		methodMetadata["type"] = method.Type.String()

		// Extract parameters
		params := make([]map[string]interface{}, 0)
		for j := 0; j < method.Type.NumIn(); j++ {
			param := method.Type.In(j)
			paramInfo := map[string]interface{}{
				"type": param.String(),
				"kind": param.Kind().String(),
			}
			params = append(params, paramInfo)
		}
		methodMetadata["parameters"] = params

		// Extract return values
		returns := make([]map[string]interface{}, 0)
		for j := 0; j < method.Type.NumOut(); j++ {
			ret := method.Type.Out(j)
			retInfo := map[string]interface{}{
				"type": ret.String(),
				"kind": ret.Kind().String(),
			}
			returns = append(returns, retInfo)
		}
		methodMetadata["returns"] = returns

		metadata[method.Name] = methodMetadata
	}

	return metadata
}

// InjectDependencies injects dependencies into a struct
func InjectDependencies(instance interface{}, serviceRegistry *ServiceRegistry) error {
	value := reflect.ValueOf(instance)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	instanceType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := field.Type()

		// Check if field can be set
		if !field.CanSet() {
			continue
		}

		// Check if field has inject tag
		structField := instanceType.Field(i)
		tags := parseStructTags(structField.Tag)

		if injectName, exists := tags["inject"]; exists {
			// Inject by name
			if dependency, found := serviceRegistry.Get(injectName); found {
				field.Set(reflect.ValueOf(dependency))
			}
		} else {
			// Try to inject by type
			if dependency, found := serviceRegistry.GetByType(fieldType); found {
				field.Set(reflect.ValueOf(dependency))
			}
		}
	}

	return nil
}

// CreateInstance creates a new instance of a type
func CreateInstance(instanceType reflect.Type) interface{} {
	if instanceType.Kind() == reflect.Ptr {
		return reflect.New(instanceType.Elem()).Interface()
	}
	return reflect.New(instanceType).Interface()
}

// IsInjectable checks if a type is injectable
func IsInjectable(instanceType reflect.Type) bool {
	// Check if type has Injectable tag
	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		tags := parseStructTags(field.Tag)
		if _, exists := tags["injectable"]; exists {
			return true
		}
	}

	// Check if type implements Injectable interface
	interfaceType := reflect.TypeOf((*Injectable)(nil)).Elem()
	return instanceType.Implements(interfaceType)
}

// GetInjectableFields returns all injectable fields from a struct
func GetInjectableFields(instanceType reflect.Type) []string {
	injectableFields := make([]string, 0)

	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		tags := parseStructTags(field.Tag)

		if _, exists := tags["inject"]; exists {
			injectableFields = append(injectableFields, field.Name)
		}
	}

	return injectableFields
}

// GetDependencyGraph builds a dependency graph for services
func GetDependencyGraph(services map[string]*Service) map[string][]string {
	graph := make(map[string][]string)

	for name, service := range services {
		dependencies := make([]string, 0)

		if service.Instance != nil {
			instanceType := reflect.TypeOf(service.Instance)
			if instanceType.Kind() == reflect.Ptr {
				instanceType = instanceType.Elem()
			}

			for i := 0; i < instanceType.NumField(); i++ {
				field := instanceType.Field(i)
				tags := parseStructTags(field.Tag)

				if injectName, exists := tags["inject"]; exists {
					dependencies = append(dependencies, injectName)
				}
			}
		}

		graph[name] = dependencies
	}

	return graph
}

// DetectCircularDependencies detects circular dependencies in the service graph
func DetectCircularDependencies(graph map[string][]string) []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	cycle := make([]string, 0)

	for node := range graph {
		if !visited[node] {
			if detectCycleUtil(node, graph, visited, recStack, &cycle) {
				return cycle
			}
		}
	}

	return nil
}

// detectCycleUtil is a utility function for detecting cycles in a graph
func detectCycleUtil(node string, graph map[string][]string, visited, recStack map[string]bool, cycle *[]string) bool {
	visited[node] = true
	recStack[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] {
			if detectCycleUtil(neighbor, graph, visited, recStack, cycle) {
				*cycle = append(*cycle, node)
				return true
			}
		} else if recStack[neighbor] {
			*cycle = append(*cycle, node)
			return true
		}
	}

	recStack[node] = false
	return false
}

// GetMethodDecorators extracts decorators from method tags
func GetMethodDecorators(method reflect.Method) map[string]interface{} {
	decorators := make(map[string]interface{})

	// Extract method name patterns for decorators
	methodName := method.Name

	// Check for HTTP method decorators
	if strings.HasPrefix(methodName, "Get") {
		decorators["method"] = "GET"
		decorators["path"] = "/" + strings.ToLower(strings.TrimPrefix(methodName, "Get"))
	} else if strings.HasPrefix(methodName, "Post") {
		decorators["method"] = "POST"
		decorators["path"] = "/" + strings.ToLower(strings.TrimPrefix(methodName, "Post"))
	} else if strings.HasPrefix(methodName, "Put") {
		decorators["method"] = "PUT"
		decorators["path"] = "/" + strings.ToLower(strings.TrimPrefix(methodName, "Put"))
	} else if strings.HasPrefix(methodName, "Delete") {
		decorators["method"] = "DELETE"
		decorators["path"] = "/" + strings.ToLower(strings.TrimPrefix(methodName, "Delete"))
	} else if strings.HasPrefix(methodName, "Patch") {
		decorators["method"] = "PATCH"
		decorators["path"] = "/" + strings.ToLower(strings.TrimPrefix(methodName, "Patch"))
	}

	return decorators
}

// ValidateStruct validates a struct using reflection
func ValidateStruct(instance interface{}, validator *DTOValidator) error {
	return validator.Validate(instance)
}

// GetStructTags extracts all tags from a struct
func GetStructTags(instanceType reflect.Type) map[string]map[string]string {
	tags := make(map[string]map[string]string)

	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		fieldTags := parseStructTags(field.Tag)
		tags[field.Name] = fieldTags
	}

	return tags
}
