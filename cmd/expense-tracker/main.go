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
	accounts_repository_postgres "github.com/loundxr/expense-tracker/internal/features/accounts/repository/postgres"
	accounts_service "github.com/loundxr/expense-tracker/internal/features/accounts/service"
	accounts_transport_http "github.com/loundxr/expense-tracker/internal/features/accounts/transport/http"
	analytics_repository_postgres "github.com/loundxr/expense-tracker/internal/features/analytics/repository/postgres"
	analytics_service "github.com/loundxr/expense-tracker/internal/features/analytics/service"
	analytics_transport_http "github.com/loundxr/expense-tracker/internal/features/analytics/transport/http"
	auth_repository_postgres "github.com/loundxr/expense-tracker/internal/features/auth/repository/postgres"
	auth_service "github.com/loundxr/expense-tracker/internal/features/auth/service"
	auth_transport_http "github.com/loundxr/expense-tracker/internal/features/auth/transport/http"
	categories_repository_postgres "github.com/loundxr/expense-tracker/internal/features/categories/repository/postgres"
	categories_service "github.com/loundxr/expense-tracker/internal/features/categories/service"
	categories_transport_http "github.com/loundxr/expense-tracker/internal/features/categories/transport/http"
	expenses_repository_postgres "github.com/loundxr/expense-tracker/internal/features/expenses/repository/postgres"
	expenses_service "github.com/loundxr/expense-tracker/internal/features/expenses/service"
	expenses_transport_http "github.com/loundxr/expense-tracker/internal/features/expenses/transport/http"
	users_repository_postgres "github.com/loundxr/expense-tracker/internal/features/users/repository/postgres"
	users_service "github.com/loundxr/expense-tracker/internal/features/users/service"
	users_transport_http "github.com/loundxr/expense-tracker/internal/features/users/transport/http"
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

	// users feature initialization
	logger.Debug("initializing feature users", slog.String("feature", "users"))
	usersRepo := users_repository_postgres.NewUsersRepository(pool)
	usersSvc := users_service.NewUsersService(usersRepo, logger)
	usersHandler := users_transport_http.NewUsersHTTPHandler(usersSvc, logger)

	// account feature initialization
	logger.Debug("initializing feature accounts", slog.String("feature", "accounts"))
	accountsRepo := accounts_repository_postgres.NewAccountsRepository(pool)
	accountsSvc := accounts_service.NewAccountsService(accountsRepo, usersRepo, logger)
	accountsHandler := accounts_transport_http.NewAccountsHTTPHandler(accountsSvc, logger)

	// auth feature initialization
	logger.Debug("initializing feature auth", slog.String("feature", "auth"))
	jwtCfg := utils.Must(core_jwt.NewConfig())
	authRepo := auth_repository_postgres.NewUsersAuthRepository(pool)
	authSvc := auth_service.NewUsersAuthService(authRepo, accountsSvc, logger, jwtCfg)
	authHandler := auth_transport_http.NewUsersAuthHandler(authSvc, logger)

	// categories feature initialization
	logger.Debug("initializing feature categories", slog.String("feature", "categories"))
	categoriesRepo := categories_repository_postgres.NewCategoriesRepository(pool)
	categoriesSvc := categories_service.NewCategoriesService(categoriesRepo, logger)
	categoriesHandler := categories_transport_http.NewCategoriesHTTPHandler(categoriesSvc, logger)

	// expenses feature initialization
	logger.Debug("initializing feature expenses", slog.String("feature", "expenses"))
	expensesRepo := expenses_repository_postgres.NewExpensesRepository(pool)
	expensesSvc := expenses_service.NewExpensesService(expensesRepo, accountsRepo, categoriesRepo, logger)
	expensesHandler := expenses_transport_http.NewExpensesHTTPHandler(expensesSvc, logger)

	//analytics feature initialization
	logger.Debug("initializing feature analytics", slog.String("feature", "analytics"))
	analyticsRepo := analytics_repository_postgres.NewAnalyticsRepository(pool)
	analyticsSvc := analytics_service.NewAnalyticsService(analyticsRepo, accountsRepo, logger)
	analyticsHandler := analytics_transport_http.NewStatisticsHTTPHandler(analyticsSvc, logger)

	// http server initialization
	logger.Debug("initializing HTTP server")
	httpConfig := utils.Must(core_http_server.NewConfig())
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	authHandler.RegisterRoutes(apiVersionRouter.Router())
	apiVersionRouter.Router().Group(func(r chi.Router) {
		r.Use(core_http_middleware.Auth(jwtCfg.Secret, logger))
		usersHandler.RegisterRoutes(r)
		accountsHandler.RegisterRoutes(r)
		categoriesHandler.RegisterRoutes(r)
		expensesHandler.RegisterRoutes(r)
		analyticsHandler.RegisterRoutes(r)
	})

	httpServer.RegisterAPIRouters(apiVersionRouter)
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("server.Run fatal error", slog.String("error", err.Error()))
	}
}
