package main

import (
	"fmt"
	"os"
	"orion/internal/runtime"
)

func main() {
	fmt.Println("#########################################")
	fmt.Println("# Orion Cognitive Runtime - Bootstrapping #")
	fmt.Println("#########################################")

	// Determine data directory
	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	// Initialize the Kernel
	kernel, err := runtime.NewKernel(dataDir)
	if err != nil {
		fmt.Printf("Critical Failure: failed to initialize kernel: %v\n", err)
		os.Exit(1)
	}

	// Start kernel
	if err := kernel.Start(); err != nil {
		fmt.Printf("Critical Failure: failed to start kernel: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Orion Cognitive Runtime is running.")

	// For bootstrap, keep running
	select {}
}
