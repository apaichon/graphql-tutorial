package main

import (
	"fmt"
	"net/http"
	"time"
	"context"
	"os"
	"os/signal"
	"log"

	"graphql-api/config"
	"graphql-api/internal/auth"
	"graphql-api/internal/middleware"
	"graphql-api/internal/monitoring"
	gql "graphql-api/pkg/graphql"

	"golang.org/x/time/rate"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/viper"
)

func main() {
	
	shutdown, err :=  monitoring.InitTracer(viper.GetString("TRACE_EXPORTER_URL"))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
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
	return auth.AuthenticationHandler(circuitBreaker)
}
