package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router wraps a Go router implementation
type Router struct {
	*httprouter.Router
}

const (
	RequestParams = "params"
)

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.HandlerFunc(http.MethodGet, path, handler)
}

func GetParams(r *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(r.Context())
}
