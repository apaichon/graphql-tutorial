package middleware

import (
	"net/http"
	"golang.org/x/time/rate"
)

// Define the rate limiter middleware
func RateLimitMiddleware(limit rate.Limit, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(limit, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request exceeds the rate limit
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Proceed to the next handler if the rate limit is not exceeded
			next.ServeHTTP(w, r)
		})
	}
}
