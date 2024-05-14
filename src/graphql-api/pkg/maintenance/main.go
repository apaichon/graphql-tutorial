package main

import (
	"os"
	"time"

	"graphql-api/config"
	"graphql-api/internal/logger"
	"log"
)

func main() {
	// Load configuration
	config := config.NewConfig()
	go moveAuditLog(config)
	// Keep the main goroutine alive
	select {}
}

func moveAuditLog(cfg *config.Config) {
	auditLog := logger.GetLogInitializer()
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	for {
		logger.Println("[info] - Run Move Audit Log")
		auditLog.MoveLogsToPostgres()
		logger.Println("[info] - End Move Audit Log")
		time.Sleep(time.Duration(cfg.LogMoveMin) * time.Minute)
	}
}
