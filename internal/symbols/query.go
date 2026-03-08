package symbols

import (
	"context"
)

type Query struct {
	Store *Store
}

func (q *Query) ResolveSymbolReference(ctx context.Context, query string) ([]Symbol, error) {
	// Exact match -> fuzzy match -> fallback
	return q.Store.FindByName(ctx, query)
}
