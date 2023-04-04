package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config contains the connection details for a PostgreSQL database
type Config struct {
	Host     string // Hostname or IP address of the PostgreSQL server
	Port     string // Port number that PostgreSQL is running on
	Password string // Password for the database user
	User     string // Database user to connect as
	DBName   string // Name of the PostgreSQL database to connect to
	SSLMode  string // SSL mode to use for the connection (e.g. "disable", "require", etc.)
}

// NewConnection creates a new connection to a PostgreSQL database using the given Config
func NewConnection(config *Config) (*gorm.DB, error) {
	// Construct a DSN string from the provided Config
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Open a new connection to the PostgreSQL database using the DSN and gorm.Open
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Return the DB object and error if any
	return db, err
}
