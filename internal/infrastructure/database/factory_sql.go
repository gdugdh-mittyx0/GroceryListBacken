package database

import (
	"errors"
	"glbackend/internal/repo"
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

const (
	InstanceGPostgres int = iota
)

func NewDatabaseSQLFactory(instance int, host string, password string, user string, port uint, database string) (repo.GSQL, error) {
	switch instance {
	case InstanceGPostgres:
		return NewGormHandler(newConfigGorm(host, database, password, user, port))
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
}
