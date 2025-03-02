package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Let's you stack middlewares.
// First in is executed first (wraps all following).
func With(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			mw := middlewares[i]
			next = mw(next)
		}
		return next
	}
}
