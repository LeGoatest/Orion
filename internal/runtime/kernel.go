package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	"orion/internal/agents"
	"orion/internal/cognition"
	"orion/internal/execution/scheduler"
	"orion/internal/execution/worker"
	"orion/internal/storage/sqlite"
	"orion/internal/types"
)

// Kernel is the central orchestrator of Orion
type Kernel struct {
	mu             sync.RWMutex
	globalDB       *sql.DB
	eventBus       *types.EventBus
	scheduler      *scheduler.Scheduler
	workerPool     *worker.WorkerPool
	workspaceMgr   *WorkspaceManager
	cognition      *cognition.CognitionEngine
	agentRegistry  *agents.Registry
	agentSupervisor *agents.Supervisor
	agentDispatcher *agents.Dispatcher
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewKernel initializes the runtime kernel with multi-agent orchestration
func NewKernel(dataDir string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 1. Storage
	dbPath := filepath.Join(dataDir, "orion.db")
	db, err := sqlite.InitializeGlobalDB(dbPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize global db: %w", err)
	}

	// 2. Core services
	eb := types.NewEventBus()
	wp := worker.NewWorkerPool(15) // Increased for multi-agent load
	sch := scheduler.NewScheduler(eb, wp)
	wm := NewWorkspaceManager(db, dataDir)

	// 3. Multi-Agent Systems
	ar := agents.NewRegistry()
	as := agents.NewSupervisor(ar, eb)
	ad := agents.NewDispatcher(ar, sch, eb)

	// 4. Cognition Engine (bridged to multi-agent dispatcher)
	ce := cognition.NewCognitionEngine(&cognition.DefaultPipeline{}, eb, ad)

	return &Kernel{
		globalDB:      db,
		eventBus:      eb,
		scheduler:     sch,
		workerPool:    wp,
		workspaceMgr:  wm,
		cognition:     ce,
		agentRegistry: ar,
		agentSupervisor: as,
		agentDispatcher: ad,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// Start boots the multi-agent Orion runtime
func (k *Kernel) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	fmt.Println("Orion Coordinated Runtime Kernel starting...")

	k.workerPool.Start(k.ctx)
	k.scheduler.Start(k.ctx)

	// Register Core System Agents
	k.agentRegistry.RegisterAgent(&ConversationAgent{})
	k.agentRegistry.RegisterAgent(&SymbolLookupAgent{})
	k.agentRegistry.RegisterAgent(&PlannerAgent{})
	k.agentRegistry.RegisterAgent(&RetrievalAgent{})
	k.agentRegistry.RegisterAgent(&MemoryGardenerAgent{})
	k.agentRegistry.RegisterAgent(&CodeIndexerAgent{})

	// Start Agents via Supervisor
	if err := k.agentSupervisor.StartAgents(k.ctx); err != nil {
		return fmt.Errorf("failed to start agents: %w", err)
	}

	fmt.Println("Cognition Engine (Dispatch-Ready): OK")
	fmt.Println("Multi-Agent Supervisor: RUNNING")

	return nil
}

// Shutdown coordinates system shutdown
func (k *Kernel) Shutdown() error {
	fmt.Println("Orion Cognitive Runtime Kernel shutting down...")
	k.agentSupervisor.StopAgents(k.ctx)
	k.cancel()
	if k.globalDB != nil { k.globalDB.Close() }
	return nil
}

func (k *Kernel) Context() context.Context { return k.ctx }
func (k *Kernel) GetEventBus() *types.EventBus { return k.eventBus }
func (k *Kernel) GetCognition() *cognition.CognitionEngine { return k.cognition }
func (k *Kernel) GetAgentRegistry() *agents.Registry { return k.agentRegistry }
