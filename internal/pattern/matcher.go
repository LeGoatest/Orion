package pattern

import (
	"context"
)

type Matcher struct {
	store *Store
}

func (m *Matcher) Match(ctx context.Context, goal string) (*Pattern, bool) {
	return m.store.MatchTrigger(ctx, goal)
}
