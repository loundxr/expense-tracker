package users

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/users/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/users/service"
	transport "github.com/loundxr/expense-tracker/internal/features/users/transport/http"
)

type Dependencies struct {
	Pool   core_postgres_pool.Pool
	Logger *slog.Logger
}

type Module struct {
	handler *transport.UsersHTTPHandler
	svc     *service.UsersService
	repo    *repository.UsersRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewUsersRepository(deps.Pool)
	svc := service.NewUsersService(repo, deps.Logger)
	h := transport.NewUsersHTTPHandler(svc, deps.Logger)

	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}

func (m *Module) Repository() *repository.UsersRepository {
	return m.repo
}
