package cache

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteOnce     sync.Once
	sqliteInstance *SQLiteInMemClient
)

// SQLiteClient represents a simple SQLite client.
type SQLiteInMemClient struct {
	db *sql.DB
}

// NewSQLiteClient creates a new SQLite client with an in-memory database.
func NewSQLiteInMemClient() (*SQLiteInMemClient, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %v", err)
	}
	return &SQLiteInMemClient{db: db}, nil
}

// GetInstance returns the singleton instance of the SQLite client.
func GetSqliteInMemInstance() (*SQLiteInMemClient, error) {
	sqliteOnce.Do(func() {
		var err error
		sqliteInstance, err = NewSQLiteInMemClient()
		if err != nil {
			log.Fatalf("Error creating SQLite client: %v", err)
		}
	})
	return sqliteInstance, nil
}

// Close closes the SQLite database connection.
func (sc *SQLiteInMemClient) Close() error {
	return sc.db.Close()
}

// CreateTable creates a table in the SQLite database.
func (sc *SQLiteInMemClient) CreateTable() error {
	_, err := sc.db.Exec(`CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value TEXT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table in SQLite database: %v", err)
	}
	return nil
}

// Get retrieves the value associated with the given key from the SQLite database.
func (sc *SQLiteInMemClient) Get(key string) (string, error) {
	var value string
	err := sc.db.QueryRow("SELECT value FROM cache WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("key '%s' not found", key)
		}
		return "", fmt.Errorf("failed to get value for key '%s': %v", key, err)
	}
	return value, nil
}

// Set sets the value associated with the given key in the SQLite database.
func (sc *SQLiteInMemClient) Set(key, value string) error {
	_, err := sc.db.Exec("INSERT INTO cache(key, value) VALUES(?, ?)", key, value)
	if err != nil {
		return fmt.Errorf("failed to set value for key '%s': %v", key, err)
	}
	return nil
}

// Remove removes the specified key from the SQLite database.
func (sc *SQLiteInMemClient) Remove(key string) error {
	res, err := sc.db.Exec("DELETE FROM cache WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to remove key '%s' from SQLite database: %v", key, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows after removing key '%s': %v", key, err)
	}
	if count == 0 {
		return fmt.Errorf("key '%s' does not exist in SQLite database", key)
	}
	return nil
}

// Remove removes the specified key from the SQLite database.
func (sc *SQLiteInMemClient) Removes(key string)  {
	sc.db.Exec("DELETE FROM cache WHERE key Like ?", "%"+key+"*%")
}