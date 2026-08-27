package core_cache_redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"6379"`
	Password string        `envconfig:"PASSWORD" default:""`
	DB       int           `envconfig:"DATABASE" default:"0"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("REDIS", &cfg); err != nil {
		return Config{}, fmt.Errorf("redis envconfig process: %w", err)
	}
	return cfg, nil
}
