package main

import (
	"context"
	"fmt"
	"log"

	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/user"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := internal.GetEntClient()
	if err != nil {
		log.Fatalf("failed creating ent Client: %v", err)
	}
	defer client.Close()
	if _, err := CreateUser(context.Background(), client); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
