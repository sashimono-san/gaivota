package mux

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
		Handlers:  make(map[string]http.Handler),
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

		fmt.Printf("Handling %s request for %s with route at %s\n", req.Method, req.URL.Path, route.Path)

		// Add path params to request context before calling the handler
		params := route.ExtractParams(req.URL.Path)
		req = requestWithPathParams(req, params)

		handler.ServeHTTP(res, req)
		return
	}

	fmt.Printf("Request method %s is not supported by route at %s\n", req.Method, route.Path)

	errorHandlers := route.ctx.Value(ErrorHandlersKey{}).(ErrorHandlers)
	if errorHandlers.MethodNotAllowed != nil {
		errorHandlers.MethodNotAllowed.ServeHTTP(res, req)
		return
	}

	http.NotFound(res, req)
}

func (route *Route) Use(methods []string, handler http.Handler) {
	for _, method := range methods {
		if _, ok := route.Handlers[method]; ok {
			panic(fmt.Sprintf("Duplicate handler for %s method for route: '%s'", method, route.Path))
		}

		route.Handlers[method] = handler
	}
}

func (route *Route) ExtractParams(path string) Params {
	pathFields := Path(path).Fields()

	params := Params{}
	for param, pos := range route.ParamsPos {
		params[param] = pathFields[pos]
	}

	return params
}
