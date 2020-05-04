package router

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

// Router wraps a Go router implementation
type Router struct {
	*httprouter.Router
}

// NewRouter returns a new wrapped router instance.
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

// Get adds a GET route to the router.
func (r *Router) Get(path string, handler http.Handler) {
	r.Handler(http.MethodGet, path, handler)
}

// Post adds a POST route to the router.
func (r *Router) Post(path string, handler http.Handler) {
	r.Handler(http.MethodPost, path, handler)
}

// GetParams returns parameters from the request context.
func GetParams(r *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(r.Context())
}

// GetURLParams returns parameters in the url.
func GetURLParams(r *http.Request) url.Values {
	return r.URL.Query()
}

// GetParamByName returns the value of the named parameter,
// and an error if it does not exist.
func GetParamByName(params *httprouter.Params, name string) (string, error) {
	value := params.ByName(name)

	if value == "" {
		return value, fmt.Errorf("parameter %s not found", name)
	}

	return value, nil
}

// GetURLParamByName returns the value of the named parameter,
// and an error if it does not exist.
func GetURLParamByName(params *url.Values, name string) (string, error) {
	value := params.Get(name)

	if value == "" {
		return value, fmt.Errorf("parameter %s not found", name)
	}

	return value, nil
}
