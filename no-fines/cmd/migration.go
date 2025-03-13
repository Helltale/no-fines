package cmd

import (
	"fmt"
	"log"

	"github.com/Helltale/no-fines/config"
	"github.com/Helltale/no-fines/internal/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "perform database migrations",
	Long:  "Perform database migrations using up/down SQL files.",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "apply all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("down")
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runMigration(direction string) {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	dbConn, err := db.ConnectWithRetry(cfg, db.ConnectToPostgre)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer dbConn.Close()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %v", err)
	}

	migrationPath := fmt.Sprintf("file://%s", cfg.POSTGRES_MIGRATION_PACKAGE)
	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to run migrations: %v", err)
		}
		log.Println("migrations applied successfully.")
	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to rollback migration: %v", err)
		}
		log.Println("migration rolled back successfully.")
	default:
		log.Fatalf("unknown migration direction: %s", direction)
	}
}
