package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"orion/ent"
	"orion/internal/agents"
	"orion/internal/api"
	"orion/internal/cognition"
	"orion/internal/execution"
	"orion/internal/storage/sqlite"
	"orion/internal/types"
	"orion/internal/workspace"

	_ "github.com/mattn/go-sqlite3"
)

type Kernel struct {
	mu           sync.RWMutex
	DB           *sql.DB
	Ent          *ent.Client
	EventBus     *types.EventBus
	Scheduler    *execution.Scheduler
	WorkerPool   *execution.WorkerPool
	WorkspaceMgr *workspace.Manager
	Cognition    *cognition.Engine
	Supervisor   *agents.Supervisor
	Dispatcher   *agents.Dispatcher
	Router       *api.Router
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewKernel(d string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := sqlite.OpenGlobal(d)
	if err != nil {
		cancel()
		return nil, err
	}

	entClient, err := ent.Open("sqlite3", fmt.Sprintf("file:%s/orion.db?cache=shared&_fk=1", d))
	if err != nil {
		cancel()
		return nil, err
	}
	if err := entClient.Schema.Create(ctx); err != nil {
		cancel()
		return nil, err
	}

	eb := types.NewEventBus()
	wp := execution.NewWorkerPool(15)
	sch := execution.NewScheduler(wp)
	wm := workspace.NewManager(db, d)
	ar := agents.NewRegistry()
	as := agents.NewSupervisor(ar, eb)
	ad := agents.NewDispatcher(ar, sch, eb)
	ce := cognition.NewEngine(eb, ad)

	pers := sqlite.NewPersistence(entClient)
	ce.SetPersistence(pers)

	router := api.NewRouter()

	return &Kernel{
		DB:           db,
		Ent:          entClient,
		EventBus:     eb,
		Scheduler:    sch,
		WorkerPool:   wp,
		WorkspaceMgr: wm,
		Cognition:    ce,
		Supervisor:   as,
		Dispatcher:   ad,
		Router:       router,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

func (k *Kernel) Start() {
	k.WorkerPool.Start(k.ctx)
	k.Supervisor.StartAgents(k.ctx)

	go http.ListenAndServe(":8080", k.Router)

	fmt.Println("Kernel started.")
}

func (k *Kernel) Shutdown() {
	k.cancel()
	if k.Ent != nil {
		k.Ent.Close()
	}
	if k.DB != nil {
		k.DB.Close()
	}
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}
