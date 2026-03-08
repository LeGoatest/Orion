package main

import (
	"fmt"
	"log"
	"os"

	"orion/internal/runtime"
	"orion/internal/runtime/goal"
	"orion/internal/bridge"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {
	fmt.Println("#########################################")
	fmt.Println("# Orion Multi-Agent Cognitive Runtime #")
	fmt.Println("#########################################")

	dataDir := os.Getenv("ORION_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	// 1. Initialize Kernel
	kernel, err := runtime.NewKernel(dataDir)
	if err != nil {
		log.Fatalf("Critical Failure: failed to initialize kernel: %v", err)
	}

	// 2. Start Kernel
	if err := kernel.Start(); err != nil {
		log.Fatalf("Critical Failure: failed to start kernel: %v", err)
	}

	// 3. Setup Wails Bridge
	appBridge := bridge.NewBridge(kernel)

	// 4. Run Desktop Shell
	err = wails.Run(&options.App{
		Title:  "Orion Cognitive Runtime",
		Width:  1024,
		Height: 768,
		Bind: []interface{}{
			appBridge,
		},
		OnStartup: appBridge.Startup,
	})

	if err != nil {
		log.Fatalf("Desktop Shell Error: %v", err)
	}
}
