package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"time"
)

type GoalEvent struct {
	ent.Schema
}

func (GoalEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.JSON("payload", map[string]interface{}{}),
		field.Time("created_at").Default(time.Now),
	}
}

func (GoalEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("goal", Goal.Type).Ref("events").Unique(),
	}
}
