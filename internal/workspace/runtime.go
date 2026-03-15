package workspace

import (
	"context"
	"fmt"
	"orion/ent"
	"orion/ent/migrate"
	_ "github.com/mattn/go-sqlite3"
)

type WorkspaceRuntime struct {
	ID     string
	Path   string
	DB     *ent.Client
}

func NewWorkspaceRuntime(id, path string) (*WorkspaceRuntime, error) {
	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?_fk=1", path))
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return &WorkspaceRuntime{
		ID:   id,
		Path: path,
		DB:   client,
	}, nil
}

func (wr *WorkspaceRuntime) Close() error {
	return wr.DB.Close()
}
