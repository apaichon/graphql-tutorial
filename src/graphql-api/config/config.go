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
    GraphQLPort int
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

         // Parse DB port as integer
         graphQlPortStr := os.Getenv("GRAPHQL_PORT")
         if graphQlPortStr != "" {
             graphQLPort, err := strconv.Atoi(graphQlPortStr)
             if err != nil {
                 fmt.Println("Invalid GRAPHQL_PORT value:", err)
                 return
             }
             instance.GraphQLPort = graphQLPort
         }
    })
    return instance
}

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
    return instance
}
