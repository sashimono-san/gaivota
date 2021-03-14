package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	t.Parallel()

	mux := &Mux{
		subRouters: []Path{Path("/investments/:id/positions")},
		Prefix: Path("/"),
		Routes: map[Path]Router{
			// Test simple route
			Path("/investments"): &Route{
				Handlers: map[string]http.Handler{
					http.MethodGet: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
						res.Write([]byte("/investments"))
					}),
				},
			},
			// Test redirect
			Path("/investments/"): &Route{
				Handlers: map[string]http.Handler{
					http.MethodGet: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
						res.Write([]byte("/investments/"))
					}),
				},
			},
			// Test route with param
			Path("/investments/:id"): &Route{
				Handlers: map[string]http.Handler{
					http.MethodGet: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
						res.Write([]byte("/investments/:id"))
					}),
				},
			},
			// Test subrouter
			Path("/investments/:id/positions"): &Mux{
				Routes: map[Path]Router{
					// Test route with same prefix as subrouter
					Path("/investments/:id/positions"): &Route{
						Handlers: map[string]http.Handler{
							http.MethodGet: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
								res.Write([]byte("/investments/:id/positions"))
							}),
						},
					},
					// Test normal route for subrouter
					Path("/investments/:investmentId/positions/:positionId"): &Route{
						Handlers: map[string]http.Handler{
							http.MethodGet: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
								res.Write([]byte("/investments/:investmentId/positions/:positionId"))
							}),
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		link       string
		statusCode int
		response   string
	}{
		{link: "http://example.com/investments", statusCode: http.StatusOK, response: "/investments"},
		{link: "http://example.com/investments/", statusCode: http.StatusPermanentRedirect, response: "/investments"},
		{link: "http://example.com/investments/some-id", statusCode: http.StatusOK, response: "/investments/:id"},
		{link: "http://example.com/investments/some-id/positions", statusCode: http.StatusOK, response: "/investments/:id/positions"},
		{link: "http://example.com/investments/some-id/positions/other-id", statusCode: http.StatusOK, response: "/investments/:investmentId/positions/:positionId"},
	}
	for _, tc := range testCases {
		t.Run(tc.link, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.link, nil)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)

			if want, got := tc.statusCode, res.Result().StatusCode; want != got {
				t.Errorf("expected status code to be %d, got %d instead", want, got)
			}
		})
	}
}
