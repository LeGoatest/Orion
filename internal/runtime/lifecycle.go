package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Lifecycle manages the startup and shutdown of the Orion system
type Lifecycle struct {
	kernel *Kernel
}

// NewLifecycle creates a new Lifecycle instance
func NewLifecycle(kernel *Kernel) *Lifecycle {
	return &Lifecycle{
		kernel: kernel,
	}
}

// StartOrion initializes and starts the Orion runtime kernel
func (l *Lifecycle) StartOrion() error {
	fmt.Println("Orion Cognitive Runtime booting up...")

	// Start kernel
	if err := l.kernel.Start(); err != nil {
		return fmt.Errorf("failed to start kernel: %w", err)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Orion Cognitive Runtime is running. Press Ctrl+C to stop.")

	// Wait for signal
	select {
	case sig := <-sigChan:
		fmt.Printf("Shutdown signal received: %v\n", sig)
	case <-l.kernel.Context().Done():
		fmt.Println("Kernel context cancelled")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := l.ShutdownOrion(ctx); err != nil {
		return fmt.Errorf("failed to shutdown Orion gracefully: %w", err)
	}

	return nil
}

// ShutdownOrion stops the Orion runtime kernel
func (l *Lifecycle) ShutdownOrion(ctx context.Context) error {
	fmt.Println("Orion Cognitive Runtime shutting down gracefully...")
	return l.kernel.Shutdown()
}
