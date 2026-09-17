package analytics

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/analytics/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/analytics/service"
	transport "github.com/loundxr/expense-tracker/internal/features/analytics/transport/http"
)

type Dependencies struct {
	Pool           core_postgres_pool.Pool
	Logger         *slog.Logger
	AccountChecker service.AccountChecker
}

type Module struct {
	handler *transport.AnalyticsHTTPHandler
	svc     *service.AnalyticsService
	repo    *repository.AnalyticsRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewAnalyticsRepository(deps.Pool)
	svc := service.NewAnalyticsService(repo, deps.AccountChecker, deps.Logger)
	h := transport.NewAnalyticsHTTPHandler(svc, deps.Logger)

	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}
