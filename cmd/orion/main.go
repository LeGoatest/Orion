package main

import (
	"fmt"
	"log"
	"os"

	"orion/internal/runtime"
)

func main() {
	fmt.Println("#########################################")
	fmt.Println("# Orion Coordinated Cognitive Runtime #")
	fmt.Println("#########################################")

	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	// 1. Initialize the Kernel (wires up dispatcher, event bus, storage)
	kernel, err := runtime.NewKernel(dataDir)
	if err != nil {
		log.Fatalf("Critical Failure: failed to initialize kernel: %v", err)
	}

	// 2. Register Agents and Stage Ownership
	ar := kernel.Dispatcher.Registry
	ar.RegisterAgent(&runtime.ConversationAgent{})
	ar.RegisterAgent(&runtime.SymbolLookupAgent{})
	ar.RegisterAgent(&runtime.PlannerAgent{})
	ar.RegisterAgent(&runtime.RetrievalAgent{})
	ar.RegisterAgent(&runtime.MemoryGardenerAgent{})
	ar.RegisterAgent(&runtime.CodeIndexerAgent{})
	ar.RegisterAgent(&runtime.PatternDetectorAgent{})

	// 3. Start Multi-Agent System
	kernel.Start()

	// 4. Trigger Goal Injection (Bootstrap Simulation)
	ctx := kernel.Context()
	testGoalID := "goal-coordinated-001"
	testDesc := "Analyze the current multi-agent coordination performance"

	fmt.Printf("Submitting Test Goal: %s\n", testDesc)
	if err := kernel.Cognition.Process(ctx, testGoalID, testDesc); err != nil {
		fmt.Printf("Initial dispatch failed: %v\n", err)
	}

	fmt.Println("Orion Coordinated Runtime is running.")

	// Keep running to allow agent jobs to process
	select {
	case <-ctx.Done():
		fmt.Println("Runtime terminated.")
	}
}
