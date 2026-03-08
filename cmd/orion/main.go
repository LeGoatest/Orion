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
	fmt.Println("# Orion Multi-Agent Cognitive Runtime #")
	fmt.Println("#########################################")

	// Determine data directory
	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	// 1. Initialize the Kernel
	kernel, err := runtime.NewKernel(dataDir)
	if err != nil {
		log.Fatalf("Critical Failure: failed to initialize kernel: %v", err)
	}

	// 2. Initialize and Register Autonomous Agents
	ar := kernel.GetAgentRegistry()

	ar.RegisterAgent(&runtime.ConversationAgent{})
	ar.RegisterAgent(&runtime.PlannerAgent{})
	ar.RegisterAgent(&runtime.CodeIndexerAgent{})
	ar.RegisterAgent(&runtime.AnalysisAgent{})
	ar.RegisterAgent(&runtime.MemoryGardenerAgent{})
	ar.RegisterAgent(&runtime.PatternDetectorAgent{})

	// 3. Start kernel
	if err := kernel.Start(); err != nil {
		log.Fatalf("Critical Failure: failed to start kernel: %v", err)
	}

	// 4. Simulate multi-agent task execution
	ctx := kernel.Context()
	testGoal := &goal.Goal{
		ID:          "goal-multi-agent-001",
		Description: "Perform repository analysis and update memory graph",
	}

	fmt.Printf("Submitting Multi-Agent Goal: %s\n", testGoal.Description)
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
