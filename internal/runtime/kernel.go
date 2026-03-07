package runtime

import (
	"context"
	"fmt"
)

// Kernel is the central orchestrator of Orion
type Kernel struct {
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewKernel(dataDir string) (*Kernel, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Kernel{ctx: ctx, cancel: cancel}, nil
}

func (k *Kernel) Context() context.Context {
	return k.ctx
}

func (k *Kernel) Start() error {
	fmt.Println("Kernel starting...")
	return nil
}

func (k *Kernel) Shutdown() error {
	fmt.Println("Kernel shutting down...")
	k.cancel()
	return nil
}
