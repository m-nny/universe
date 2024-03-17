package brain

import "gorm.io/gorm"

type Brain struct {
	gormDb *gorm.DB
}

func New(gormDb *gorm.DB) *Brain {
	return &Brain{gormDb: gormDb}
}
