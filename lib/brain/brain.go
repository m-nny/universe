package brain

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type Brain struct {
	gormDb *gorm.DB
	sqlxDb *sqlx.DB
}

func New(gormDsn, sqlxDsn string, enableLogging bool) (*Brain, error) {
	gormDb, err := initGormDb(gormDsn, enableLogging)
	if err != nil {
		return nil, err
	}
	sqlxDb, err := sqlx.Connect("libsql", sqlxDsn)
	if err != nil {
		return nil, err
	}
	if err := initSqlx(sqlxDb); err != nil {
		return nil, err
	}
	return &Brain{gormDb: gormDb, sqlxDb: sqlxDb}, nil
}
