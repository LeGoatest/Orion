package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("goal_id"),
		field.String("stage"),
		field.String("assigned_agent"),
		field.String("status"),
		field.Int("retry_count").Default(0),
		field.Time("created_at").Default(time.Now),
		field.Time("finished_at").Optional().Nillable(),
	}
}
