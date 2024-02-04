package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Album holds the schema definition for the Album entity.
type Album struct {
	ent.Schema
}

// Fields of the Album.
func (Album) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty(),
		field.String("name").NotEmpty(),
		field.JSON("artistNames", []string{}),
		field.JSON("artistIds", []string{}),
	}
}

// Edges of the Album.
func (Album) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tracks", Track.Type),
	}
}
