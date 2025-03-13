package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Helltale/no-fines/config"
)

func ConnectWithRetry[T any](cfg *config.Config, connectFunc func(*config.Config) (T, error)) (T, error) {
	var result T
	var err error

	// get data from cfg
	initialDelay := cfg.GetInitialDelay()
	multiplier := cfg.GetMultiplier()
	maxDelay := cfg.GetMaxDelay()
	maxAttempts := cfg.GetMaxAttempts()

	delay := initialDelay

	for attempts := 0; attempts < maxAttempts; attempts++ {
		result, err = connectFunc(cfg)
		if err == nil {
			return result, nil
		}

		slog.Warn("can not connect to postgres", "attempt", attempts+1, "error", err.Error())
		time.Sleep(delay)

		if delay < maxDelay {
			delay *= time.Duration(multiplier)
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return result, fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, err)
}
