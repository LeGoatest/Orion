package main

import (
	"context"
	"fmt"
	"log"
	"orion/ent"
	"orion/internal/runtime"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Initialize Global DB
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}
	client, err := ent.Open("sqlite3", "file:data/orion.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// 2. Initialize Kernel
	k := runtime.NewKernel(client, "data")

	// 3. Bootstrap Kernel
	if err := k.Bootstrap(); err != nil {
		log.Fatalf("failed to bootstrap kernel: %v", err)
	}

	fmt.Println("Orion Cognitive Runtime Kernel started.")
	k.Start()

	fmt.Println("Press Ctrl+C to shut down.")

	select {
	case <-k.Context().Done():
		fmt.Println("Kernel context cancelled.")
	case <-os.Interrupt:
		fmt.Println("Interrupt received.")
	}

	k.Shutdown()
	fmt.Println("Orion shutdown complete.")
}
