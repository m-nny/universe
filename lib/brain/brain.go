package brain

import (
	"database/sql"

	"gorm.io/gorm"
)

type Brain struct {
	gormDb *gorm.DB
	turso  *sql.DB
}

func New(databasePath string, enableLogging bool) (*Brain, error) {
	gormDb, err := initGormDb(databasePath, enableLogging)
	if err != nil {
		return nil, err
	}
	turso, err := initTurso()
	if err != nil {
		return nil, err
	}
	return &Brain{gormDb: gormDb, turso: turso}, nil
}
