package brain

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Brain struct {
	gormDb *gorm.DB
}

func New(databasePath string) (*Brain, error) {
	gormDb, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{Logger: getLogger()})
	if err != nil {
		return nil, err
	}
	if err := gormDb.AutoMigrate(&Artist{}, &Album{}, &Track{}); err != nil {
		return nil, err
	}
	return &Brain{gormDb: gormDb}, nil
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
}
