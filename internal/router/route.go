package router

import (
	"context"
	"fmt"
	"net/http"
)

func NewRoute(path Path, ctx context.Context) *Route {
	paramsPos := path.extractParamsPos()

	return &Route{
		Path:      path,
		ParamsPos: paramsPos,
		ctx:       ctx,
	}
}

type Route struct {
	Path      Path
	ParamsPos ParamsPos
	Handlers  map[string]http.Handler
	ctx       context.Context
}

func (route *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if handler, ok := route.Handlers[req.Method]; ok {

		// Add path params to request context before calling the handler
		params := route.ExtractParams(req.URL.Path)
		req = requestWithPathParams(req, params)

		handler.ServeHTTP(res, req)
		return
	}

	errorHandlers := route.ctx.Value(ErrorHandlersKey{}).(ErrorHandlers)
	if errorHandlers.MethodNotAllowed != nil {
		errorHandlers.MethodNotAllowed.ServeHTTP(res, req)
		return
	}

	http.NotFound(res, req)
}

func (route *Route) Use(path string, methods []string, handler http.Handler) {
	handlers := make(map[string]http.Handler)
	for _, method := range methods {
		if _, ok := route.Handlers[method]; ok {
			panic(fmt.Sprintf("Duplicate handler for %s method for route: '%s'", method, path))
		}

		handlers[method] = handler
	}

	route.Handlers = handlers
}

func (route *Route) ExtractParams(path string) Params {
	pathFields := Path(path).Fields()

	params := Params{}
	for param, pos := range route.ParamsPos {
		params[param] = pathFields[pos]
	}

	return params
}
