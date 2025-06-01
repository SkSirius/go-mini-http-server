package router

import "net/http"

type Router struct {
	routes map[string]map[string]http.Handler
}

func New() *Router {
	return &Router{
		routes: make(map[string]map[string]http.Handler),
	}
}

func (r *Router) Handle(method, path string, handler http.Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.Handler)
	}
	r.routes[method][path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodRoutes, ok := r.routes[req.Method]
	if !ok {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	handler, ok := methodRoutes[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}

	handler.ServeHTTP(w, req)
}
