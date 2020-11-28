package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
)

type forkedSession struct {
	sess              Session
	logger            *logrus.Logger
	forkedName        string
	parentPostgresCfg config.PostgresConfig
}

func (self forkedSession) GetDB() *sqlx.DB {
	return self.sess.GetDB()
}

// Close closes the database, freeing up any resources.
func (self forkedSession) Close() error {
	if err := self.sess.Close(); err != nil {
		return err
	}

	// we need to remove the database from postgres
	self.logger.Info("Remove the forked database from postgres")

	dbParent, err := NewPostgresConnection(self.logger, self.parentPostgresCfg)
	if err != nil {
		panic(err)
	}
	defer dbParent.Close()

	// drop the forked database from postgres
	self.logger.Infof("Drop database: %s", self.forkedName)
	dbParent.MustExec(fmt.Sprintf(`REVOKE CONNECT ON DATABASE %s FROM public`, self.forkedName))
	dbParent.MustExec(fmt.Sprintf(`SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s'`, self.forkedName))
	dbParent.MustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", self.forkedName))

	return nil
}

func NewForkedSession(logger *logrus.Logger, cfg config.Config) Session {
	dbParent, err := NewPostgresConnection(logger, cfg.Postgres)
	if err != nil {
		panic(err)
	}

	defer dbParent.Close()

	tempName := fmt.Sprintf("%s_%s", cfg.Postgres.DBName, xid.New().String())
	logger.Infof("Fork child database: %s", tempName)
	dbParent.MustExec(fmt.Sprintf("CREATE DATABASE %s", tempName))
	tempCfg := cfg
	tempCfg.Postgres.DBName = tempName
	if sess, err := NewSession(logger, tempCfg); err != nil {
		panic(err)
	} else {
		return forkedSession{
			sess:              sess,
			logger:            logger,
			forkedName:        tempName,
			parentPostgresCfg: cfg.Postgres,
		}
	}
}
