# Lab3 - Module 4: Query
## Prepare Database and sampling data.
1. Prepare sqlite database and table.
- Create file "event.db" on folder "data".
- Connect by SQL Tools.
- Run sql create database.
```sql
CREATE TABLE IF NOT EXISTS contact (
  contact_id INTEGER PRIMARY KEY,
  name TEXT,
  first_name TEXT,
  last_name TEXT,
  gender_id INTEGER,
  dob DATE,
  email TEXT,
  phone TEXT,
  address TEXT,
  photo_path TEXT,
  created_at DATETIME,
  created_by TEXT
);
```
- Run sql sampling data
```sql
INSERT INTO contact (contact_id, name, first_name, last_name, gender_id, dob, email, phone, address, photo_path, created_at, created_by)
VALUES 
(1, 'John Doe', 'John', 'Doe', 1, '1990-01-01', 'john@example.com', '123-456-7890', '123 Main St', 'path/to/photo.jpg', '2024-04-16 12:00:00', 'Admin'),
(2, 'Jane Smith', 'Jane', 'Smith', 2, '1992-05-15', 'jane@example.com', '987-654-3210', '456 Oak St', 'path/to/photo2.jpg', '2024-04-16 12:30:00', 'Admin');

```
2. Create Configuration object to keep database configuration on config/config.go
```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "sync"

    "github.com/joho/godotenv"
)

// Config represents application configuration
type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
	DBName string
    SecretKey string
    TokenAge int
}

var instance *Config
var once sync.Once

// LoadConfig loads the configuration from environment variables
func NewConfig() *Config {
    once.Do(func() {
        // Load environment variables from .env file
        if err := godotenv.Load(); err != nil {
            fmt.Println("Failed to load env variables:", err)
            return
        }

        // Parse environment variables and create config object
        instance = &Config{
            DBHost:     os.Getenv("DB_HOST"),
            DBUser:     os.Getenv("DB_USER"),
            DBPassword: os.Getenv("DB_PASSWORD"),
			DBName: os.Getenv("DB_NAME"),
            SecretKey: os.Getenv("SECRET_KEY"),
        }

        // Parse DB port as integer
        dbPortStr := os.Getenv("DB_PORT")
        if dbPortStr != "" {
            dbPort, err := strconv.Atoi(dbPortStr)
            if err != nil {
                fmt.Println("Invalid DB_PORT value:", err)
                return
            }
            instance.DBPort = dbPort
        }
        tokenAge, err :=strconv.Atoi(os.Getenv("TOKEN_AGE"))
        if (err != nil) {
            fmt.Println("Invalid TOKEN_AGE value:", err)
            return
        }
        instance.TokenAge = tokenAge
    })
    return instance
}

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
    return instance
}

```
3. Install go library for read environtment variable.
```sh
 go get github.com/joho/godotenv
```
4. Create .env at /cmd/server/.env and input configuration following.
```sh
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=<dbpath>
SECRET_KEY=your-secret-key
TOKEN_AGE=60
```
5. Create Database object to manage sqlite database at database/db.go
```go
package database

import (
    "database/sql"

    "fmt"
	"log"
	"sync"
    _ "github.com/mattn/go-sqlite3"
	"github.com/apaichon/graphql-tutorial/config"
)

// DB represents the SQLite database
type DB struct {
    Connection *sql.DB
}

var instance *DB
var once sync.Once


// NewDB initializes a new instance of the DB struct
func NewDB() (*DB) {
	once.Do (func() {
		config := config.NewConfig()
		conn, err := sql.Open("sqlite3", config.DBName)
		if err != nil {
			log.Fatal(err)
		}
		instance = &DB{conn}
	})

    return instance
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

```
6. Download Sqlite Libary github.com/mattn/go-sqlite3
```sh
go get github.com/mattn/go-sqlite3
```

Create base and contact repository

Create test
Run test


## Create GraphQL Query 

Schema

- Queries

    - ContactQuery
    
        - Type
            
            - Fields
        
        - Agruments
            
            - (searchText, limit, offset)
        
        - Resolve
            
            - ContactRepo.GetContactsByTextSearch(Agruments)
 
 