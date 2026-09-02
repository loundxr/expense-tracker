package core_logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/loundxr/expense-tracker/internal/core/logger/handlers/slogpretty"
	"github.com/loundxr/expense-tracker/internal/core/logger/handlers/teehandler"
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
	var handlers []slog.Handler

	if cfg.Env == envLocal {
		handlers = append(handlers, setupPrettyHandler(os.Stdout))
	} else {
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	logDir := os.Getenv("LOCAL_LOG_DIR")
	if cfg.Env == envLocal && logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err == nil {
			timestamp := time.Now().Format("2006-01-02T15-04-05.000")
			fileName := filepath.Join(logDir, fmt.Sprintf("%s.log", timestamp))

			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})
				handlers = append(handlers, fileHandler)
			}
		}
	}

	combinedHandler := teehandler.New(handlers...)
	return slog.New(combinedHandler)
}

func setupPrettyHandler(out io.Writer) slog.Handler {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	return opts.NewPrettyHandler(out)
}
