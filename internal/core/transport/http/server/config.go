package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig process server cfg: %w", err)
	}
	return cfg, nil
}
