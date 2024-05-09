package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)


// Higher-order function that wraps a resolver with timeout handling
func CircuitBreakerResolver(resolver graphql.FieldResolveFn, timeout time.Duration) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(p.Context, timeout)
		defer cancel()

		// Create a channel to signal completion or timeout
		done := make(chan struct{}, 1)
		defer close(done)

		// Execute the resolver within a goroutine
		go func() {
			// Call the resolver
			_, err := resolver(p)
			if err != nil {
				// If resolver returns an error, send it through the channel
				done <- struct{}{}
			}
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			return nil, errors.New("resolver execution error")
		case <-ctx.Done():
			// Timeout occurred
			return nil, errors.New("timeout occurred")
		}
	}
}
// Middleware function signature
type Middleware func(http.Handler) http.Handler

func CircuitBreakerMiddleware(timeout time.Duration) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Create a context with timeout
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()

            // Create a channel to signal completion or timeout
            done := make(chan struct{}, 1)
            // defer close(done)

            // Execute the handler within a goroutine
            go func() {
                next.ServeHTTP(w, r)
                done <- struct{}{}
            }()

            // Wait for either completion or timeout
            select {
            case <-done:
                return
            case <-ctx.Done():
                // Timeout occurred
                http.Error(w, "Request timeout", http.StatusRequestTimeout)
            }
        })
    }
}