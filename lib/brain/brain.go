package brain

import (
	"gorm.io/gorm"
)

type Brain struct {
	gormDb *gorm.DB
}

func New(databasePath string, enableLogging bool) (*Brain, error) {
	gormDb, err := initGormDb(databasePath, enableLogging)
	if err != nil {
		return nil, err
	}
	return &Brain{gormDb: gormDb}, nil
}
