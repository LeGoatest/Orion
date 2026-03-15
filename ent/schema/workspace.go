package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Workspace struct {
	ent.Schema
}

func (Workspace) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("path").Unique(),
	}
}
