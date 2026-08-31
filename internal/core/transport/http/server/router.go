package core_http_server

import (
	"github.com/go-chi/chi/v5"
	core_http_middleware "github.com/loundxr/expense-tracker/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	router     chi.Router
	version    APIVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	v APIVersion,
	m ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		router:     chi.NewRouter(),
		version:    v,
		middleware: m,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, m := range r.middleware {
		r.router.Use(m)
	}

	for _, route := range routes {
		r.router.Method(route.Method, route.Path, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) Router() chi.Router {
	return r.router
}

func (r *APIVersionRouter) Group(fn func(chi.Router)) {
	r.router.Group(fn)
}
