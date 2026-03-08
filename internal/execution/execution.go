package execution

import "context"
import "sync"
import "fmt"

type Job interface { Execute(ctx context.Context) error; ID() string; Type() string }

type WorkerPool struct {
	maxWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
}

func NewWorkerPool(mw int) *WorkerPool { return &WorkerPool{maxWorkers: mw, jobQueue: make(chan Job, 100)} }

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.wg.Add(wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case job := <-wp.jobQueue:
					fmt.Printf("Job %s (%s) started\n", job.ID(), job.Type())
					if err := job.Execute(ctx); err != nil { fmt.Printf("Job failed: %v\n", err) }
				case <-ctx.Done(): return
				}
			}
		}()
	}
}

func (wp *WorkerPool) Submit(job Job) { wp.jobQueue <- job }

type Scheduler struct { pool *WorkerPool }
func NewScheduler(p *WorkerPool) *Scheduler { return &Scheduler{pool: p} }
func (s *Scheduler) Schedule(job Job) { s.pool.Submit(job) }
