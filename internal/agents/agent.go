package agents

import (
	"context"
)

type Agent interface {
	Name() string
	Type() string
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

type Registry struct {
	agents map[string]Agent
}

func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]Agent),
	}
}

func (r *Registry) Register(a Agent) {
	r.agents[a.Name()] = a
}

func (r *Registry) Get(name string) Agent {
	return r.agents[name]
}
