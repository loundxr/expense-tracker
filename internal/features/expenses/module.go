package expenses

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/expenses/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/expenses/service"
	transport "github.com/loundxr/expense-tracker/internal/features/expenses/transport/http"
)

type Dependencies struct {
	Pool            core_postgres_pool.Pool
	Logger          *slog.Logger
	AccountChecker  service.AccountAccessChecker
	CategoryChecker service.CategoryAccessChecker
}

type Module struct {
	handler *transport.ExpensesHTTPHandler
	svc     *service.ExpensesService
	repo    *repository.ExpensesRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewExpensesRepository(deps.Pool)
	svc := service.NewExpensesService(repo, deps.AccountChecker, deps.CategoryChecker, deps.Logger)
	h := transport.NewExpensesHTTPHandler(svc, deps.Logger)
	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}
