package user

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// User represents a user entity in the system
type User struct {
	ID        string    `json:"id" validate:"required"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=8"`
	FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
	Role      string    `json:"role" validate:"required,oneof=user admin moderator"`
	Status    string    `json:"status" validate:"required,oneof=active inactive suspended"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// NewUser creates a new user instance
func NewUser(username, email, password, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:        generateID(),
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email, password, firstName, lastName string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if username already exists
	if s.usernameExists(username) {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if s.emailExists(email) {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := NewUser(username, email, password, firstName, lastName)

	// Store user
	s.users[user.ID] = user

	s.logger.Infof("Created user: %s (%s)", user.Username, user.ID)
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

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id string, firstName, lastName, email, status string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if email != "" {
		// Check if email is already taken by another user
		if email != user.Email && s.emailExists(email) {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}
	if status != "" {
		user.Status = status
	}

	user.UpdatedAt = time.Now()
	s.users[id] = user

	s.logger.Infof("Updated user: %s (%s)", user.Username, user.ID)
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	s.logger.Infof("Deleted user: %s (%s)", user.Username, user.ID)
	return nil
}

// ListUsers retrieves a list of users
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

func (s *UserService) emailExists(email string) bool {
	for _, user := range s.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// User methods
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsActive() bool {
	return u.Status == "active"
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) UpdateProfile(firstName, lastName string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}

func (u *User) ChangeStatus(status string) {
	u.Status = status
	u.UpdatedAt = time.Now()
}

// generateID generates a unique ID (simplified for example)
func generateID() string {
	return time.Now().Format("20060102150405")
}
