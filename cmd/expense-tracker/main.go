package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	core_cache_redis "github.com/loundxr/expense-tracker/internal/core/cache/redis"
	core_logger "github.com/loundxr/expense-tracker/internal/core/logger"
	core_pgx_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool/pgx"
	core_http_server "github.com/loundxr/expense-tracker/internal/core/transport/http/server"
	auth_postgres_repository "github.com/loundxr/expense-tracker/internal/features/auth/repository/postgres"
	auth_service "github.com/loundxr/expense-tracker/internal/features/auth/service"
	auth_transport_http "github.com/loundxr/expense-tracker/internal/features/auth/transport/http"
	"github.com/loundxr/expense-tracker/internal/utils"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	// logger initialization
	logCfg, err := core_logger.NewConfig()
	if err != nil {
		slog.Error("failed to load logger config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger := core_logger.Setup(logCfg)

	logger.Info("expense-tracker app", slog.String("env", logCfg.Env))

	// postgres pool connection initialization
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.New(ctx, utils.Must(core_pgx_pool.NewConfig()))
	if err != nil {
		logger.Error("failed to init postgres conn pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("postgres pool established")

	// redis connection initialization
	logger.Debug("initializing redis connection")
	cache, err := core_cache_redis.New(utils.Must(core_cache_redis.NewConfig()))
	if err != nil {
		logger.Error("failed to init redis connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cache.Close()
	logger.Info("redis connection established")

	// auth feature initialization
	logger.Debug("initializing feature auth", slog.String("feature", "auth"))
	authRepo := auth_postgres_repository.NewUsersAuthRepository(pool)
	authSvc := auth_service.NewUsersAuthService(authRepo, logger, utils.Must(core_jwt.NewConfig()))
	authHandler := auth_transport_http.NewUsersAuthHandler(authSvc, logger)

	// http server initialization
	logger.Debug("initializing HTTP server")
	httpConfig := utils.Must(core_http_server.NewConfig())
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	authHandler.RegisterRoutes(apiVersionRouter.Router())

	httpServer.RegisterAPIRouters(apiVersionRouter)
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("server.Run fatal error", slog.String("error", err.Error()))
	}
}
