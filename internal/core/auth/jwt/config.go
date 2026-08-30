package core_jwt

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type JWTConfig struct {
	Secret string        `envconfig:"SECRET" required:"true"`
	TTL    time.Duration `envconfig:"TTL" default:"24h"`
}

func NewConfig() (JWTConfig, error) {
	var cfg JWTConfig
	if err := envconfig.Process("JWT", &cfg); err != nil {
		return JWTConfig{}, err
	}
	return cfg, nil
}
