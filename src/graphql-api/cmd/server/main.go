package main

import (
	"fmt"
	"net/http"

	"graphql-api/config"
	"graphql-api/internal/auth"
	gql "graphql-api/pkg/graphql"

	"graphql-api/internal/middleware"

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

	// Serve GraphQL API at /graphql endpoint

	// authenticated := auth.AuthenticationHandler(graphqlHandler)
	// auditlog := middleware.AuditLogMiddleware(authenticated)
	// http.Handle("/graphql", auditlog)
	http.Handle("/graphql", auth.AuthenticationHandler(middleware.AuditLogMiddleware(graphqlHandler)))
	http.HandleFunc("/login", auth.LoginHandler)

	// Start the HTTP server
	fmt.Printf(`Server is running at http://localhost:%v/graphql`, config.GraphQLPort)
	http.ListenAndServe(fmt.Sprintf(`:%v`, config.GraphQLPort), nil)

}
