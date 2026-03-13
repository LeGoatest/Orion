package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("goal_id"),
		field.String("stage"),
		field.String("assigned_agent"),
		field.String("status").Default("PENDING"),
		field.Int("retry_count").Default(0),
		field.Time("created_at").Default(time.Now),
		field.Time("finished_at").Optional().Nillable(),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
