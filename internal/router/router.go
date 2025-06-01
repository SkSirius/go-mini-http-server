package router

import (
	"context"
	"net/http"
	"strings"
)

// key is a custom type for context keys to avoid collisions.
type key int

// paramsKey is the context key used to store route parameters.
const paramsKey key = 0

// WithParams returns a new context with the given params map stored in it.
func WithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

// Params retrieves the route parameters from the context.
// Returns an empty map if no parameters are found.
func Params(ctx context.Context) map[string]string {
	if val, ok := ctx.Value(paramsKey).(map[string]string); ok {
		return val
	}
	return map[string]string{}
}

// Param retrieves a single parameter by name from the context.
func Param(ctx context.Context, name string) string {
	return Params(ctx)[name]
}

// Router holds the registered routes, organized by HTTP method.
type Router struct {
	routes map[string][]routeEntry
}

// routeEntry represents a single route pattern and its handler.
type routeEntry struct {
	pattern string
	handler http.Handler
}

// New creates and returns a new Router instance.
func New() *Router {
	return &Router{
		routes: make(map[string][]routeEntry),
	}
}

// Handle registers a new route with a method, path pattern, and handler.
func (r *Router) Handle(method, path string, handler http.Handler) {
	r.routes[method] = append(r.routes[method], routeEntry{
		pattern: path,
		handler: handler,
	})
}

// ServeHTTP implements the http.Handler interface for Router.
// It matches incoming requests to registered routes and dispatches them.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	entries, ok := r.routes[req.Method]
	if !ok {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	for _, entry := range entries {
		if params, ok := match(entry.pattern, req.URL.Path); ok {
			// Attach route parameters to the request context.
			req = req.WithContext(WithParams(req.Context(), params))
			entry.handler.ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

// match checks if the given path matches the pattern.
// It supports patterns with named parameters (e.g., "/users/:id").
// Returns a map of parameter names to values if matched, or false otherwise.
func match(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			// This part is a parameter, extract its value.
			params[patternParts[i][1:]] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			// Not a match.
			return nil, false
		}
	}

	return params, true
}
