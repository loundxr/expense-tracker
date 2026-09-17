package accounts

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/accounts/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/accounts/service"
	transport "github.com/loundxr/expense-tracker/internal/features/accounts/transport/http"
)

type Dependencies struct {
	Pool         core_postgres_pool.Pool
	Logger       *slog.Logger
	UserProvider service.UserProvider
}

type Module struct {
	handler *transport.AccountsHTTPHandler
	svc     *service.AccountsService
	repo    *repository.AccountsRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewAccountsRepository(deps.Pool)
	svc := service.NewAccountsService(repo, deps.UserProvider, deps.Logger)
	h := transport.NewAccountsHTTPHandler(svc, deps.Logger)

	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}

func (m *Module) Service() *service.AccountsService {
	return m.svc
}

func (m *Module) Repository() *repository.AccountsRepository {
	return m.repo
}
