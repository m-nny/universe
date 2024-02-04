package internal

import (
	"context"
	"fmt"

	"github.com/m-nny/universe/ent"
	_ "github.com/mattn/go-sqlite3"
)

func GetEntClient() (*ent.Client, error) {
	client, err := ent.Open("sqlite3", "file:data/ent.db?cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}
	return client, nil
}
