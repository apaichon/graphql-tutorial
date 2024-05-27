package data

import (
	"database/sql"

	"fmt"
	"graphql-api/config"
	"log"
	"sync"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

// DB represents the SQLite database
type DB struct {
	Connection *sql.DB
}

var instance *DB
var once sync.Once

// NewDB initializes a new instance of the DB struct
func NewDB() *DB {
	once.Do(func() {
		config := config.NewConfig()
		fmt.Println(config.DBName)
		conn, err := sql.Open("sqlite3", config.DBName)
		if err != nil {
			fmt.Printf("err:%v", err)
			log.Fatal(err)
		}
		instance = &DB{conn}
	})
	return instance
}

func (db *DB) Open() error {
	if db.Connection == nil {
		config := config.NewConfig()
		conn, err := sql.Open("sqlite3", config.DBName)
		if err != nil {
			return err
		}
		instance = &DB{conn}
	}
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.Connection == nil {
		return nil
	}
	return db.Connection.Close()
}

// Insert inserts data into the specified table
func (db *DB) Insert(query string, args ...interface{}) (sql.Result, error) {
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
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return rows, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *DB) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	row := db.Connection.QueryRow(query, args...)
	return row, nil
}

// Delete executes a delete statement
func (db *DB) Delete(query string, args ...interface{}) (sql.Result, error) {
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
func (db *DB) Update(query string, args ...interface{}) (sql.Result, error) {
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

func (db *DB) Begin() (*sql.Tx, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return tx, nil
}

func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	return stmt, nil
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Connection.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %v", err)
	}
	return result, nil
}
