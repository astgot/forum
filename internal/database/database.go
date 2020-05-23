package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database ...
type Database struct {
	config     *Config
	db         *sql.DB
	usersStore *UsersStore
}

// NewDB ...
func NewDB(config *Config) *Database {
	return &Database{
		config: config,
	}
}

// InitDB ...
func (d *Database) InitDB() error {
	db, err := sql.Open("sqlite3", d.config.DatabaseAddress)
	if err != nil {
		return err
	}
	// Complete checking of connection with DB
	if err := db.Ping(); err != nil {
		return err
	}
	d.db = db // fill "db" field with completely configured DB
	d.BuildSchema()
	return nil
}

// Close DB
func (d *Database) Close() {
	d.db.Close()
}

// User ---> to use Users' info in the service --->
func (d *Database) User() *UsersStore {
	if d.usersStore != nil {
		return d.usersStore
	}
	d.usersStore = &UsersStore{
		database: d,
	}
	return d.usersStore
}
