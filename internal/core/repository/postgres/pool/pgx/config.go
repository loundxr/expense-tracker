package core_pgx_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

func New() (Config, error) {
	var cfg Config

	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return cfg, nil
}

func ConfigMust(val Config, err error) Config {
	if err != nil {
		panic(err)
	}
	return val
}
