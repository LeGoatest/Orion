package graph

import (
	"time"
)

type Node struct {
	ID         string
	Type       string
	Content    string
	Importance float64
	CreatedAt  time.Time
	Archived   bool
}

type Link struct {
	ID       string
	FromID   string
	ToID     string
	LinkType string
}

type MemoryGraph struct {
	// Storage integration
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{}
}

func (mg *MemoryGraph) AddNode(node Node) error {
	// Persist to DB
	return nil
}

func (mg *MemoryGraph) LinkNodes(from, to, linkType string) error {
	// Persist to DB
	return nil
}
