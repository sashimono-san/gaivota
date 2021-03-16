package mux

import (
	"context"
	"fmt"
	"net/http"
)

// TODO: Implement Middleware struct

func New(p string) *Router {
	return &Router{
		Prefix: NewPath(p),
		Routes: make(map[Path]http.Handler),
		// Error:  ErrorHandlers{},
	}
}

type Router struct {
	Prefix     Path
	Routes     map[Path]http.Handler
	subRouters []Path
	Error      ErrorHandlers
	ctx        context.Context
}

type ErrorHandlers struct {
	NotFound         http.Handler
	MethodNotAllowed http.Handler
}

type ErrorHandlersKey struct{}

// Creates a new Router bound to the parent one.
// It has the same context as the parent and is relative to the parent prefix.
func (router *Router) NewSubrouter(p string) *Router {
	prefix := router.Prefix.Join(p)

	if path, ok := router.FindSubrouter(prefix); ok {
		return router.Routes[path].(*Router).NewSubrouter(p)
	}

	if _, ok := router.Routes[prefix]; ok {
		panic(fmt.Sprintf("Duplicate subrouter for prefix: '%s'", prefix))
	}

	router.subRouters = append(router.subRouters, prefix)

	subrouter := New(string(prefix))
	subrouter.ctx = router.Context()
	router.Routes[subrouter.Prefix] = subrouter

	return subrouter
}

func (router *Router) FindSubrouter(path Path) (Path, bool) {
	founPath := NewPath("")
	ok := false

	for _, candidatePath := range router.subRouters {
		// Only test prefix if candidate path is bigger than current one
		if len(founPath) < len(candidatePath) {
			if path.HasPrefix(candidatePath) {
				founPath = candidatePath
				ok = true
			}
		}
	}

	return founPath, ok
}

// Attempts to match the given request path against the registered routes.
func (router *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reqPath := NewPath(req.URL.Path)

	if req.URL.Path != string(reqPath) {
		// If after cleanPath, the paths are different, redirect. This is needed for SEO
		redirect(res, req, string(reqPath))
		return
	}

	if path, ok := router.FindSubrouter(reqPath); ok {
		router.Routes[path].ServeHTTP(res, req)
		return
	}

	fmt.Printf("Got %s request for %s\n", req.Method, req.URL.Path)

	for path, route := range router.Routes {
		if path.Match(reqPath) {
			route.ServeHTTP(res, req)
			return
		}
	}

	if router.Error.NotFound != nil {
		router.Error.NotFound.ServeHTTP(res, req)
		return
	}

	ErrorHandlers := router.Context().Value(ErrorHandlersKey{}).(ErrorHandlers)
	if ErrorHandlers.NotFound != nil {
		ErrorHandlers.NotFound.ServeHTTP(res, req)
		return
	}

	http.NotFound(res, req)
}

func redirect(res http.ResponseWriter, req *http.Request, toPath string) {
	rc := *req.URL
	rc.User = req.URL.User
	rc.Path = toPath

	// Use 'moved permanently' for HEAD|GET requests.
	code := http.StatusPermanentRedirect
	if req.Method != http.MethodHead && req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}

	http.Redirect(res, req, rc.String(), code)
}

// Returns the router's context.
//
// The returned context is always non-nil; it defaults to the background context with the default handlers value.
// To retrieve the default handlers, access the context value with ErrorHandlersKey struct
func (router *Router) Context() context.Context {
	if router.ctx != nil {
		return router.ctx
	}

	return context.WithValue(context.Background(), ErrorHandlersKey{}, router.Error)
}

// ----------------------------------------------------------------------------
// Route factories
// ----------------------------------------------------------------------------

func (router *Router) Use(p string, methods []string, handler http.Handler) {
	routePath := router.Prefix.Join(p)

	if path, ok := router.FindSubrouter(routePath); ok {
		router.Routes[path].(*Router).Use(p, methods, handler)
		return
	}

	route, ok := router.Routes[routePath]

	if !ok {
		fmt.Printf("Creating new Route at %s\n", routePath)
		route = NewRoute(routePath, router.Context())
		router.Routes[routePath] = route
	}

	route.(*Route).Use(methods, handler)
}

func (router *Router) Connect(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodConnect}, handler)
}

func (router *Router) Delete(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodDelete}, handler)
}

func (router *Router) Get(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodGet}, handler)
}

func (router *Router) Head(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodHead}, handler)
}

func (router *Router) Options(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodOptions}, handler)
}

func (router *Router) Patch(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodPatch}, handler)
}

func (router *Router) Post(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodPost}, handler)
}

func (router *Router) Put(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodPut}, handler)
}

func (router *Router) Trace(path string, handler http.Handler) {
	router.Use(path, []string{http.MethodTrace}, handler)
}
