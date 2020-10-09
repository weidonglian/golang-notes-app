package store

import (
	"errors"
	"fmt"
	"os"

	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type DB struct {
	conn   *sqlx.DB
	logger *logrus.Logger
}

func NewDB(logger *logrus.Logger, cfg config.Config) (DB*, error) {
	logger.Info("Connecting to database")
	conn, err := sqlx.Connect("postgres", cfg.Postgres.GetDataSourceName())
	if err != nil {
		return nil, err
	}

	if err := runMigrations(conn, conf.MigrationsPath, logger); err != nil {
		return nil, err
	}

	return &DB{
		conn: conn,
	}, nil
}

func runMigrations(conn *sqlx.DB, migrationsDir string, logger logging.Logger) error {
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		logger.Warnw("Migration directory not found, skipping migrations", "migrations_dir", migrationsDir)
		return nil
	}

	logger.Infow("Running migrations")
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Infow("No change when running migrations.")
			return nil
		}
		return err
	}
	logger.Infow("Migrations done.")
	return nil
}

// Close closes the database, freeing up any resources.
func (db *DB) Close() {
	db.conn.Close()
}
