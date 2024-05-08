package logger

import (
	"graphql-api/pkg/data/models"
	"os"
	"fmt"
	"path/filepath"
	"sync"
	"log"
	"encoding/json"
)

// Write log to a file
func WriteLogToFile(logChan <-chan models.LogModel, wg *sync.WaitGroup) {
	defer wg.Done()

	var currentLogFile *os.File
	var currentLogFilePath string

	for logEntry := range logChan {
		// Determine the log file name based on the current timestamp (every 5 minutes)
		logFileName := fmt.Sprintf("%04d-%02d-%02d-%02d_%02d.log",
			logEntry.Timestamp.Year(),
			logEntry.Timestamp.Month(),
			logEntry.Timestamp.Day(),
			logEntry.Timestamp.Hour(),
			(logEntry.Timestamp.Minute()/5)*5,
		)

		logFilePath := filepath.Join("../../logs", logFileName)

		// Check if we're still using the same log file, or need to switch
		if currentLogFilePath != logFilePath {
			if currentLogFile != nil {
				currentLogFile.Close() // Close the previous log file
			}
			currentLogFilePath = logFilePath

			// Ensure the "logs" directory exists
			if err := os.MkdirAll("logs", 0755); err != nil {
				log.Fatalf("Error creating logs directory: %v", err)
			}

			var err error
			currentLogFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error opening log file: %v", err)
			}
		}

		// Write the log entry as JSON to the file
		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			log.Fatalf("Error marshaling log data: %v", err)
		}

		_, err = currentLogFile.Write(append(jsonData, '\n')) // Add newline
		if err != nil {
			log.Fatalf("Error writing to log file: %v", err)
		}
	}
}