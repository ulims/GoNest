package user

import (
	"errors"
	"sync"
	"time"
	"github.com/sirupsen/logrus"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserService handles user business logic
type UserService struct {
	users  map[string]*User
	logger *logrus.Logger
	mutex  sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService(logger *logrus.Logger) *UserService {
	return &UserService{
		users:  make(map[string]*User),
		logger: logger,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Check if username already exists
	if s.usernameExists(username) {
		return nil, errors.New("username already exists")
	}
	
	// Create new user
	user := &User{
		ID:        time.Now().Format("20060102150405"),
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
	
	// Store user
	s.users[user.ID] = user
	
	s.logger.Infof("Created user: %%s (%%s)", user.Username, user.ID)
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	
	return users, nil
}

// Helper methods
func (s *UserService) usernameExists(username string) bool {
	for _, user := range s.users {
		if user.Username == username {
			return true
		}
	}
	return false
}
