package db

import (
	"github.com/weidonglian/golang-notes-app/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Session struct {
	conn   *sqlx.DB
	logger *logrus.Logger
}

func NewSession(logger *logrus.Logger, cfg config.Config) (*Session, error) {
	logger.Info("Connecting to database")
	conn, err := sqlx.Connect("postgres", cfg.Postgres.GetDataSourceName())
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(conn, cfg, logger); err != nil {
		return nil, err
	}

	return &Session{
		conn:   conn,
		logger: logger,
	}, nil
}

// Close closes the database, freeing up any resources.
func (sess Session) Close() {
	sess.conn.Close()
}
