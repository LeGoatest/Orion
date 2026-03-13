package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"time"
)

type Goal struct {
	ent.Schema
}

func (Goal) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("description"),
		field.String("current_stage").Default("OBSERVE"),
		field.String("status").Default("pending"),
		field.String("assigned_agent").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Goal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("events", GoalEvent.Type),
	}
}
