package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"orion/internal/agents"
	"orion/internal/cognition"
	"orion/internal/execution"
	"orion/internal/storage/sqlite"
	"orion/internal/workspace"
	"path/filepath"
	"sync"
)

type Kernel struct {
	mu           sync.RWMutex
	DB           *sql.DB
	EventBus     *EventBus
	Scheduler    *execution.Scheduler
	WorkerPool   *execution.WorkerPool
	WorkspaceMgr *workspace.Manager
	Cognition    *cognition.Engine
	Supervisor   *agents.Supervisor
	Dispatcher   *agents.Dispatcher
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewKernel(d string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := sqlite.InitializeGlobalDB(filepath.Join(d, "orion.db"))
	if err != nil {
		cancel()
		return nil, err
	}

	eb := NewEventBus()
	wp := execution.NewWorkerPool(15)
	sch := execution.NewScheduler(wp)
	wm := workspace.NewManager(db, d)
	ar := agents.NewRegistry()
	as := agents.NewSupervisor(ar, eb)
	ad := agents.NewDispatcher(ar, sch, eb)
	ce := cognition.NewEngine(eb, ad)

	return &Kernel{
		DB:           db,
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

func (k *Kernel) Start() {
	k.WorkerPool.Start(k.ctx)
	k.Supervisor.StartAgents(k.ctx)
	fmt.Println("Kernel started.")
}

func (k *Kernel) Shutdown() {
	k.cancel()
	if k.DB != nil {
		k.DB.Close()
	}
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}
