package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/config"
)

func RunMigrations(conn *sqlx.DB, cfg config.Config, logger *logrus.Logger) error {
	if _, err := os.Stat(cfg.MigrationsPath); os.IsNotExist(err) {
		workDir, _ := os.Getwd()
		logger.Errorf("Skipping migration, could not find migration path:%s with current working directory:%s", cfg.MigrationsPath, workDir)
		return err
	}

	var absMigrationDir string
	if mgdir, err := filepath.Abs(cfg.MigrationsPath); err != nil {
		return err
	} else {
		absMigrationDir = mgdir
	}

	logger.Infof("Running migrations from %s", absMigrationDir)

	var driver database.Driver
	if d, err := postgres.WithInstance(conn.DB, &postgres.Config{}); err != nil {
		return err
	} else {
		driver = d
	}

	var mg *migrate.Migrate
	if m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filepath.ToSlash(absMigrationDir)),
		cfg.Postgres.DBName, driver); err != nil {
		return err
	} else {
		mg = m
	}

	if err := mg.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No change when running migrations.")
			return nil
		}
		return err
	}
	logger.Info("Migrations done.")
	return nil
}
