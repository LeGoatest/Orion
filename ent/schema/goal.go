package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Goal holds the schema definition for the Goal entity.
type Goal struct {
	ent.Schema
}

// Fields of the Goal.
func (Goal) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Text("description"),
		field.String("current_stage").Default("OBSERVE"),
		field.String("status").Default("PENDING"),
		field.String("assigned_agent").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Goal.
func (Goal) Edges() []ent.Edge {
	return nil
}
