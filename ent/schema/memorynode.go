package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type MemoryNode struct {
	ent.Schema
}

func (MemoryNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.Text("content"),
		field.Float("importance").Default(0),
		field.Int("usage_count").Default(0),
		field.Bool("archived").Default(false),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
