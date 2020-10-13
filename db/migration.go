package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/config"
)

func RunMigrations(conn *sqlx.DB, cfg config.Config, logger *logrus.Logger) error {
	migrationsDir := cfg.MigrationsPath
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		workDir, _ := os.Getwd()
		logger.Errorf("Skipping migration, could not find migration path:%s with current working directory:%s", migrationsDir, workDir)
		return err
	}
	absMigrationDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	logger.Info("Running migrations")

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absMigrationDir),
		cfg.Postgres.DBName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No change when running migrations.")
			return nil
		}
		return err
	}
	logger.Info("Migrations done.")
	return nil
}
