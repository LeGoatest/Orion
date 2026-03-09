package main

import (
	"fmt"
	"orion/internal/runtime"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Orion Cognitive Runtime starting...")
	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	k, err := runtime.NewKernel(dataDir)
	if err != nil {
		fmt.Printf("Fatal: %v\n", err)
		os.Exit(1)
	}
	k.Supervisor.Reg.RegisterAgent(&runtime.ConversationAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.SymbolLookupAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.PlannerAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.RetrievalAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.ToolExecutionAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.AnalysisAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.CodeIndexerAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.MemoryGardenerAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Supervisor.Reg.RegisterAgent(&runtime.PatternDetectorAgent{BaseAgent: runtime.BaseAgent{EventBus: k.EventBus}})
	k.Start()
	fmt.Println("Runtime operational.")
	go func() { k.Cognition.Process(k.Context(), "goal-1", "Harden the coordination runtime") }()
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-s:
		k.Shutdown()
	case <-k.Context().Done():
	}
	fmt.Println("Shutdown complete.")
}
