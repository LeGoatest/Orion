package prompt

import (
	"context"
	"fmt"
	"orion/internal/retrieval"
	"orion/internal/types"
	"time"
)

// Builder assembles a ContextBundle into a deterministic Envelope
type Builder struct {
	eb *types.EventBus
}

func NewBuilder(eb *types.EventBus) *Builder {
	return &Builder{eb: eb}
}

// BuildContext converts retrieval output into sections of the prompt envelope
func (b *Builder) BuildContext(ctx context.Context, bundle *retrieval.ContextBundle) (*Envelope, error) {
	fmt.Println("Prompt Builder: Assembling envelope from ContextBundle")

	env := NewEnvelope()
	env.Sections[FactsSection] = b.formatList(bundle.Facts)
	env.Sections[PatternsSection] = b.formatList(bundle.Patterns)
	env.Sections[InsightsSection] = b.formatList(bundle.Insights)
	env.Sections[CodeSection] = b.formatList(bundle.CodeSymbols)
	env.Sections[CallGraphSection] = b.formatList(bundle.CallGraphEdges)

	b.eb.Publish(types.Event{
		Type: "prompt.context_built",
		CreatedAt: time.Now(),
	})

	return env, nil
}

func (b *Builder) formatList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return "- " + items[0] // Simple placeholder
}
