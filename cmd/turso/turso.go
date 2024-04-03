package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/m-nny/universe/lib/brain"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("could not load .env: %v", err)
	}
	dsn := os.Getenv("turso_db_dsn")
	authToken := os.Getenv("turso_db_token")
	if dsn == "" {
		log.Fatalf("turso_db_name is empty")
	}
	if authToken != "" {
		dsn += "authToken=" + authToken
	}
	log.Printf("url: %s", dsn)

	db, err := sqlx.Connect("libsql", dsn)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", dsn, err)
	}
	defer db.Close()

	if err := queryUsers(db); err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	log.Printf("Done.")
}

func queryUsers(db *sqlx.DB) error {
	var users []brain.User
	if err := db.Select(&users, "SELECT username, spotify_token_str FROM users"); err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("user: %+v", user)
	}
	log.Printf("Found %d users", len(users))
	return nil
}
