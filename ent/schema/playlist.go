package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Playlist struct {
	ent.Schema
}

// Fields of the Playlist.
func (Playlist) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("snaphot_id").NotEmpty(),
	}
}

// Edges of the Playlist.
func (Playlist) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("playlists").Unique(),
	}
}
