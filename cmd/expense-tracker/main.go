package main

import (
	"fmt"
	"log/slog"
	"os"

	core_logger "github.com/loundxr/expense-tracker/internal/core/logger"
)

func main() {
	logCfg, err := core_logger.NewConfig()
	if err != nil {
		fmt.Printf("failed to load logger config: %v", err)
		os.Exit(1)
	}
	logger := core_logger.Setup(logCfg)

	logger.Info("expense-tracker app", slog.String("env", logCfg.Env))
}
