package agents

import ("context"; "fmt"; "orion/internal/types"; "sync"; "time")

type Supervisor struct { Reg *Registry; EB *types.EventBus; mu sync.RWMutex; cancel context.CancelFunc }
func NewSupervisor(r *Registry, eb *types.EventBus) *Supervisor { return &Supervisor{Reg: r, EB: eb} }
func (s *Supervisor) StartAgents(ctx context.Context) error { for _, a := range s.Reg.List() { s.EB.Publish(types.Event{Type: "agent.started", Payload: a.Name(), CreatedAt: time.Now()}) }; return nil }
func (s *Supervisor) StopAgents(ctx context.Context) error { return nil }
