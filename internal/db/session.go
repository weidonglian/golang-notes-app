package db

import (
	"github.com/weidonglian/notes-app/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Session interface {
	GetDB() *sqlx.DB
	Close() error
}

type sessionImpl struct {
	db     *sqlx.DB
	dbName string
	Logger *logrus.Logger
}

func (sess sessionImpl) GetDB() *sqlx.DB {
	return sess.db
}

// Close closes the database, freeing up any resources.
func (sess sessionImpl) Close() error {
	sess.Logger.Infof("close database session '%s'", sess.dbName)
	return sess.db.Close()
}

// NewSession will return a ready-to-use db session that already performs migration-up.
func NewSession(logger *logrus.Logger, cfg config.Config) (Session, error) {
	db, err := NewPostgresConnection(logger, cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db, cfg, logger); err != nil {
		return nil, err
	}

	return &sessionImpl{
		db:     db,
		dbName: cfg.Postgres.DBName,
		Logger: logger,
	}, nil
}

func NewPostgresConnection(logger *logrus.Logger, cfgPostgres config.PostgresConfig) (*sqlx.DB, error) {
	logger.Infof("Connecting to database: %s", cfgPostgres.DBName)
	if db, err := sqlx.Connect("postgres", cfgPostgres.GetDataSourceName()); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
