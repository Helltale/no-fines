package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Helltale/no-fines/config"
	_ "github.com/lib/pq"
)

func ConnectToPostgre(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_NAME, cfg.POSTGRES_SSL_MODE)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("successfully connected to the database")
	return db, nil
}
