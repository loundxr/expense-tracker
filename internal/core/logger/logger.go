package core_logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/loundxr/expense-tracker/internal/core/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env   string `envconfig:"ENV" default:"local"`
	Level string `envconfig:"LEVEL" default:"DEBUG"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return Config{}, fmt.Errorf("create logger config: %w", err)
	}
	return cfg, nil
}

type Logger struct {
}

func Setup(cfg Config) *slog.Logger {
	var log *slog.Logger

	localLogPath := os.Getenv("LOCAL_LOG_PATH")

	var output io.Writer = os.Stdout

	if cfg.Env == envLocal && localLogPath != "" {
		if err := os.MkdirAll(filepath.Dir(localLogPath), 0755); err == nil {
			file, err := os.OpenFile(localLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				output = io.MultiWriter(os.Stdout, file)
			}
		}
	}

	switch cfg.Env {
	case envLocal:
		log = setupPrettySlog(output)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog(out io.Writer) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(out)
	return slog.New(handler)
}
