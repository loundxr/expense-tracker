package categories

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/categories/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/categories/service"
	transport "github.com/loundxr/expense-tracker/internal/features/categories/transport/http"
)

type Dependencies struct {
	Pool   core_postgres_pool.Pool
	Logger *slog.Logger
}

type Module struct {
	handler *transport.CategoriesHTTPHandler
	svc     *service.CategoriesService
	repo    *repository.CategoriesRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewCategoriesRepository(deps.Pool)
	svc := service.NewCategoriesService(repo, deps.Logger)
	h := transport.NewCategoriesHTTPHandler(svc, deps.Logger)

	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}

func (m *Module) Repository() *repository.CategoriesRepository {
	return m.repo
}
