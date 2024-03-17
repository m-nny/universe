package brain

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Brain struct {
	gormDb *gorm.DB
}

func New(databasePath string) (*Brain, error) {
	gormDb, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := gormDb.AutoMigrate(&Artist{}, &Album{}, &Track{}); err != nil {
		return nil, err
	}
	return &Brain{gormDb: gormDb}, nil
}
