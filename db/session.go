package db

import (
	"github.com/weidonglian/golang-notes-app/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
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
func (sess Session) Close() {
	sess.db.Close()
}

func NewSession(logger *logrus.Logger, cfg config.Config) (*Session, error) {
	logger.Info("Connecting to database")
	conn, err := sqlx.Connect("postgres", cfg.Postgres.GetDBDataSource())
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(conn, cfg, logger); err != nil {
		return nil, err
	}

	return &Session{
		db:     conn,
		Logger: logger,
	}, nil
}
