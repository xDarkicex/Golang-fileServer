package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routes = []func(*Context){}

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router{router: mux.NewRouter()}
}
func (r *Router) ListenAndServe(address string) error {
	server := &http.Server{
		Handler: r.router,
		Addr:    address,
	}
	return server.ListenAndServe()
}

func (r *Router) Group(template string) *Router {
	return &Router{router: r.router.PathPrefix(template).Subrouter()}
}

func (r *Router) route(path string, method string, handler func(*Context)) {
	handle := func(response http.ResponseWriter, request *http.Request) {
		ctx := NewContext(response, request)
		for _, route := range routes {
			route(ctx)
		}
		handler(ctx)
	}
	if path == "" {
		r.router.Methods(method).HandlerFunc(handle)
	} else {
		r.router.Methods(method).Path(path).HandlerFunc(handle)
	}
}

func (r *Router) GET(path string, handler func(*Context)) {
	r.route(path, "GET", handler)
}

func (r *Router) POST(path string, handler func(*Context)) {
	r.route(path, "POST", handler)

}
func (r *Router) OPTIONS(path string, handler func(*Context)) {
	r.route(path, "OPTIONS", handler)
}

func (r *Router) PUT(path string, handler func(*Context)) {
	r.route(path, "PUT", handler)
}
func (r *Router) AllRoutes(f func(*Context)) {
	routes = append(routes, f)
}

func ContructRequest(c *Context) {

	c.Header("Access-Controll-Allow-Methods", "GET, POST")
	c.Header("Access-Controll-Allow-Max-Age", "86400")
	c.Header("Access-Controll-Allow-Credentials", "true")
	c.Header("Access-Controll-Allow-Origin", "*")
	c.Header("Access-Controll-Allow-Headers", "content-type, authorization")

}

func DestructRequest(c *Context) {
	c.Header("Access-Controll-Allow-Origin", "*")

}

func (r *Router) Static(pathPrefix, path string) {
	r.router.PathPrefix(pathPrefix).Handler(http.FileServer(http.Dir(path)))
}
