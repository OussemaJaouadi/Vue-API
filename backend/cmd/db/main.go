package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"vue-api/backend/internal/config"
	gormstorage "vue-api/backend/internal/storage/gorm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: db <generate|plan|migrate>")
		os.Exit(2)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	db, err := gormstorage.Open(cfg.Database)
	if err != nil {
		logger.Error("database open failed", "error", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		name := "schema"
		if len(os.Args) > 2 {
			name = os.Args[2]
		}
		plan, err := gormstorage.GenerateMigrationPlan(db, migrationsDir())
		if err != nil {
			logger.Error("database plan failed", "error", err)
			os.Exit(1)
		}
		path, err := gormstorage.GenerateMigrationFile(migrationsDir(), name, plan, time.Now())
		if err != nil {
			logger.Error("database migration generation failed", "error", err)
			os.Exit(1)
		}
		if path == "" {
			fmt.Println("Database schema is up to date. No migration generated.")
			return
		}
		fmt.Println("Generated migration:", path)
	case "plan":
		plan, err := gormstorage.PendingMigrationPlan(db, migrationsDir())
		if err != nil {
			logger.Error("database plan failed", "error", err)
			os.Exit(1)
		}
		if !plan.HasChanges() {
			schemaPlan, err := gormstorage.Plan(db)
			if err != nil {
				logger.Error("database plan failed", "error", err)
				os.Exit(1)
			}
			if schemaPlan.HasChanges() {
				fmt.Println("No generated migration files are pending. Run `task db:generate`.")
				fmt.Println(schemaPlan.String())
				return
			}
		}
		fmt.Println(plan.String())
	case "migrate":
		if err := gormstorage.ApplyMigrationFiles(db, migrationsDir()); err != nil {
			logger.Error("database migration failed", "error", err)
			os.Exit(1)
		}
		fmt.Println("Pending migrations applied.")
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func migrationsDir() string {
	return filepath.Join("migrations")
}
