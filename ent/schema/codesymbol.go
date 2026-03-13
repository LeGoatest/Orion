package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type CodeSymbol struct {
	ent.Schema
}

func (CodeSymbol) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("kind"),
		field.String("language"),
		field.String("path"),
		field.Int("line"),
	}
}
