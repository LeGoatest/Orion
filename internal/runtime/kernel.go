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
	"orion/internal/pattern"
	"orion/internal/retrieval"
	"orion/internal/retrieval/graph"
	"orion/internal/storage/sqlite"
	"orion/internal/symbols"
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
	patternEngine  *pattern.Engine
	retrieval      *retrieval.RetrievalEngine
	agentRegistry  *agents.Registry
	agentSupervisor *agents.Supervisor
	agentDispatcher *agents.Dispatcher
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewKernel initializes the runtime kernel with all functional subsystems
func NewKernel(dataDir string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 1. Initialize global database
	dbPath := filepath.Join(dataDir, "orion.db")
	db, err := sqlite.InitializeGlobalDB(dbPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize global db: %w", err)
	}

	// 2. Core services
	eb := types.NewEventBus()
	wp := worker.NewWorkerPool(10)
	sch := scheduler.NewScheduler(eb, wp)
	wm := NewWorkspaceManager(db, dataDir)

	// 3. Multi-Agent Systems
	ar := agents.NewRegistry()
	as := agents.NewSupervisor(ar, eb)
	ad := agents.NewDispatcher(ar, sch, eb)

	// 4. Cognitive subsystems
	ps := pattern.NewStore(db)
	pe := pattern.NewEngine(ps, eb)

	ge := &graph.Expander{Db: db}
	sq := &symbols.Query{Store: &symbols.Store{Db: db}}
	re := retrieval.NewRetrievalEngine(eb, sq, ge)

	ce := cognition.NewCognitionEngine(&cognition.DefaultPipeline{}, eb, re, pe, ad)

	return &Kernel{
		globalDB:      db,
		eventBus:      eb,
		scheduler:     sch,
		workerPool:    wp,
		workspaceMgr:  wm,
		cognition:     ce,
		patternEngine: pe,
		retrieval:     re,
		agentRegistry: ar,
		agentSupervisor: as,
		agentDispatcher: ad,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// Start boots the Orion cognitive runtime
func (k *Kernel) Start() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	fmt.Println("Orion Cognitive Runtime Kernel starting...")

	k.workerPool.Start(k.ctx)
	k.scheduler.Start(k.ctx)

	// Start Agents
	if err := k.agentSupervisor.StartAgents(k.ctx); err != nil {
		return fmt.Errorf("failed to start agents: %w", err)
	}

	fmt.Println("Cognition Engine: READY")
	fmt.Println("Pattern Engine: READY")
	fmt.Println("Retrieval Engine: READY")
	fmt.Println("Multi-Agent Layer: READY")

	return nil
}

// Shutdown coordinates system shutdown
func (k *Kernel) Shutdown() error {
	fmt.Println("Orion Cognitive Runtime Kernel shutting down...")

	k.agentSupervisor.StopAgents(k.ctx)

	k.cancel()

	if k.globalDB != nil {
		k.globalDB.Close()
	}

	return nil
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}

func (k *Kernel) GetGlobalDB() *sql.DB {
	return k.globalDB
}

func (k *Kernel) GetEventBus() *types.EventBus {
	return k.eventBus
}

func (k *Kernel) GetCognition() *cognition.CognitionEngine {
	return k.cognition
}

func (k *Kernel) GetAgentRegistry() *agents.Registry {
	return k.agentRegistry
}
