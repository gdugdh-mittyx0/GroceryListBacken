package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"glbackend/internal/repo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormHandler struct {
	db *gorm.DB
}

func injectGormTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractGormTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 1000:
			pageSize = 1000
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func NewGormHandler(c *config) (*gormHandler, error) {
	ds := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Yakutsk",
		c.host,
		c.user,
		c.password,
		c.database,
		c.port,
	)
	db, err := gorm.Open(postgres.Open(ds), &gorm.Config{})
	if err != nil {
		return &gormHandler{}, err
	}
	// Проверка соединения с базой данных
	sqlDB, err := db.DB()
	if err != nil {
		return &gormHandler{}, fmt.Errorf("не удалось получить соединение с базой данных: %w", err)
	}

	// Пинг базы данных для проверки соединения
	if err := sqlDB.Ping(); err != nil {
		return &gormHandler{}, fmt.Errorf("не удалось выполнить ping базы данных: %w", err)
	} else {
		fmt.Println("ping db success")
	}

	return &gormHandler{db: db}, nil
}

func (g gormHandler) Create(ctx context.Context, table string, data interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	fmt.Println("IN CREATE")
	if err := db.WithContext(ctx).Table(table).Create(data).Error; err != nil {
		fmt.Println("ERROR", err)
		return err
	}
	return nil
}

func (g gormHandler) Update(ctx context.Context, table string, data interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	if err := db.WithContext(ctx).Table(table).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (g gormHandler) UpdateOne(ctx context.Context, table string, data, query interface{}, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	if err := db.WithContext(ctx).Table(table).Where(query, args...).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (g gormHandler) FindOne(ctx context.Context, table string, result interface{}, query interface{}, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(table).Where(query, args...).First(result).Error
	if err != nil {
		return err
	}
	return nil
}

func (g gormHandler) Delete(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(table).Where(condition, args...).Delete(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (g gormHandler) DeleteByQuery(ctx context.Context, table string, data, query interface{}, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(table).Where(query, args...).Delete(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (g gormHandler) Exec(ctx context.Context, query string, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Exec(query, args...).Error
}
func (g gormHandler) Preload(ctx context.Context, table string, query string, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(table).Preload(query, args...).Error
	if err != nil {
		return err
	}
	return nil
}

func (g gormHandler) RawQuery(ctx context.Context, table string, scanner interface{}, query string, args ...interface{}) error {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Table(table).Raw(query, args...).Scan(scanner).Error
	if err != nil {
		return err
	}
	return nil
}

func (g gormHandler) BeginFind(ctx context.Context, table string) repo.Find {
	db := g.db
	tx := extractGormTx(ctx)
	if tx != nil {
		db = tx
	}
	db = db.WithContext(ctx).Table(table)
	return &FindBuilder{
		db: db,
	}
}

type FindBuilder struct {
	db *gorm.DB
}

func (f *FindBuilder) OrderBy(query string) repo.Find {
	f.db = f.db.Order(query)
	return f
}

func (f *FindBuilder) Group(query string) repo.Find {
	f.db = f.db.Group(query)
	return f
}

func (f *FindBuilder) Limit(limit int) repo.Find {
	f.db = f.db.Limit(limit)
	return f
}
func (f *FindBuilder) Preload(query string, args ...interface{}) repo.Find {
	f.db = f.db.Preload(query, args...)
	return f
}

func (f *FindBuilder) Where(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Where(query, args...)
	return f
}

func (f *FindBuilder) Page(current, limit int) repo.Find {
	if limit > 0 {
		f.db = f.db.Scopes(paginate(current, limit))
	}
	return f
}

func (f *FindBuilder) Or(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Or(query, args...)
	return f
}

func (f *FindBuilder) Select(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Select(query, args...)
	return f
}

func (f *FindBuilder) Join(query string, args ...interface{}) repo.Find {
	f.db = f.db.Joins(query, args...)
	return f
}

func (f *FindBuilder) Having(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Having(query, args...)
	return f
}

func (f *FindBuilder) Scan(result interface{}) error {
	return f.db.Scan(result).Error
}

func (f *FindBuilder) Not(query interface{}, args ...interface{}) repo.Find {
	f.db = f.db.Not(query, args...)
	return f
}

func (f *FindBuilder) Count(total *int) (err error) {
	var count int64
	err = f.db.Count(&count).Error
	*total = int(count)
	return err
}

func (f *FindBuilder) Distinct(tables []string) repo.Find {
	f.db = f.db.Distinct(tables)
	return f
}

func (f *FindBuilder) Find(result interface{}, args ...interface{}) error {
	return f.db.Find(result, args...).Error
}

func (f *FindBuilder) Rows() (*sql.Rows, error) {
	return f.db.Rows()
}

func (f *FindBuilder) First(result interface{}, args ...interface{}) error {
	return f.db.First(result, args...).Error
}

func (f *FindBuilder) Clause(strength string) repo.Find {
	f.db = f.db.Clauses(clause.Locking{
		Strength: strength,
	})
	return f
}

func (g gormHandler) GetInstance() interface{} {
	return g.db
}

var ErrTransactionNotBegin = errors.New("TRANSACTION_NOT_BEGIN")

type gormTransaction struct {
	db       *gorm.DB
	isCommit bool
}

func NewGormTransaction(db repo.GSQL) repo.GTransaction {
	gorm := db.GetInstance().(*gorm.DB)
	return &gormTransaction{
		db:       gorm,
		isCommit: false,
	}
}

func (t *gormTransaction) Begin(ctx context.Context) context.Context {
	tx := t.db.Begin()
	return injectGormTx(ctx, tx)
}

func (t *gormTransaction) Commit(ctx context.Context) error {
	tx := extractGormTx(ctx)
	if tx == nil {
		return ErrTransactionNotBegin
	}
	t.isCommit = true
	return tx.Commit().Error
}

func (t *gormTransaction) Rollback(ctx context.Context) error {
	tx := extractGormTx(ctx)
	if tx == nil {
		return ErrTransactionNotBegin
	}
	return tx.Rollback().Error
}

func (t *gormTransaction) DeferRollback(ctx context.Context) {
	tx := extractGormTx(ctx)
	if tx == nil {
		return
	}
	if err := recover(); err != nil {
		tx.Rollback()
	} else if !t.isCommit {
		tx.Rollback()
	}
}
