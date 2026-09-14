package core_config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
}

func NewConfig() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %s: %w", tz, err)
	}

	return &Config{
		TimeZone: zone,
	}, nil
}
