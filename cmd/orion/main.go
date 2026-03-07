package main

import (
	"fmt"
	"log"
	"os"

	"orion/internal/runtime"
	"orion/internal/workspace"
	"orion/internal/tools"
	"orion/internal/cognition"
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
		log.Fatalf("Critical Failure: failed to initialize kernel: %v", err)
	}

	// Initialize Workspace Manager
	wsManager := workspace.NewManager(kernel.GetGlobalDB(), dataDir)
	err = wsManager.LoadWorkspaces(kernel.Context())
	if err != nil {
		log.Printf("Warning: failed to load existing workspaces: %v", err)
	}

	// Register tools
	kernel.GetToolRegistry().Register(tools.NewShellTool())
	kernel.GetToolRegistry().Register(tools.NewFileSystemTool())

	// Initialize Cognition Engine
	pipeline := &cognition.DefaultPipeline{}
	engine := cognition.NewEngine(pipeline)

	// Register agents
	conversationAgent := runtime.NewConversationAgent("agent-001", "ConversationAgent", kernel.GetEventBus())
	plannerAgent := runtime.NewPlannerAgent("agent-002", "PlannerAgent", kernel.GetEventBus(), engine)
	codeIndexerAgent := runtime.NewCodeIndexerAgent("agent-003", "CodeIndexerAgent", kernel.GetEventBus())
	analysisAgent := runtime.NewAnalysisAgent("agent-004", "AnalysisAgent", kernel.GetEventBus())
	memoryGardenerAgent := runtime.NewMemoryGardenerAgent("agent-005", "MemoryGardenerAgent", kernel.GetEventBus())
	patternDetectorAgent := runtime.NewPatternDetectorAgent("agent-006", "PatternDetectorAgent", kernel.GetEventBus())

	// Start agents
	agents := []runtime.Agent{
		conversationAgent,
		plannerAgent,
		codeIndexerAgent,
		analysisAgent,
		memoryGardenerAgent,
		patternDetectorAgent,
	}

	for _, agent := range agents {
		if err := agent.Start(kernel.Context()); err != nil {
			log.Printf("Warning: failed to start %s: %v", agent.Name(), err)
		}
	}

	// Ensure at least one workspace exists for initial use
	if len(wsManager.ListWorkspaces()) == 0 {
		fmt.Println("No workspaces found. Creating default workspace 'alpha'...")
		_, err := wsManager.CreateWorkspace(kernel.Context(), "alpha")
		if err != nil {
			log.Printf("Warning: failed to create default workspace: %v", err)
		}
	}

	// Start the Orion runtime via Lifecycle manager
	lifecycle := runtime.NewLifecycle(kernel)

	// UI is skipped in this environment due to missing X11 dependencies
	if os.Getenv("ORION_UI") == "true" {
		fmt.Println("Warning: UI skipped due to missing environment dependencies")
	}

	if err := lifecycle.StartOrion(); err != nil {
		log.Fatalf("Critical Failure: Orion encountered an error: %v", err)
	}
}
