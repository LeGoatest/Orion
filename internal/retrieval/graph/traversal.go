package graph

import (
	"context"
	"fmt"
)

type Traversal struct {
	expander *Expander
}

func (t *Traversal) Traverse(ctx context.Context, startIDs []string) ([]string, error) {
	fmt.Println("Traversing memory graph...")
	return t.expander.Expand(ctx, startIDs, 2)
}
