package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	core_cache_redis "github.com/loundxr/expense-tracker/internal/core/cache/redis"
	core_logger "github.com/loundxr/expense-tracker/internal/core/logger"
	core_pgx_pool "github.com/loundxr/expense-tracker/internal/core/repository/postgres/pool/pgx"
	"github.com/loundxr/expense-tracker/internal/utils"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	logCfg, err := core_logger.NewConfig()
	if err != nil {
		slog.Error("failed to load logger config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger := core_logger.Setup(logCfg)

	logger.Info("expense-tracker app", slog.String("env", logCfg.Env))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.New(ctx, utils.Must(core_pgx_pool.NewConfig()))
	if err != nil {
		logger.Error("failed to init postgres conn pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()
	logger.Debug("postgres pool established")

	logger.Debug("initializing redis connection")
	cache, err := core_cache_redis.New(utils.Must(core_cache_redis.NewConfig()))
	if err != nil {
		logger.Error("failed to init redis connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cache.Close()
	logger.Debug("redis connection established")
}
