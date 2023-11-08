package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Struct that encapsulates the sql DB object
// Lowercase db means it is private so it is not available to packages outside of db
type Database struct {
	db *sql.DB
}

// Returns a pointer to the database struct
func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgres://root:password@localhost:5433/go-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// Create an additional method so we can close the database
// Because our db is encapsulated, th only way to access it is through the methods we define

func (d *Database) GetDB() *sql.DB {
	return d.db
}
