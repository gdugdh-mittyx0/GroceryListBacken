package repo

import (
	"context"
	"database/sql"
	"glbackend/internal/entities"
)

type Repo interface {
	Grouper() entities.GrouperRepository
	Product() entities.ProductRepository
}

type repo struct {
	db GSQL
	entities.GrouperRepository
	entities.ProductRepository
}

func NewRepo(db GSQL) Repo {
	return &repo{
		db:                db,
		GrouperRepository: NewGrouperGSQL(db),
		ProductRepository: NewProductGSQL(db),
	}
}

func (r *repo) Grouper() entities.GrouperRepository {
	return r.GrouperRepository
}

func (r *repo) Product() entities.ProductRepository {
	return r.ProductRepository
}

type GSQL interface {
	Create(ctx context.Context, table string, data interface{}) error
	UpdateOne(ctx context.Context, table string, data, query interface{}, args ...interface{}) error
	FindOne(ctx context.Context, table string, result interface{}, query interface{}, args ...interface{}) error
	BeginFind(ctx context.Context, tableName string) Find
	Delete(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) error
	DeleteByQuery(ctx context.Context, table string, data, query interface{}, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
	RawQuery(ctx context.Context, table string, scanner interface{}, query string, args ...interface{}) error
	Preload(ctx context.Context, table string, query string, args ...interface{}) error
	GetInstance() interface{}
}

type Find interface {
	Where(query interface{}, args ...interface{}) Find
	Having(query interface{}, args ...interface{}) Find
	Page(current, limit int) Find
	Join(query string, args ...interface{}) Find
	Or(query interface{}, args ...interface{}) Find
	Not(query interface{}, args ...interface{}) Find
	Count(total *int) error
	Find(result interface{}, args ...interface{}) error
	First(result interface{}, args ...interface{}) error
	Select(query interface{}, args ...interface{}) Find
	Preload(query string, args ...interface{}) Find
	Scan(result interface{}) error
	OrderBy(query string) Find
	Group(query string) Find
	Limit(limit int) Find
	Clause(strength string) Find
	Rows() (*sql.Rows, error)
	Distinct(tables []string) Find
}

type GTransaction interface {
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	DeferRollback(ctx context.Context)
}
