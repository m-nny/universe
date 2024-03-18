package brain

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Brain struct {
	gormDb *gorm.DB
}

func New(databasePath string) (*Brain, error) {
	gormDb, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{Logger: logger.Default})
	if err != nil {
		return nil, err
	}
	if err := gormDb.AutoMigrate(&Artist{}, &Album{}, &Track{}); err != nil {
		return nil, err
	}
	gormDb.Logger = gormDb.Logger.LogMode(logger.Info)
	return &Brain{gormDb: gormDb}, nil
}
