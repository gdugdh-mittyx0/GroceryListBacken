package database

import (
	"errors"
	"glbackend/internal/repo"
)

var (
	errInvalidCacheInstance = errors.New("invalid cache instance")
)

func NewTransactionFactory(instance int, db repo.GSQL) (repo.GTransaction, error) {
	switch instance {
	case InstanceGPostgres:
		return NewGormTransaction(db), nil
	default:
		return nil, errInvalidCacheInstance
	}
}
