package brain

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	gormBatchSize = 100
)

func initGormDb(databasePath string, enableLogging bool) (*gorm.DB, error) {
	dbName := os.Getenv("turso_db_name")
	authToken := os.Getenv("turso_db_token")
	if dbName == "" {
		log.Fatalf("turso_db_name is empty")
	}
	if authToken == "" {
		log.Fatalf("turso_db_token is empty")
	}
	dsn := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, authToken)
	log.Printf("Connecting to %s", dsn)
	gormDb, err := gorm.Open(sqlite.Dialector{
		DriverName: "libsql",
		DSN:        dsn,
	}, &gorm.Config{})
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
