package db

import (
	"fmt"
	"github.com/weidonglian/golang-notes-app/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	db     *sqlx.DB
	Logger *logrus.Logger
}

func (sess Session) GetDB() *sqlx.DB {
	return sess.db
}

// Close closes the database, freeing up any resources.
func (sess Session) Close() {
	if err := sess.db.Close(); err != nil {
		sess.Logger.Fatalf("Failed to close the db session: '%v'", err)
	}
}

func NewSession(logger *logrus.Logger, cfg config.Config) (*Session, error) {
	logger.Infof("Connecting to database '%s'", cfg.DatabaseDriver)

	var dataSourceName string
	if cfg.DatabaseDriver == config.DatabaseDriverSqlite3 {
		dataSourceName = cfg.Sqlite3.SourceName
	} else if cfg.DatabaseDriver == config.DatabaseDriverPostgres {
		dataSourceName = cfg.Postgres.GetDataSourceName()
	} else {
		panic(fmt.Sprintf("Unknown required DatabaseDriver:%s", cfg.DatabaseDriver))
	}

	var db *sqlx.DB
	if conn, err := sqlx.Connect(cfg.DatabaseDriver, dataSourceName); err != nil {
		return nil, err
	} else {
		db = conn
	}

	if err := RunMigrations(db, cfg, logger); err != nil {
		return nil, err
	}

	return &Session{
		db:     db,
		Logger: logger,
	}, nil
}
