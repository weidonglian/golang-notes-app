package db

import (
	"github.com/weidonglian/golang-notes-app/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Session struct {
	db     *sqlx.DB
	Logger *logrus.Logger
}

func (sess Session) GetDB() *sqlx.DB {
	return sess.db
}

// Close closes the database, freeing up any resources.
func (sess Session) Close() error {
	return sess.db.Close()
}

// NewSession will return a ready-to-use db session that already performs migration-up.
func NewSession(logger *logrus.Logger, cfg config.Config) (*Session, error) {
	db, err := newPostgresConnection(logger, cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db, cfg, logger); err != nil {
		return nil, err
	}

	return &Session{
		db:     db,
		Logger: logger,
	}, nil
}

func newPostgresConnection(logger *logrus.Logger, cfgPostgres config.PostgresConfig) (*sqlx.DB, error) {
	logger.Infof("Connecting to database: %s", cfgPostgres.DBName)
	if db, err := sqlx.Connect("postgres", cfgPostgres.GetDataSourceName()); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
