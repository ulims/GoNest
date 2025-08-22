package gonest

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// DatabaseInterface interface for database operations
type DatabaseInterface interface {
	Connect() error
	Disconnect() error
	GetDB() *sql.DB
	IsConnected() bool
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
	MaxOpen  int
	MaxIdle  int
	Timeout  time.Duration
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "password",
		Database: "gonest",
		SSLMode:  "disable",
		MaxOpen:  25,
		MaxIdle:  5,
		Timeout:  30 * time.Second,
	}
}

// RepositoryInterface interface for data access
type RepositoryInterface interface {
	Create(model interface{}) error
	FindByID(id interface{}, model interface{}) error
	FindAll(models interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error
	Where(query string, args ...interface{}) RepositoryInterface
	Order(order string) RepositoryInterface
	Limit(limit int) RepositoryInterface
	Offset(offset int) RepositoryInterface
}

// BaseRepository provides basic repository functionality
type BaseRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB, logger *logrus.Logger) *BaseRepository {
	return &BaseRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new record
func (r *BaseRepository) Create(model interface{}) error {
	// Implementation would use reflection to generate SQL
	r.logger.Infof("Creating record: %T", model)
	return nil
}

// FindByID finds a record by ID
func (r *BaseRepository) FindByID(id interface{}, model interface{}) error {
	// Implementation would use reflection to generate SQL
	r.logger.Infof("Finding record by ID: %v", id)
	return nil
}

// FindAll finds all records
func (r *BaseRepository) FindAll(models interface{}) error {
	// Implementation would use reflection to generate SQL
	r.logger.Info("Finding all records")
	return nil
}

// Update updates a record
func (r *BaseRepository) Update(model interface{}) error {
	// Implementation would use reflection to generate SQL
	r.logger.Infof("Updating record: %T", model)
	return nil
}

// Delete deletes a record
func (r *BaseRepository) Delete(model interface{}) error {
	// Implementation would use reflection to generate SQL
	r.logger.Infof("Deleting record: %T", model)
	return nil
}

// Where adds a WHERE clause
func (r *BaseRepository) Where(query string, args ...interface{}) RepositoryInterface {
	// Implementation would build query
	return r
}

// Order adds an ORDER BY clause
func (r *BaseRepository) Order(order string) RepositoryInterface {
	// Implementation would build query
	return r
}

// Limit adds a LIMIT clause
func (r *BaseRepository) Limit(limit int) RepositoryInterface {
	// Implementation would build query
	return r
}

// Offset adds an OFFSET clause
func (r *BaseRepository) Offset(offset int) RepositoryInterface {
	// Implementation would build query
	return r
}

// Model interface for database models
type Model interface {
	TableName() string
	GetID() interface{}
	SetID(id interface{})
	GetCreatedAt() time.Time
	SetCreatedAt(t time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

// BaseModel provides basic model functionality
type BaseModel struct {
	ID        interface{} `json:"id" db:"id"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name
func (bm *BaseModel) TableName() string {
	return "base_models"
}

// GetID returns the ID
func (bm *BaseModel) GetID() interface{} {
	return bm.ID
}

// SetID sets the ID
func (bm *BaseModel) SetID(id interface{}) {
	bm.ID = id
}

// GetCreatedAt returns the created at time
func (bm *BaseModel) GetCreatedAt() time.Time {
	return bm.CreatedAt
}

// SetCreatedAt sets the created at time
func (bm *BaseModel) SetCreatedAt(t time.Time) {
	bm.CreatedAt = t
}

// GetUpdatedAt returns the updated at time
func (bm *BaseModel) GetUpdatedAt() time.Time {
	return bm.UpdatedAt
}

// SetUpdatedAt sets the updated at time
func (bm *BaseModel) SetUpdatedAt(t time.Time) {
	bm.UpdatedAt = t
}

// DatabaseService provides database functionality
type DatabaseService struct {
	config *DatabaseConfig
	db     *sql.DB
	logger *logrus.Logger
}

// NewDatabaseService creates a new database service
func NewDatabaseService(config *DatabaseConfig, logger *logrus.Logger) *DatabaseService {
	return &DatabaseService{
		config: config,
		logger: logger,
	}
}

// Connect connects to the database
func (ds *DatabaseService) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ds.config.Host, ds.config.Port, ds.config.Username, ds.config.Password,
		ds.config.Database, ds.config.SSLMode)

	db, err := sql.Open(ds.config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(ds.config.MaxOpen)
	db.SetMaxIdleConns(ds.config.MaxIdle)
	db.SetConnMaxLifetime(ds.config.Timeout)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	ds.db = db
	ds.logger.Info("Database connected successfully")
	return nil
}

// Disconnect disconnects from the database
func (ds *DatabaseService) Disconnect() error {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %v", err)
		}
		ds.logger.Info("Database disconnected successfully")
	}
	return nil
}

// GetDB returns the database connection
func (ds *DatabaseService) GetDB() *sql.DB {
	return ds.db
}

// IsConnected checks if the database is connected
func (ds *DatabaseService) IsConnected() bool {
	if ds.db == nil {
		return false
	}
	return ds.db.Ping() == nil
}

// Repository decorators
type RepositoryDecorator struct {
	Model interface{}
}

// Repository decorator for marking repositories
func Repository(model interface{}) RepositoryDecorator {
	return RepositoryDecorator{Model: model}
}

// Database decorators
type DatabaseDecorator struct {
	Config *DatabaseConfig
}

// Database decorator for database configuration
func Database(config *DatabaseConfig) DatabaseDecorator {
	return DatabaseDecorator{Config: config}
}

// Migration interface for database migrations
type Migration interface {
	Up() error
	Down() error
	Version() int
	Description() string
}

// MigrationService manages database migrations
type MigrationService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *sql.DB, logger *logrus.Logger) *MigrationService {
	return &MigrationService{
		db:     db,
		logger: logger,
	}
}

// RunMigrations runs all pending migrations
func (ms *MigrationService) RunMigrations(migrations []Migration) error {
	// Create migrations table if it doesn't exist
	if err := ms.createMigrationsTable(); err != nil {
		return err
	}

	for _, migration := range migrations {
		if err := ms.runMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

// createMigrationsTable creates the migrations table
func (ms *MigrationService) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := ms.db.Exec(query)
	return err
}

// runMigration runs a single migration
func (ms *MigrationService) runMigration(migration Migration) error {
	// Check if migration already applied
	var count int
	query := "SELECT COUNT(*) FROM migrations WHERE version = $1"
	err := ms.db.QueryRow(query, migration.Version()).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		ms.logger.Infof("Migration %d already applied", migration.Version())
		return nil
	}

	// Run migration
	if err := migration.Up(); err != nil {
		return fmt.Errorf("failed to run migration %d: %v", migration.Version(), err)
	}

	// Record migration
	insertQuery := "INSERT INTO migrations (version, description) VALUES ($1, $2)"
	_, err = ms.db.Exec(insertQuery, migration.Version(), migration.Description())
	if err != nil {
		return err
	}

	ms.logger.Infof("Migration %d applied successfully: %s", migration.Version(), migration.Description())
	return nil
}

// Built-in database drivers
const (
	DriverPostgreSQL = "postgres"
	DriverMySQL      = "mysql"
	DriverSQLite     = "sqlite3"
	DriverSQLServer  = "sqlserver"
)
