package config

import (
	"fmt"
	"path/filepath"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	relativePath := "../../config/.env"

		// Get the absolute path
		absolutePath, err := filepath.Abs(relativePath)
        if err != nil {
            fmt.Println(err)
            return
        }
		// Load environment variables from .env file
		if err := godotenv.Load(absolutePath); err != nil {
			fmt.Println("Failed to load env variables:", err)
			return
		}

		viper.AutomaticEnv()
}
