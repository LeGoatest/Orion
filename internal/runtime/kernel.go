package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"orion/internal/execution"
	"orion/internal/storage/sqlite"
	"orion/internal/tools"
)

// Kernel is the core runtime kernel for Orion
type Kernel struct {
	mu            sync.RWMutex
	globalDB      *sql.DB
	eventBus      *EventBus
	workerPool    *execution.WorkerPool
	scheduler     *execution.Scheduler
	toolRegistry  *tools.Registry
	dataDir       string
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewKernel creates a new Kernel instance
func NewKernel(dataDir string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize global database
	dbPath := filepath.Join(dataDir, "orion.db")
	db, err := sqlite.InitializeDB(dbPath, sqlite.GlobalDBSchema)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize global database: %w", err)
	}

	// Initialize subsystems
	eb := NewEventBus(db)
	wp := execution.NewWorkerPool(10, 100)
	sch := execution.NewScheduler(wp)
	tr := tools.NewRegistry()

	return &Kernel{
		globalDB:     db,
		eventBus:     eb,
		workerPool:   wp,
		scheduler:    sch,
		toolRegistry: tr,
		dataDir:      dataDir,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// Start boots the Orion cognitive runtime kernel
func (k *Kernel) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	fmt.Println("Orion Cognitive Runtime Kernel starting...")

	// Start worker pool
	k.workerPool.Start(k.ctx)
	fmt.Println("Worker pool started")

	// Subsystems initialization...
	fmt.Println("Event bus initialized")
	fmt.Println("Scheduler initialized")
	fmt.Println("Tool registry initialized")

	return nil
}

// Shutdown gracefully shuts down the kernel
func (k *Kernel) Shutdown() error {
	fmt.Println("Orion Cognitive Runtime Kernel shutting down...")
	k.cancel()
	k.workerPool.Stop()

	if k.globalDB != nil {
		k.globalDB.Close()
	}

	fmt.Println("Kernel shutdown complete")
	return nil
}

// GetGlobalDB returns the global database handle
func (k *Kernel) GetGlobalDB() *sql.DB {
	return k.globalDB
}

// GetEventBus returns the event bus instance
func (k *Kernel) GetEventBus() *EventBus {
	return k.eventBus
}

// GetScheduler returns the job scheduler
func (k *Kernel) GetScheduler() *execution.Scheduler {
	return k.scheduler
}

// GetToolRegistry returns the tool registry
func (k *Kernel) GetToolRegistry() *tools.Registry {
	return k.toolRegistry
}

// GetDataDir returns the base data directory
func (k *Kernel) GetDataDir() string {
	return k.dataDir
}

// Context returns the kernel context
func (k *Kernel) Context() context.Context {
	return k.ctx
}
