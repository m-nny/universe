package brain

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Brain struct {
	sqlxDb *sqlx.DB
}

func New(sqlxDsn string, enableLogging bool) (*Brain, error) {
	sqlxDb, err := sqlx.Connect("libsql", sqlxDsn)
	if err != nil {
		return nil, err
	}
	if err := initSqlx(sqlxDb); err != nil {
		return nil, err
	}
	return &Brain{sqlxDb: sqlxDb}, nil
}
