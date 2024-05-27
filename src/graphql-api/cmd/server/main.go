package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	"graphql-api/config"
	"graphql-api/internal/auth"
	"graphql-api/internal/middleware"

	"graphql-api/internal/monitoring"
	"graphql-api/internal/subscription"
	gql "graphql-api/pkg/graphql"
	"graphql-api/pkg/graphql/directives"
)

var cfg *config.Config

func init() {
	// Load configuration
	cfg = config.NewConfig()
}

func main() {
	shutdown, err := monitoring.InitTracer(viper.GetString("TRACE_EXPORTER_URL"))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Channel to signal when the server has shut down
	done := make(chan bool, 1)

	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        gql.RootQuery,
		Mutation:     gql.RootMutation,
		Subscription: gql.RootSubscription,
		Directives:   append(graphql.SpecifiedDirectives, directives.SubstringDirective),
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
	http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/subscribe", subscription.NewMessageHandler)
	http.HandleFunc("/ws-subscribe", middleware.CorsHandler(subscription.SubscribeWsHandler(schema)))

	// Create the server
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", cfg.GraphQLPort),
	}

	go func() {
		// Start the HTTP server
		fmt.Printf("Server is running at http://localhost:%v/graphql\n", cfg.GraphQLPort)
		server.ListenAndServe()

	}()

	// Listen for interrupt signal
	<-stop
	fmt.Println("Shutting down server...")

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	close(done)

	// Wait for the server to shutdown
	<-done
	fmt.Println("Server stopped")

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
	limit := rate.Every((time.Duration(rateLimitReqSec) * time.Second))
	execTimeOut := viper.GetInt("EXEC_TIME_OUT")
	auditLog := middleware.AuditLogMiddleware(graphqlHandler)
	rateLimit := middleware.RateLimitMiddleware(limit, rateLimitBurst)(auditLog)
	circuitBreaker := middleware.CircuitBreakerMiddleware(time.Duration(execTimeOut) * time.Second)(rateLimit)
	return middleware.CorsHandler(auth.AuthenticationHandler(circuitBreaker))

}
