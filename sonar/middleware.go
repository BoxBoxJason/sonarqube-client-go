package sonar

import (
	"net/http"
	"slices"
)

// Middleware is a function that wraps an http.RoundTripper.
// Use WithMiddleware to attach middleware to the client transport chain.
type Middleware func(http.RoundTripper) http.RoundTripper

// applyMiddlewares wraps base with each middleware in reverse order so that
// middlewares[0] is the outermost (first to handle a request).
func applyMiddlewares(base http.RoundTripper, middlewares []Middleware) http.RoundTripper {
	transport := base
	for _, mw := range slices.Backward(middlewares) {
		transport = mw(transport)
	}

	return transport
}
