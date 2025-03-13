package config

import (
	"time"
)

func (c *Config) GetInitialDelay() time.Duration {
	return time.Duration(c.DB_CONNECTION_INITIAL_DELAY) * time.Second
}

func (c *Config) GetMaxDelay() time.Duration {
	return time.Duration(c.DB_CONNECTION_MAX_DELAY) * time.Second
}

func (c *Config) GetMultiplier() float64 {
	return float64(c.DB_CONNECTION_MULTIPLIER)
}

func (c *Config) GetMaxAttempts() int {
	return c.DB_CONNECTION_MAX_ATTEMPTS
}
