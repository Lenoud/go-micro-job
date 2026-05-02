package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

func (Department) Config() ent.Config {
	return ent.Config{
		Table: "b_department",
	}
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("title").
			Default(""),
		field.String("description").
			Default("").
			Optional(),
		field.Int("parent_id").
			Default(0).
			Optional(),
		field.Time("create_time").
			Default(time.Now).
			Optional(),
	}
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return nil
}
