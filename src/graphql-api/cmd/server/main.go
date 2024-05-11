package main

import (
	"fmt"
	"net/http"
	"time"

	"graphql-api/config"
	"graphql-api/internal/auth"
	"graphql-api/internal/middleware"
	gql "graphql-api/pkg/graphql"

	"golang.org/x/time/rate"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/viper"
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
	http.Handle("/graphql", handlers(graphqlHandler))
	// http.Handle("/graphql", graphqlHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	// Start the HTTP server
	fmt.Printf(`Server is running at http://localhost:%v/graphql`, config.GraphQLPort)
	http.ListenAndServe(fmt.Sprintf(`:%v`, config.GraphQLPort), nil)

}

func handlers(graphqlHandler *handler.Handler) http.Handler {
	/* High Order Functions
	* Authen
	* RateLimit
	* AuditLog
	* GraphQLHanler

	 */
	rateLimitReqSec := viper.GetInt("RATE_LIMIT_REQ_SEC")
	rateLimitBurst := viper.GetInt("RATE_LIMIT_BURST")
	limit := rate.Every( (time.Duration(rateLimitReqSec) * time.Second))
	execTimeOut := viper.GetInt("EXEC_TIME_OUT")
	auditLog := middleware.AuditLogMiddleware(graphqlHandler)
	rateLimit := middleware.RateLimitMiddleware(limit, rateLimitBurst)(auditLog)
	circuitBreaker := middleware.CircuitBreakerMiddleware(time.Duration(execTimeOut) * time.Second)(rateLimit)
	// return auth.AuthenticationHandler(circuitBreaker)
	return circuitBreaker
}
