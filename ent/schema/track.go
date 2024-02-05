package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Track holds the schema definition for the Track entity.
type Track struct {
	ent.Schema
}

// Fields of the Track.
func (Track) Fields() []ent.Field {
	return []ent.Field{
		field.Strings("spotifyIds"),
		field.String("name").NotEmpty(),
		field.Int("trackNumber").Positive(),
		field.String("simplifiedName").NotEmpty(),
	}
}

// Edges of the Track.
func (Track) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("savedBy", User.Type).Ref("savedTracks"),
		edge.From("album", Album.Type).Ref("tracks").Unique(),
		edge.From("artists", Artist.Type).Ref("tracks"),
	}
}
