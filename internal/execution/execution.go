package execution

import ("context"; "fmt"; "sync")

type Job interface { Execute(ctx context.Context) error; ID() string; Type() string }
type WorkerPool struct { max int; queue chan Job; wg sync.WaitGroup }
func NewWorkerPool(m int) *WorkerPool { return &WorkerPool{max: m, queue: make(chan Job, 100)} }
func (wp *WorkerPool) Start(ctx context.Context) { wp.wg.Add(wp.max); for i := 0; i < wp.max; i++ { go func() { defer wp.wg.Done(); for { select { case j := <-wp.queue: fmt.Printf("Job %s starting\n", j.ID()); if err := j.Execute(ctx); err != nil { fmt.Printf("Job %s failed: %v\n", j.ID(), err) }; case <-ctx.Done(): return } } }() } }
func (wp *WorkerPool) Submit(j Job) { wp.queue <- j }
type Scheduler struct { pool *WorkerPool }
func NewScheduler(p *WorkerPool) *Scheduler { return &Scheduler{pool: p} }
func (s *Scheduler) Schedule(j Job) { s.pool.Submit(j) }
