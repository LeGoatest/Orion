package pattern

import (
	"context"
	"fmt"
)

// Matcher performs context-aware pattern matching
type Matcher struct {
	store *Store
}

// FindMatch finds a pattern that corresponds to the given context characteristics
func (m *Matcher) FindMatch(ctx context.Context, characteristics string) (*Pattern, bool) {
	fmt.Printf("Pattern Matcher: Searching for characteristics: %s\n", characteristics)
	return m.store.MatchTrigger(ctx, characteristics)
}
