package router

import (
	"context"
	"fmt"
	"net/http"
)

// TODO: Make /users/:id and /users/:uuid point to the same router
// TODO: Implement Middleware struct

func New(p string) *Mux {
	return &Mux{
		Prefix: NewPath(p),
		Routes: make(map[Path]Router),
		// Error:  ErrorHandlers{},
	}
}

type Mux struct {
	Prefix Path
	Routes map[Path]Router
	Error  ErrorHandlers
	ctx    context.Context
}

type ErrorHandlers struct {
	NotFound         http.Handler
	MethodNotAllowed http.Handler
}

type ErrorHandlersKey struct{}

// Creates a new Mux bound to the parent one.
// It has the same context as the parent and is relative to the parent prefix.
func (mux *Mux) NewSubrouter(p string) *Mux {
	prefix := mux.Prefix.Join(p)

	if _, ok := mux.Routes[prefix]; ok {
		panic(fmt.Sprintf("Duplicate subrouter for prefix: '%s'", prefix))
	}

	subrouter := New(string(prefix))
	subrouter.ctx = mux.Context()
	mux.Routes[subrouter.Prefix] = subrouter
	return subrouter
}

// Attempts to match the given request path against the registered routes.
func (mux *Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reqPath := NewPath(req.URL.Path)

	if req.URL.Path != string(reqPath) {
		// If after cleanPath, the paths are different, redirect; This is needed for SEO
		redirect(res, req, string(reqPath))
		return
	}

	reqPathFields := reqPath.Fields()

	for path, route := range mux.Routes {
		if path.Match(reqPathFields) {
			route.ServeHTTP(res, req)
			return
		}
	}

	if mux.Error.NotFound != nil {
		mux.Error.NotFound.ServeHTTP(res, req)
		return
	}

	ErrorHandlers := mux.Context().Value(ErrorHandlersKey{}).(ErrorHandlers)
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

// Returns the mux's context.
//
// The returned context is always non-nil; it defaults to the background context with the default handlers value.
// To retrieve the default handlers, access the context value with ErrorHandlersKey struct
func (mux *Mux) Context() context.Context {
	if mux.ctx != nil {
		return mux.ctx
	}

	return context.WithValue(context.Background(), ErrorHandlersKey{}, mux.Error)
}

// ----------------------------------------------------------------------------
// Route factories
// ----------------------------------------------------------------------------

func (mux *Mux) Use(p string, methods []string, handler http.Handler) {
	routePath := mux.Prefix.Join(p)

	route, ok := mux.Routes[routePath]

	if !ok {
		route = NewRoute(routePath, mux.Context())
		mux.Routes[routePath] = route
	}

	route.Use(string(routePath), methods, handler)
}

func (mux *Mux) Connect(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodConnect}, handler)
}

func (mux *Mux) Delete(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodDelete}, handler)
}

func (mux *Mux) Get(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodGet}, handler)
}

func (mux *Mux) Head(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodHead}, handler)
}

func (mux *Mux) Options(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodOptions}, handler)
}

func (mux *Mux) Patch(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodPatch}, handler)
}

func (mux *Mux) Post(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodPost}, handler)
}

func (mux *Mux) Put(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodPut}, handler)
}

func (mux *Mux) Trace(path string, handler http.Handler) {
	mux.Use(path, []string{http.MethodTrace}, handler)
}
