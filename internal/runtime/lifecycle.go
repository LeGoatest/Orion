package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type LifecycleManager struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewLifecycleManager() *LifecycleManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &LifecycleManager{ctx: ctx, cancel: cancel}
}
func (l *LifecycleManager) Wait() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-s:
		l.cancel()
	case <-l.ctx.Done():
	}
	fmt.Println("Shutting down...")
}
func (l *LifecycleManager) Context() context.Context { return l.ctx }
