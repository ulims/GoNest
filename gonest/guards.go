package gonest

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

// Guard interface for route protection
type Guard interface {
	CanActivate(ctx echo.Context) (bool, error)
}

// GuardFunc is a function type that implements Guard interface
type GuardFunc func(ctx echo.Context) (bool, error)

// CanActivate implements Guard interface for GuardFunc
func (f GuardFunc) CanActivate(ctx echo.Context) (bool, error) {
	return f(ctx)
}

// GuardMetadata contains guard information
type GuardMetadata struct {
	Guard    Guard
	Priority int
	Metadata map[string]interface{}
}

// GuardRegistry manages all guards
type GuardRegistry struct {
	guards map[string]*GuardMetadata
}

// NewGuardRegistry creates a new guard registry
func NewGuardRegistry() *GuardRegistry {
	return &GuardRegistry{
		guards: make(map[string]*GuardMetadata),
	}
}

// Register registers a guard
func (gr *GuardRegistry) Register(name string, guard Guard, priority int) {
	gr.guards[name] = &GuardMetadata{
		Guard:    guard,
		Priority: priority,
		Metadata: make(map[string]interface{}),
	}
}

// Get retrieves a guard by name
func (gr *GuardRegistry) Get(name string) (Guard, bool) {
	if metadata, exists := gr.guards[name]; exists {
		return metadata.Guard, true
	}
	return nil, false
}

// GetAll returns all registered guards
func (gr *GuardRegistry) GetAll() map[string]*GuardMetadata {
	return gr.guards
}

// Guard decorators
type GuardDecorator struct {
	Guards []string
}

// UseGuards decorator for applying guards to routes
func UseGuards(guards ...string) GuardDecorator {
	return GuardDecorator{Guards: guards}
}

// AuthGuard is a basic authentication guard
type AuthGuard struct {
	JWTSecret string
}

// NewAuthGuard creates a new authentication guard
func NewAuthGuard(jwtSecret string) *AuthGuard {
	return &AuthGuard{JWTSecret: jwtSecret}
}

// CanActivate checks if the request is authenticated
func (ag *AuthGuard) CanActivate(ctx echo.Context) (bool, error) {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		return false, echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}

	// Basic token validation (in real app, validate JWT)
	if len(token) < 10 {
		return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return true, nil
}

// RoleGuard is a role-based authorization guard
type RoleGuard struct {
	RequiredRoles []string
}

// NewRoleGuard creates a new role-based guard
func NewRoleGuard(roles ...string) *RoleGuard {
	return &RoleGuard{RequiredRoles: roles}
}

// CanActivate checks if the user has required roles
func (rg *RoleGuard) CanActivate(ctx echo.Context) (bool, error) {
	// Extract user roles from context (set by AuthGuard)
	userRoles, ok := ctx.Get("user_roles").([]string)
	if !ok {
		return false, echo.NewHTTPError(http.StatusForbidden, "User roles not found")
	}

	// Check if user has any of the required roles
	for _, requiredRole := range rg.RequiredRoles {
		for _, userRole := range userRoles {
			if userRole == requiredRole {
				return true, nil
			}
		}
	}

	return false, echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
}

// RateLimitGuard is a rate limiting guard
type RateLimitGuard struct {
	MaxRequests int
	Window      int // in seconds
}

// NewRateLimitGuard creates a new rate limiting guard
func NewRateLimitGuard(maxRequests, window int) *RateLimitGuard {
	return &RateLimitGuard{
		MaxRequests: maxRequests,
		Window:      window,
	}
}

// CanActivate checks if the request is within rate limits
func (rlg *RateLimitGuard) CanActivate(ctx echo.Context) (bool, error) {
	// In a real implementation, this would check against Redis or memory store
	// For now, we'll just return true
	return true, nil
}

// GuardMiddleware creates middleware from guards
func GuardMiddleware(guards ...Guard) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, guard := range guards {
				canActivate, err := guard.CanActivate(c)
				if err != nil {
					return err
				}
				if !canActivate {
					return echo.NewHTTPError(http.StatusForbidden, "Access denied")
				}
			}
			return next(c)
		}
	}
}

// Guard decorator for struct tags
type GuardTag struct {
	Guards []string
}

// ParseGuardTags parses guard tags from struct fields
func ParseGuardTags(tag reflect.StructTag) *GuardTag {
	guardTag := tag.Get("guard")
	if guardTag == "" {
		return nil
	}

	// Parse comma-separated guards
	guards := []string{}
	// Implementation would parse the guard tag
	return &GuardTag{Guards: guards}
}
