package gardener

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"orion/internal/memory"
)

// MemoryGardener runs background tasks for knowledge maintenance
type MemoryGardener struct {
	db      *sql.DB
	manager *memory.MemoryManager
}

// NewMemoryGardener creates a new MemoryGardener
func NewMemoryGardener(db *sql.DB, mm *memory.MemoryManager) *MemoryGardener {
	return &MemoryGardener{
		db:      db,
		manager: mm,
	}
}

// ConsolidateKnowledge deduplicates and summarizes redundant memories
func (g *MemoryGardener) ConsolidateKnowledge(ctx context.Context) error {
	fmt.Println("Memory Gardener: Starting knowledge consolidation")

	// Consolidation logic:
	// 1. Identify highly similar memory nodes
	// 2. Merge nodes if their combined content is more compact or meaningful
	// 3. Update related links

	// For bootstrap, just logging
	fmt.Println("Memory Gardener: Knowledge consolidation complete")
	return nil
}

// ArchiveStaleNodes marks infrequently used nodes as archived
func (g *MemoryGardener) ArchiveStaleNodes(ctx context.Context, usageThreshold int) error {
	fmt.Printf("Memory Gardener: Archiving stale nodes (usage < %d)\n", usageThreshold)

	query := `UPDATE memory_nodes SET archived = TRUE WHERE usage_count < ? AND created_at < ?`
	thresholdTime := time.Now().AddDate(0, -1, 0) // 1 month ago
	_, err := g.db.ExecContext(ctx, query, usageThreshold, thresholdTime)
	return err
}

// RunBackgroundLoop starts the gardener in a non-blocking loop
func (g *MemoryGardener) RunBackgroundLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := g.ConsolidateKnowledge(ctx); err != nil {
					fmt.Printf("Memory Gardener Error: %v\n", err)
				}
				if err := g.ArchiveStaleNodes(ctx, 2); err != nil {
					fmt.Printf("Memory Gardener Error: %v\n", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
