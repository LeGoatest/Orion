package gardener

import (
	"context"
	"fmt"
	"time"
)

type MemoryGardener struct {
	// Periodic logic
}

func NewMemoryGardener() *MemoryGardener {
	return &MemoryGardener{}
}

func (mg *MemoryGardener) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mg.PerformGardening(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (mg *MemoryGardener) PerformGardening(ctx context.Context) {
	fmt.Println("Gardening: deduplication, consolidation, archiving.")
	// Deduplication
	// Consolidation
	// Archiving
	// Repair links
}
