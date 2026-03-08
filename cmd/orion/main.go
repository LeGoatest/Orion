package main

import (
	"fmt"
	"log"
	"os"

	"orion/internal/runtime"
	"orion/internal/runtime/goal"
)

func main() {
	fmt.Println("#########################################")
	fmt.Println("# Orion Cognitive Runtime - Start #")
	fmt.Println("#########################################")

	// Determine data directory
	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	// 1. Initialize the Kernel (wires up all cognitive subsystems)
	kernel, err := runtime.NewKernel(dataDir)
	if err != nil {
		log.Fatalf("Critical Failure: failed to initialize kernel: %v", err)
	}

	// 2. Start kernel
	if err := kernel.Start(); err != nil {
		log.Fatalf("Critical Failure: failed to start kernel: %v", err)
	}

	// 3. Simulate a goal being processed
	ctx := kernel.Context()
	testGoal := &goal.Goal{
		ID:          "goal-001",
		Description: "Analyze the current repository for patterns",
	}

	fmt.Printf("Submitting Test Goal: %s\n", testGoal.Description)
	if err := kernel.GetCognition().Process(ctx, testGoal); err != nil {
		fmt.Printf("Goal processing failed: %v\n", err)
	}

	fmt.Println("Orion Cognitive Runtime is running.")

	// For simulation, we stay alive
	select {
	case <-ctx.Done():
		fmt.Println("Runtime terminated.")
	}
}
