package db

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/config"
)

type SessionPool struct {
	parent          *Session
	cfg             config.Config
	Logger          *logrus.Logger
	forkedDatabases []string
}

func NewSessionPool(logger *logrus.Logger, cfg config.Config) SessionPool {
	sess, err := NewSession(logger, cfg)
	if err != nil {
		panic(err)
	}
	return SessionPool{
		parent: sess,
		Logger: logger,
		cfg:    cfg,
	}
}

func (i SessionPool) ForkNewSession() *Session {
	tempName := fmt.Sprintf("%s_%s", i.cfg.Postgres.DBName, xid.New().String())
	i.Logger.Infof("Fork child database: %s", tempName)
	i.parent.GetDB().MustExec(fmt.Sprintf("CREATE DATABASE %s", tempName))
	tempCfg := i.cfg
	tempCfg.Postgres.DBName = tempName
	if sess, err := NewSession(i.Logger, tempCfg); err != nil {
		panic(err)
	} else {
		i.forkedDatabases = append(i.forkedDatabases, tempName)
		return sess
	}
}

func (i SessionPool) Close() {
	i.Logger.Info("Close session pool database")

	for _, dbName := range i.forkedDatabases {
		i.Logger.Infof("Drop table: %s", dbName)
		i.parent.GetDB().MustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	}

	if err := i.parent.Close(); err != nil {
		panic(err)
	}
}
