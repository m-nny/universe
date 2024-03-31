package brain

import (
	"database/sql"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func initTurso() (*sql.DB, error) {
	url := "libsql://localhost:8080"
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
