package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"orion/internal/agents"
	"orion/internal/cognition"
	"orion/internal/execution"
	"orion/internal/storage/sqlite"
	"orion/internal/types"
	"orion/internal/workspace"
)

// Kernel is the central orchestrator of Orion
type Kernel struct {
	mu             sync.RWMutex
	GlobalDB       *sql.DB
	EventBus       *types.EventBus
	Scheduler      *execution.Scheduler
	WorkerPool     *execution.WorkerPool
	WorkspaceMgr   *workspace.Manager
	Cognition      *cognition.CognitionEngine
	Supervisor     *agents.Supervisor
	Dispatcher     *agents.Dispatcher
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewKernel initializes the runtime kernel
func NewKernel(dataDir string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	dbPath := filepath.Join(dataDir, "orion.db")
	db, err := sqlite.InitializeGlobalDB(dbPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to init global db: %w", err)
	}

	eb := types.NewEventBus()
	wp := execution.NewWorkerPool(15)
	sch := execution.NewScheduler(wp)
	wm := workspace.NewManager(db, dataDir)
	ar := agents.NewRegistry()
	as := agents.NewSupervisor(ar, eb)
	ad := agents.NewDispatcher(ar, sch, eb)

	ce := cognition.NewCognitionEngine(&cognition.DefaultPipeline{}, eb, ad)

	return &Kernel{
		GlobalDB:     db,
		EventBus:     eb,
		Scheduler:    sch,
		WorkerPool:   wp,
		WorkspaceMgr: wm,
		Cognition:    ce,
		Supervisor:   as,
		Dispatcher:   ad,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// Start boots the Orion kernel
func (k *Kernel) Start() {
	k.WorkerPool.Start(k.ctx)
	k.Supervisor.StartAgents(k.ctx)
	fmt.Println("Orion Kernel started.")
}

func (k *Kernel) Shutdown() error {
	fmt.Println("Orion Kernel shutting down...")
	k.cancel()
	if k.GlobalDB != nil {
		k.GlobalDB.Close()
	}
	return nil
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}
