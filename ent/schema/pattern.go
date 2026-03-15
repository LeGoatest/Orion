package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type Pattern struct {
	ent.Schema
}

func (Pattern) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Text("description"),
		field.JSON("metadata", map[string]interface{}{}),
		field.Time("created_at").Default(time.Now),
	}
}
