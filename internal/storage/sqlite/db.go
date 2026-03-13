package sqlite

import (
	"orion/ent"
	"context"
	_ "github.com/mattn/go-sqlite3"
)

func InitGlobal(ctx context.Context, client *ent.Client) error {
	return client.Schema.Create(ctx)
}
