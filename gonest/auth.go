package gonest

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey               string
	PublicKey               *rsa.PublicKey
	PrivateKey              *rsa.PrivateKey
	SigningMethod           jwt.SigningMethod
	TokenLookup             string
	TokenLookupFuncs        []func(echo.Context) (string, error)
	AuthScheme              string
	ContextKey              string
	Claims                  jwt.Claims
	TokenGenerator          func(user interface{}) (string, error)
	ParseTokenFunc          func(token string) (jwt.Claims, error)
	Skipper                 func(echo.Context) bool
	BeforeFunc              func(echo.Context)
	SuccessHandler          func(echo.Context)
	ErrorHandler            func(echo.Context, error) error
	ErrorHandlerWithContext func(error, echo.Context) error
	ParseTokenFilter        func(string, echo.Context) (string, error)
	TokenExpiry             time.Duration
	RefreshExpiry           time.Duration
}

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SigningMethod: jwt.SigningMethodHS256,
		ContextKey:    "user",
		TokenLookup:   "header:Authorization,query:token,cookie:token",
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
		TokenExpiry:   24 * time.Hour,
		RefreshExpiry: 7 * 24 * time.Hour,
	}
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID   string                 `json:"sub"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Roles    []string               `json:"roles"`
	Metadata map[string]interface{} `json:"metadata"`
	jwt.RegisteredClaims
}

// AuthUser represents an authenticated user
type AuthUser struct {
	ID       string                 `json:"id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Roles    []string               `json:"roles"`
	Metadata map[string]interface{} `json:"metadata"`
}

// AuthService provides authentication functionality
type AuthService struct {
	config *JWTConfig
	logger *logrus.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(config *JWTConfig, logger *logrus.Logger) *AuthService {
	if config == nil {
		config = DefaultJWTConfig()
	}

	return &AuthService{
		config: config,
		logger: logger,
	}
}

// GenerateToken generates a JWT token for a user
func (as *AuthService) GenerateToken(user *AuthUser) (string, error) {
	if as.config.TokenGenerator != nil {
		return as.config.TokenGenerator(user)
	}

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		Metadata: user.Metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(as.config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gonest",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(as.config.SigningMethod, claims)

	if as.config.SecretKey != "" {
		return token.SignedString([]byte(as.config.SecretKey))
	} else if as.config.PrivateKey != nil {
		return token.SignedString(as.config.PrivateKey)
	}

	return "", errors.New("no signing key configured")
}

// GenerateRefreshToken generates a refresh token
func (as *AuthService) GenerateRefreshToken(user *AuthUser) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(as.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gonest-refresh",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(as.config.SigningMethod, claims)

	if as.config.SecretKey != "" {
		return token.SignedString([]byte(as.config.SecretKey))
	} else if as.config.PrivateKey != nil {
		return token.SignedString(as.config.PrivateKey)
	}

	return "", errors.New("no signing key configured")
}

// ValidateToken validates a JWT token
func (as *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	if as.config.ParseTokenFunc != nil {
		claims, err := as.config.ParseTokenFunc(tokenString)
		if err != nil {
			return nil, err
		}
		if jwtClaims, ok := claims.(*JWTClaims); ok {
			return jwtClaims, nil
		}
		return nil, errors.New("invalid token claims")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if as.config.SecretKey != "" {
			return []byte(as.config.SecretKey), nil
		} else if as.config.PublicKey != nil {
			return as.config.PublicKey, nil
		}
		return nil, errors.New("no verification key configured")
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractToken extracts token from request
func (as *AuthService) ExtractToken(c echo.Context) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	if auth != "" {
		l := len(as.config.AuthScheme)
		if len(auth) > l+1 && strings.EqualFold(auth[:l], as.config.AuthScheme) {
			return auth[l+1:], nil
		}
	}

	// Try query parameter
	token := c.QueryParam("token")
	if token != "" {
		return token, nil
	}

	// Try cookie
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", errors.New("token not found")
}

// JWTMiddleware returns JWT middleware
func (as *AuthService) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if as.config.Skipper != nil && as.config.Skipper(c) {
				return next(c)
			}

			if as.config.BeforeFunc != nil {
				as.config.BeforeFunc(c)
			}

			tokenString, err := as.ExtractToken(c)
			if err != nil {
				if as.config.ErrorHandler != nil {
					return as.config.ErrorHandler(c, err)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			claims, err := as.ValidateToken(tokenString)
			if err != nil {
				if as.config.ErrorHandler != nil {
					return as.config.ErrorHandler(c, err)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			user := &AuthUser{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    claims.Email,
				Roles:    claims.Roles,
				Metadata: claims.Metadata,
			}

			c.Set(as.config.ContextKey, user)

			if as.config.SuccessHandler != nil {
				as.config.SuccessHandler(c)
			}

			return next(c)
		}
	}
}

// RequireRoles middleware that requires specific roles
func (as *AuthService) RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get(as.config.ContextKey).(*AuthUser)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
			}

			if len(roles) == 0 {
				return next(c)
			}

			userRoles := make(map[string]bool)
			for _, role := range user.Roles {
				userRoles[role] = true
			}

			for _, requiredRole := range roles {
				if userRoles[requiredRole] {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

// OptionalAuth middleware that optionally authenticates
func (as *AuthService) OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := as.ExtractToken(c)
			if err != nil {
				// No token found, continue without authentication
				return next(c)
			}

			claims, err := as.ValidateToken(tokenString)
			if err != nil {
				// Invalid token, continue without authentication
				return next(c)
			}

			user := &AuthUser{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    claims.Email,
				Roles:    claims.Roles,
				Metadata: claims.Metadata,
			}

			c.Set(as.config.ContextKey, user)
			return next(c)
		}
	}
}

// GetCurrentUser extracts current user from context
func GetCurrentUser(c echo.Context) (*AuthUser, error) {
	user, ok := c.Get("user").(*AuthUser)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

// AuthStrategy interface for different authentication strategies
type AuthStrategy interface {
	Authenticate(ctx context.Context, credentials interface{}) (*AuthUser, error)
	GetName() string
}

// LocalStrategy implements username/password authentication
type LocalStrategy struct {
	validator func(username, password string) (*AuthUser, error)
	name      string
}

// NewLocalStrategy creates a new local strategy
func NewLocalStrategy(validator func(username, password string) (*AuthUser, error)) *LocalStrategy {
	return &LocalStrategy{
		validator: validator,
		name:      "local",
	}
}

// Authenticate authenticates using username and password
func (ls *LocalStrategy) Authenticate(ctx context.Context, credentials interface{}) (*AuthUser, error) {
	creds, ok := credentials.(map[string]string)
	if !ok {
		return nil, errors.New("invalid credentials format")
	}

	username, ok := creds["username"]
	if !ok {
		return nil, errors.New("username required")
	}

	password, ok := creds["password"]
	if !ok {
		return nil, errors.New("password required")
	}

	return ls.validator(username, password)
}

// GetName returns strategy name
func (ls *LocalStrategy) GetName() string {
	return ls.name
}

// PassportService manages authentication strategies
type PassportService struct {
	strategies map[string]AuthStrategy
	logger     *logrus.Logger
}

// NewPassportService creates a new passport service
func NewPassportService(logger *logrus.Logger) *PassportService {
	return &PassportService{
		strategies: make(map[string]AuthStrategy),
		logger:     logger,
	}
}

// Use registers an authentication strategy
func (ps *PassportService) Use(strategy AuthStrategy) {
	ps.strategies[strategy.GetName()] = strategy
	ps.logger.Infof("Registered auth strategy: %s", strategy.GetName())
}

// Authenticate authenticates using a specific strategy
func (ps *PassportService) Authenticate(ctx context.Context, strategyName string, credentials interface{}) (*AuthUser, error) {
	strategy, exists := ps.strategies[strategyName]
	if !exists {
		return nil, fmt.Errorf("strategy '%s' not found", strategyName)
	}

	return strategy.Authenticate(ctx, credentials)
}

// AuthController provides authentication endpoints
type AuthController struct {
	authService     *AuthService
	passportService *PassportService
	logger          *logrus.Logger
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *AuthService, passportService *PassportService, logger *logrus.Logger) *AuthController {
	return &AuthController{
		authService:     authService,
		passportService: passportService,
		logger:          logger,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Strategy string `json:"strategy,omitempty"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *AuthUser `json:"user"`
}

// Login handles user login
func (ac *AuthController) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return BadRequestException("Invalid request format")
	}

	if err := ValidateStruct(&req, nil); err != nil {
		return BadRequestException(fmt.Sprintf("Validation failed: %v", err))
	}

	strategy := req.Strategy
	if strategy == "" {
		strategy = "local"
	}

	credentials := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}

	user, err := ac.passportService.Authenticate(c.Request().Context(), strategy, credentials)
	if err != nil {
		ac.logger.WithError(err).Warn("Authentication failed")
		return UnauthorizedException("Invalid credentials")
	}

	accessToken, err := ac.authService.GenerateToken(user)
	if err != nil {
		ac.logger.WithError(err).Error("Failed to generate access token")
		return InternalServerErrorException("Failed to generate token")
	}

	refreshToken, err := ac.authService.GenerateRefreshToken(user)
	if err != nil {
		ac.logger.WithError(err).Error("Failed to generate refresh token")
		return InternalServerErrorException("Failed to generate refresh token")
	}

	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(ac.authService.config.TokenExpiry.Seconds()),
		User:         user,
	}

	return c.JSON(http.StatusOK, response)
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshToken handles token refresh
func (ac *AuthController) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return BadRequestException("Invalid request format")
	}

	if err := ValidateStruct(&req, nil); err != nil {
		return BadRequestException(fmt.Sprintf("Validation failed: %v", err))
	}

	claims, err := ac.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		return UnauthorizedException("Invalid refresh token")
	}

	// Verify it's a refresh token
	if claims.RegisteredClaims.Issuer != "gonest-refresh" {
		return UnauthorizedException("Invalid refresh token")
	}

	user := &AuthUser{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
		Metadata: claims.Metadata,
	}

	accessToken, err := ac.authService.GenerateToken(user)
	if err != nil {
		ac.logger.WithError(err).Error("Failed to generate access token")
		return InternalServerErrorException("Failed to generate token")
	}

	response := &LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(ac.authService.config.TokenExpiry.Seconds()),
		User:        user,
	}

	return c.JSON(http.StatusOK, response)
}

// GetProfile returns current user profile
func (ac *AuthController) GetProfile(c echo.Context) error {
	user, err := GetCurrentUser(c)
	if err != nil {
		return UnauthorizedException("User not authenticated")
	}

	return c.JSON(http.StatusOK, user)
}
