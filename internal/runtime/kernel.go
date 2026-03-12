package runtime

import (
	"context"
	"fmt"
	"sync"
	"orion/ent"
	"orion/internal/execution"
	"orion/internal/workspace"
	"orion/internal/types"
	"orion/internal/agents"
	"orion/internal/cognition"
)

type Kernel struct {
	ctx      context.Context
	cancel   context.CancelFunc
	EventBus *types.EventBus
	Pool     *execution.WorkerPool
	Scheduler *execution.Scheduler
	DB       *ent.Client
	WorkspaceMgr *workspace.Manager
	Registry *agents.Registry
	Engine   *cognition.Engine
	mu       sync.Mutex
}

func NewKernel(db *ent.Client, dataDir string) *Kernel {
	ctx, cancel := context.WithCancel(context.Background())
	eb := types.NewEventBus()
	pool := execution.NewWorkerPool(10)
	sch := execution.NewScheduler(pool)
	reg := agents.NewRegistry()
	eng := cognition.NewEngine(db)

	return &Kernel{
		ctx:      ctx,
		cancel:   cancel,
		EventBus: eb,
		Pool:     pool,
		Scheduler: sch,
		DB:       db,
		WorkspaceMgr: workspace.NewManager(nil, dataDir),
		Registry: reg,
		Engine: eng,
	}
}

func (k *Kernel) Start() {
	// Startup logic
}

func (k *Kernel) Shutdown() {
	k.cancel()
	k.Pool.Shutdown()
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}

func (k *Kernel) Bootstrap() error {
	fmt.Println("Bootstrapping Orion Kernel...")
	// Register initial skeleton agents
	k.Registry.Register(agents.NewConversationAgent())
	k.Registry.Register(agents.NewPlannerAgent())
	k.Registry.Register(agents.NewAnalysisAgent())
	k.Registry.Register(agents.NewCodeIndexerAgent())
	k.Registry.Register(agents.NewMemoryGardenerAgent())
	k.Registry.Register(agents.NewPatternDetectorAgent())
	return nil
}
