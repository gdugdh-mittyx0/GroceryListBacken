package infrastructure

import (
	"errors"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/config"
	"glbackend/internal/infrastructure/database"
	"glbackend/internal/infrastructure/log"
	"glbackend/internal/infrastructure/router"
	"glbackend/internal/repo"
	"time"
)

type app struct {
	cfg         config.Config
	logger      logging.Logger
	dbGSQL      repo.GSQL
	transaction repo.GTransaction
	ctxTimeout  time.Duration
	webServer   router.Server
}

func NewConfig(config config.Config) *app {
	return &app{
		cfg: config,
	}
}

func (a *app) ContextTimeout(t time.Duration) *app {
	a.ctxTimeout = t
	return a
}

func (a *app) Logger(instance int) *app {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}
	a.logger = log
	a.logger.Infof("Success log configured")

	return a
}

func (a *app) DBGSql(instance int) *app {
	db, err := database.NewDatabaseSQLFactory(instance, a.cfg.DatabaseHost, a.cfg.DatabasePassword, a.cfg.DatabaseUser, a.cfg.DatabasePort, a.cfg.DatabaseDB)
	if err != nil {
		a.logger.Fatalln(err, "Could not make a connection database")
	}

	a.logger.Infof("Success connected to database")
	a.dbGSQL = db
	return a
}

func (a *app) GTransaction(instance int) *app {
	if a.dbGSQL == nil {
		a.logger.Fatalln(errors.New("Not db transaction"), "Could not make transaction instance")
	}
	transaction, err := database.NewTransactionFactory(instance, a.dbGSQL)
	if err != nil {
		a.logger.Fatalln(err, "Could not make transaction instance")
	}

	a.logger.Infof("Success transaction created")
	a.transaction = transaction
	return a
}

func (a *app) WebServer(instance int) *app {
	s, err := router.NewWebServerFactory(
		a.cfg,
		instance,
		a.logger,
		a.dbGSQL,
		a.transaction,
		a.ctxTimeout,
	)
	if err != nil {
		a.logger.Fatalln(err)
	}

	a.logger.Infof("Success router server configured")

	a.webServer = s
	return a
}

func (a *app) Start() {
	a.webServer.Listen()
}
