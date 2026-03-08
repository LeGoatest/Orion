package embedding

import (
	"context"
)

// Job defines a unit of work for the embedding service
type Job interface {
	Execute(ctx context.Context, provider Provider) error
	ID() string
}

// TextJob handles simple text embeddings
type TextJob struct {
	id   string
	text string
}

func NewTextJob(id, text string) *TextJob {
	return &TextJob{id: id, text: text}
}

func (j *TextJob) ID() string { return j.id }

func (j *TextJob) Execute(ctx context.Context, provider Provider) error {
	_, err := provider.Generate(ctx, j.text)
	return err
}

// MemoryNodeJob handles node embedding and persistence
type MemoryNodeJob struct {
	NodeID string
	Text   string
	Store  interface {
		StoreEmbedding(ctx context.Context, id string, embedding []float32) error
	}
}

func (j *MemoryNodeJob) ID() string { return j.NodeID }

func (j *MemoryNodeJob) Execute(ctx context.Context, provider Provider) error {
	emb, err := provider.Generate(ctx, j.Text)
	if err != nil {
		return err
	}
	return j.Store.StoreEmbedding(ctx, j.NodeID, emb)
}
