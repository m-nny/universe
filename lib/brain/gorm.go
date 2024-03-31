package brain

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	gormBatchSize = 100
)

func initGormDb(databasePath string, enableLogging bool) (*gorm.DB, error) {
	gormDb, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	allModels := []any{&Artist{}, &SpotifyAlbum{}, &MetaAlbum{}, &SpotifyTrack{}, &MetaTrack{}, &User{}, &DiscogsRelease{}, &DiscogsSeller{}}
	if err := gormDb.AutoMigrate(allModels...); err != nil {
		return nil, err
	}
	if enableLogging {
		gormDb.Logger = getLogger()
	}
	return gormDb, nil
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
