package auth

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	core_postgres_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool"
	repository "github.com/loundxr/expense-tracker/internal/features/auth/repository/postgres"
	"github.com/loundxr/expense-tracker/internal/features/auth/service"
	transport "github.com/loundxr/expense-tracker/internal/features/auth/transport/http"
)

type Dependencies struct {
	JWTConfig      core_jwt.JWTConfig
	Pool           core_postgres_pool.Pool
	Logger         *slog.Logger
	AccountCreator service.AccountCreator
}

type Module struct {
	handler *transport.AuthHTTPHandler
	svc     *service.UsersAuthService
	repo    *repository.UsersAuthRepository
}

func New(deps Dependencies) *Module {
	repo := repository.NewUsersAuthRepository(deps.Pool)
	svc := service.NewUsersAuthService(repo, deps.AccountCreator, deps.Pool, deps.Logger, deps.JWTConfig)
	h := transport.NewUsersAuthHandler(svc, deps.Logger)

	return &Module{
		handler: h,
		svc:     svc,
		repo:    repo,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}
