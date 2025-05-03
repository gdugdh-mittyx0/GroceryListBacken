package router

import (
	"errors"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/config"
	"glbackend/internal/repo"
	"time"
)

type Server interface {
	Listen()
}

var errInvalidWebServerInstance = errors.New("invalid router server instance")

const (
	InstanceGin int = iota
)

func NewWebServerFactory(
	cfg config.Config,
	instance int,
	log logging.Logger,
	dbGSQL repo.GSQL,
	transaction repo.GTransaction,
	ctxTimeout time.Duration,
) (Server, error) {
	switch instance {
	case InstanceGin:
		return newGinServer(cfg, log, dbGSQL, transaction, ctxTimeout), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
