package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("could not load .env: %v", err)
	}
	dbName := os.Getenv("turso_db_name")
	authToken := os.Getenv("turso_db_token")
	if dbName == "" {
		log.Fatalf("turso_db_name is empty")
	}
	if authToken == "" {
		log.Fatalf("turso_db_token is empty")
	}
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, authToken)
	log.Printf("url: %s", url)

	db, err := sql.Open("libsql", url)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", url, err)
	}
	defer db.Close()
	if err := queryUsers(db); err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	log.Printf("Done.")
}

type User struct {
	ID   int
	Name string
}

func queryUsers(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		users = append(users, user)
		fmt.Println(user.ID, user.Name)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %v", err)
	}
	log.Printf("Found %d users", len(users))
	return nil
}
