package runtime

import (
	"context"
	"orion/ent"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestKernelCognitionLoop(t *testing.T) {
	os.MkdirAll("testdata", 0755)
	defer os.RemoveAll("testdata")

	client, err := ent.Open("sqlite3", "file:testdata/orion.db?cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}

	k := NewKernel(client, "testdata")
	if err := k.Bootstrap(); err != nil {
		t.Fatalf("failed to bootstrap: %v", err)
	}
	k.Start()
	defer k.Shutdown()

	err = k.Engine.Run(context.Background(), "Test intent")
	if err != nil {
		t.Errorf("Cognition loop failed: %v", err)
	}
}
