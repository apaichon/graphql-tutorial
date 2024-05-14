package data

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"graphql-api/config"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// DB represents the PostgreSQL database
type PostgresDB struct {
	Connection *sql.DB
}

var postgresInstance *PostgresDB
var postgresOnce sync.Once

// NewDB initializes a new instance of the DB struct
func NewPostgresDB() *PostgresDB {
	postgresOnce.Do(func() {
		config := config.NewConfig()
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword, viper.GetString("METRIC_DB"))
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		postgresInstance = &PostgresDB{conn}
	})
	return postgresInstance
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	if db.Connection == nil {
		return nil
	}
	return db.Connection.Close()
}

// Insert inserts data into the specified table
func (db *PostgresDB) Insert(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

// Query executes a query and returns rows
func (db *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return rows, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	row := db.Connection.QueryRow(query, args...)
	return row
}

// Delete executes a delete statement
func (db *PostgresDB) Delete(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

// Update executes an update statement
func (db *PostgresDB) Update(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}

	return result, nil
}

func (db *PostgresDB) Begin() (*sql.Tx, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return tx, nil
}

func (db *PostgresDB) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	return stmt, nil
}

func (db *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Connection.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}
	return result, nil
}