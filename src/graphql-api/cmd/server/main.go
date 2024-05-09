package main

import (
	"fmt"
	"net/http"
	"time"

	"graphql-api/config"
	"graphql-api/internal/auth"
	"graphql-api/internal/logger"
	"graphql-api/internal/middleware"
	gql "graphql-api/pkg/graphql"


	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {

	// Load configuration
	config := config.NewConfig()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    gql.RootQuery,
		Mutation: gql.RootMutation,
	})
	if err != nil {
		panic(err)
	}
	// Create a GraphQL handler for HTTP requests
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false, // Disable GraphiQL for subscriptions endpoint
		Playground: true,
	})

	go moveAuditLog(config)
	// Serve GraphQL API at /graphql endpoint
	http.Handle("/graphql", auth.AuthenticationHandler(middleware.AuditLogMiddleware(graphqlHandler)))
	http.HandleFunc("/login", auth.LoginHandler)
	// Start the HTTP server
	fmt.Printf(`Server is running at http://localhost:%v/graphql`, config.GraphQLPort)
	http.ListenAndServe(fmt.Sprintf(`:%v`, config.GraphQLPort), nil)
	
}

func moveAuditLog(cfg *config.Config) {
	auditLog:= logger.GetLogInitializer()
	 for {
		fmt.Println("Run Move Audit Log")
		auditLog.MoveLogsToSQLite()
		fmt.Println("End Move Audit Log")
		time.Sleep(time.Duration(cfg.LogMoveMin) * time.Minute)
	 }
}
