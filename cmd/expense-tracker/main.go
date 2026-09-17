package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	core_jwt "github.com/loundxr/expense-tracker/internal/core/auth/jwt"
	core_cache_redis "github.com/loundxr/expense-tracker/internal/core/cache/redis"
	core_config "github.com/loundxr/expense-tracker/internal/core/config"
	core_logger "github.com/loundxr/expense-tracker/internal/core/logger"
	core_pgx_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/loundxr/expense-tracker/internal/core/transport/http/middleware"
	core_http_server "github.com/loundxr/expense-tracker/internal/core/transport/http/server"
	"github.com/loundxr/expense-tracker/internal/features/accounts"
	"github.com/loundxr/expense-tracker/internal/features/analytics"
	"github.com/loundxr/expense-tracker/internal/features/auth"
	"github.com/loundxr/expense-tracker/internal/features/categories"
	"github.com/loundxr/expense-tracker/internal/features/expenses"
	"github.com/loundxr/expense-tracker/internal/features/users"
	"github.com/loundxr/expense-tracker/internal/utils"
)

func main() {
	cfg := utils.Must(core_config.NewConfig())
	time.Local = cfg.TimeZone

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	// logger
	logCfg, err := core_logger.NewConfig()
	if err != nil {
		slog.Error("failed to load logger config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger := core_logger.Setup(logCfg)

	logger.Info("expense-tracker app", slog.String("env", logCfg.Env))

	// postgres pool connection
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.New(ctx, utils.Must(core_pgx_pool.NewConfig()))
	if err != nil {
		logger.Error("failed to init postgres conn pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("postgres pool established")

	// redis connection
	logger.Debug("initializing redis connection")
	cache, err := core_cache_redis.New(utils.Must(core_cache_redis.NewConfig()))
	if err != nil {
		logger.Error("failed to init redis connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cache.Close()
	logger.Info("redis connection established")

	// users feature
	logger.Debug("initializing feature users", slog.String("feature", "users"))
	usersModule := users.New(users.Dependencies{
		Pool:   pool,
		Logger: logger,
	})

	// account feature
	logger.Debug("initializing feature accounts", slog.String("feature", "accounts"))
	accountsModule := accounts.New(accounts.Dependencies{
		Pool:         pool,
		Logger:       logger,
		UserProvider: usersModule.Repository(),
	})

	// auth feature
	logger.Debug("initializing feature auth", slog.String("feature", "auth"))
	jwtCfg := utils.Must(core_jwt.NewConfig())
	authModule := auth.New(auth.Dependencies{
		JWTConfig:      jwtCfg,
		Logger:         logger,
		Pool:           pool,
		AccountCreator: accountsModule.Service(),
	})

	// categories feature
	logger.Debug("initializing feature categories", slog.String("feature", "categories"))
	categoriesModule := categories.New(categories.Dependencies{
		Pool:   pool,
		Logger: logger,
	})

	// expenses feature
	logger.Debug("initializing feature expenses", slog.String("feature", "expenses"))
	expensesModule := expenses.New(expenses.Dependencies{
		Pool:            pool,
		Logger:          logger,
		AccountChecker:  accountsModule.Repository(),
		CategoryChecker: categoriesModule.Repository(),
	})

	//analytics feature
	logger.Debug("initializing feature analytics", slog.String("feature", "analytics"))
	analyticsModule := analytics.New(analytics.Dependencies{
		Pool:           pool,
		Logger:         logger,
		AccountChecker: accountsModule.Repository(),
	})

	// http server
	logger.Debug("initializing HTTP server")
	httpConfig := utils.Must(core_http_server.NewConfig())
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	authModule.RegisterRoutes(apiVersionRouter.Router())
	apiVersionRouter.Router().Group(func(r chi.Router) {
		r.Use(core_http_middleware.Auth(jwtCfg.Secret, logger))
		usersModule.RegisterRoutes(r)
		accountsModule.RegisterRoutes(r)
		categoriesModule.RegisterRoutes(r)
		expensesModule.RegisterRoutes(r)
		analyticsModule.RegisterRoutes(r)
	})

	httpServer.RegisterAPIRouters(apiVersionRouter)
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("server.Run fatal error", slog.String("error", err.Error()))
	}
}
