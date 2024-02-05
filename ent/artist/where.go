// Code generated by ent, DO NOT EDIT.

package artist

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-nny/universe/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Artist {
	return predicate.Artist(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Artist {
	return predicate.Artist(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Artist {
	return predicate.Artist(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Artist {
	return predicate.Artist(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Artist {
	return predicate.Artist(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Artist {
	return predicate.Artist(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Artist {
	return predicate.Artist(sql.FieldLTE(FieldID, id))
}

// SpotifyId applies equality check predicate on the "spotifyId" field. It's identical to SpotifyIdEQ.
func SpotifyId(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldSpotifyId, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldName, v))
}

// SpotifyIdEQ applies the EQ predicate on the "spotifyId" field.
func SpotifyIdEQ(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldSpotifyId, v))
}

// SpotifyIdNEQ applies the NEQ predicate on the "spotifyId" field.
func SpotifyIdNEQ(v string) predicate.Artist {
	return predicate.Artist(sql.FieldNEQ(FieldSpotifyId, v))
}

// SpotifyIdIn applies the In predicate on the "spotifyId" field.
func SpotifyIdIn(vs ...string) predicate.Artist {
	return predicate.Artist(sql.FieldIn(FieldSpotifyId, vs...))
}

// SpotifyIdNotIn applies the NotIn predicate on the "spotifyId" field.
func SpotifyIdNotIn(vs ...string) predicate.Artist {
	return predicate.Artist(sql.FieldNotIn(FieldSpotifyId, vs...))
}

// SpotifyIdGT applies the GT predicate on the "spotifyId" field.
func SpotifyIdGT(v string) predicate.Artist {
	return predicate.Artist(sql.FieldGT(FieldSpotifyId, v))
}

// SpotifyIdGTE applies the GTE predicate on the "spotifyId" field.
func SpotifyIdGTE(v string) predicate.Artist {
	return predicate.Artist(sql.FieldGTE(FieldSpotifyId, v))
}

// SpotifyIdLT applies the LT predicate on the "spotifyId" field.
func SpotifyIdLT(v string) predicate.Artist {
	return predicate.Artist(sql.FieldLT(FieldSpotifyId, v))
}

// SpotifyIdLTE applies the LTE predicate on the "spotifyId" field.
func SpotifyIdLTE(v string) predicate.Artist {
	return predicate.Artist(sql.FieldLTE(FieldSpotifyId, v))
}

// SpotifyIdContains applies the Contains predicate on the "spotifyId" field.
func SpotifyIdContains(v string) predicate.Artist {
	return predicate.Artist(sql.FieldContains(FieldSpotifyId, v))
}

// SpotifyIdHasPrefix applies the HasPrefix predicate on the "spotifyId" field.
func SpotifyIdHasPrefix(v string) predicate.Artist {
	return predicate.Artist(sql.FieldHasPrefix(FieldSpotifyId, v))
}

// SpotifyIdHasSuffix applies the HasSuffix predicate on the "spotifyId" field.
func SpotifyIdHasSuffix(v string) predicate.Artist {
	return predicate.Artist(sql.FieldHasSuffix(FieldSpotifyId, v))
}

// SpotifyIdEqualFold applies the EqualFold predicate on the "spotifyId" field.
func SpotifyIdEqualFold(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEqualFold(FieldSpotifyId, v))
}

// SpotifyIdContainsFold applies the ContainsFold predicate on the "spotifyId" field.
func SpotifyIdContainsFold(v string) predicate.Artist {
	return predicate.Artist(sql.FieldContainsFold(FieldSpotifyId, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Artist {
	return predicate.Artist(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Artist {
	return predicate.Artist(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Artist {
	return predicate.Artist(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Artist {
	return predicate.Artist(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Artist {
	return predicate.Artist(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Artist {
	return predicate.Artist(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Artist {
	return predicate.Artist(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Artist {
	return predicate.Artist(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Artist {
	return predicate.Artist(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Artist {
	return predicate.Artist(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Artist {
	return predicate.Artist(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Artist {
	return predicate.Artist(sql.FieldContainsFold(FieldName, v))
}

// HasTracks applies the HasEdge predicate on the "tracks" edge.
func HasTracks() predicate.Artist {
	return predicate.Artist(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, TracksTable, TracksPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTracksWith applies the HasEdge predicate on the "tracks" edge with a given conditions (other predicates).
func HasTracksWith(preds ...predicate.Track) predicate.Artist {
	return predicate.Artist(func(s *sql.Selector) {
		step := newTracksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAlbums applies the HasEdge predicate on the "albums" edge.
func HasAlbums() predicate.Artist {
	return predicate.Artist(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, AlbumsTable, AlbumsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAlbumsWith applies the HasEdge predicate on the "albums" edge with a given conditions (other predicates).
func HasAlbumsWith(preds ...predicate.Album) predicate.Artist {
	return predicate.Artist(func(s *sql.Selector) {
		step := newAlbumsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Artist) predicate.Artist {
	return predicate.Artist(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Artist) predicate.Artist {
	return predicate.Artist(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Artist) predicate.Artist {
	return predicate.Artist(sql.NotPredicates(p))
}
