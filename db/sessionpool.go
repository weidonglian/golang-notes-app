package db

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
)

type SessionPool struct {
	parent          *Session
	cfg             config.Config
	Logger          *logrus.Logger
	forkedDatabases []string
}

var currentSessionPool *SessionPool

func LoadSessionPool(logger *logrus.Logger, cfg config.Config) *SessionPool {
	if currentSessionPool != nil {
		return currentSessionPool
	}

	sess, err := NewSession(logger, cfg)
	if err != nil {
		panic(err)
	}
	currentSessionPool = &SessionPool{
		parent: sess,
		Logger: logger,
		cfg:    cfg,
	}
	return currentSessionPool
}

func (i *SessionPool) ForkNewSession() *Session {
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

func UnloadSessionPool() {
	if currentSessionPool == nil {
		return
	}

	p := currentSessionPool

	p.Logger.Info("Close session pool database")

	db := p.parent.GetDB()
	for _, dbName := range p.forkedDatabases {
		p.Logger.Infof("Drop table: %s", dbName)
		db.MustExec(fmt.Sprintf(`REVOKE CONNECT ON DATABASE %s FROM public`, dbName))
		db.MustExec(fmt.Sprintf(`SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s'`, dbName))
		db.MustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	}

	if err := p.parent.Close(); err != nil {
		panic(err)
	}
	currentSessionPool = nil
}
