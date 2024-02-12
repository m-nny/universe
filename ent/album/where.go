// Code generated by ent, DO NOT EDIT.

package album

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-nny/universe/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Album {
	return predicate.Album(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Album {
	return predicate.Album(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Album {
	return predicate.Album(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Album {
	return predicate.Album(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Album {
	return predicate.Album(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Album {
	return predicate.Album(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Album {
	return predicate.Album(sql.FieldLTE(FieldID, id))
}

// DiscogsMasterId applies equality check predicate on the "discogsMasterId" field. It's identical to DiscogsMasterIdEQ.
func DiscogsMasterId(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldDiscogsMasterId, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldName, v))
}

// SimplifiedName applies equality check predicate on the "simplifiedName" field. It's identical to SimplifiedNameEQ.
func SimplifiedName(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldSimplifiedName, v))
}

// DiscogsMasterIdEQ applies the EQ predicate on the "discogsMasterId" field.
func DiscogsMasterIdEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdNEQ applies the NEQ predicate on the "discogsMasterId" field.
func DiscogsMasterIdNEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldNEQ(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdIn applies the In predicate on the "discogsMasterId" field.
func DiscogsMasterIdIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldIn(FieldDiscogsMasterId, vs...))
}

// DiscogsMasterIdNotIn applies the NotIn predicate on the "discogsMasterId" field.
func DiscogsMasterIdNotIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldNotIn(FieldDiscogsMasterId, vs...))
}

// DiscogsMasterIdGT applies the GT predicate on the "discogsMasterId" field.
func DiscogsMasterIdGT(v string) predicate.Album {
	return predicate.Album(sql.FieldGT(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdGTE applies the GTE predicate on the "discogsMasterId" field.
func DiscogsMasterIdGTE(v string) predicate.Album {
	return predicate.Album(sql.FieldGTE(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdLT applies the LT predicate on the "discogsMasterId" field.
func DiscogsMasterIdLT(v string) predicate.Album {
	return predicate.Album(sql.FieldLT(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdLTE applies the LTE predicate on the "discogsMasterId" field.
func DiscogsMasterIdLTE(v string) predicate.Album {
	return predicate.Album(sql.FieldLTE(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdContains applies the Contains predicate on the "discogsMasterId" field.
func DiscogsMasterIdContains(v string) predicate.Album {
	return predicate.Album(sql.FieldContains(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdHasPrefix applies the HasPrefix predicate on the "discogsMasterId" field.
func DiscogsMasterIdHasPrefix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasPrefix(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdHasSuffix applies the HasSuffix predicate on the "discogsMasterId" field.
func DiscogsMasterIdHasSuffix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasSuffix(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdIsNil applies the IsNil predicate on the "discogsMasterId" field.
func DiscogsMasterIdIsNil() predicate.Album {
	return predicate.Album(sql.FieldIsNull(FieldDiscogsMasterId))
}

// DiscogsMasterIdNotNil applies the NotNil predicate on the "discogsMasterId" field.
func DiscogsMasterIdNotNil() predicate.Album {
	return predicate.Album(sql.FieldNotNull(FieldDiscogsMasterId))
}

// DiscogsMasterIdEqualFold applies the EqualFold predicate on the "discogsMasterId" field.
func DiscogsMasterIdEqualFold(v string) predicate.Album {
	return predicate.Album(sql.FieldEqualFold(FieldDiscogsMasterId, v))
}

// DiscogsMasterIdContainsFold applies the ContainsFold predicate on the "discogsMasterId" field.
func DiscogsMasterIdContainsFold(v string) predicate.Album {
	return predicate.Album(sql.FieldContainsFold(FieldDiscogsMasterId, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Album {
	return predicate.Album(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Album {
	return predicate.Album(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Album {
	return predicate.Album(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Album {
	return predicate.Album(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Album {
	return predicate.Album(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Album {
	return predicate.Album(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Album {
	return predicate.Album(sql.FieldContainsFold(FieldName, v))
}

// SimplifiedNameEQ applies the EQ predicate on the "simplifiedName" field.
func SimplifiedNameEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldEQ(FieldSimplifiedName, v))
}

// SimplifiedNameNEQ applies the NEQ predicate on the "simplifiedName" field.
func SimplifiedNameNEQ(v string) predicate.Album {
	return predicate.Album(sql.FieldNEQ(FieldSimplifiedName, v))
}

// SimplifiedNameIn applies the In predicate on the "simplifiedName" field.
func SimplifiedNameIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldIn(FieldSimplifiedName, vs...))
}

// SimplifiedNameNotIn applies the NotIn predicate on the "simplifiedName" field.
func SimplifiedNameNotIn(vs ...string) predicate.Album {
	return predicate.Album(sql.FieldNotIn(FieldSimplifiedName, vs...))
}

// SimplifiedNameGT applies the GT predicate on the "simplifiedName" field.
func SimplifiedNameGT(v string) predicate.Album {
	return predicate.Album(sql.FieldGT(FieldSimplifiedName, v))
}

// SimplifiedNameGTE applies the GTE predicate on the "simplifiedName" field.
func SimplifiedNameGTE(v string) predicate.Album {
	return predicate.Album(sql.FieldGTE(FieldSimplifiedName, v))
}

// SimplifiedNameLT applies the LT predicate on the "simplifiedName" field.
func SimplifiedNameLT(v string) predicate.Album {
	return predicate.Album(sql.FieldLT(FieldSimplifiedName, v))
}

// SimplifiedNameLTE applies the LTE predicate on the "simplifiedName" field.
func SimplifiedNameLTE(v string) predicate.Album {
	return predicate.Album(sql.FieldLTE(FieldSimplifiedName, v))
}

// SimplifiedNameContains applies the Contains predicate on the "simplifiedName" field.
func SimplifiedNameContains(v string) predicate.Album {
	return predicate.Album(sql.FieldContains(FieldSimplifiedName, v))
}

// SimplifiedNameHasPrefix applies the HasPrefix predicate on the "simplifiedName" field.
func SimplifiedNameHasPrefix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasPrefix(FieldSimplifiedName, v))
}

// SimplifiedNameHasSuffix applies the HasSuffix predicate on the "simplifiedName" field.
func SimplifiedNameHasSuffix(v string) predicate.Album {
	return predicate.Album(sql.FieldHasSuffix(FieldSimplifiedName, v))
}

// SimplifiedNameEqualFold applies the EqualFold predicate on the "simplifiedName" field.
func SimplifiedNameEqualFold(v string) predicate.Album {
	return predicate.Album(sql.FieldEqualFold(FieldSimplifiedName, v))
}

// SimplifiedNameContainsFold applies the ContainsFold predicate on the "simplifiedName" field.
func SimplifiedNameContainsFold(v string) predicate.Album {
	return predicate.Album(sql.FieldContainsFold(FieldSimplifiedName, v))
}

// HasTracks applies the HasEdge predicate on the "tracks" edge.
func HasTracks() predicate.Album {
	return predicate.Album(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TracksTable, TracksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTracksWith applies the HasEdge predicate on the "tracks" edge with a given conditions (other predicates).
func HasTracksWith(preds ...predicate.Track) predicate.Album {
	return predicate.Album(func(s *sql.Selector) {
		step := newTracksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasArtists applies the HasEdge predicate on the "artists" edge.
func HasArtists() predicate.Album {
	return predicate.Album(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ArtistsTable, ArtistsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasArtistsWith applies the HasEdge predicate on the "artists" edge with a given conditions (other predicates).
func HasArtistsWith(preds ...predicate.Artist) predicate.Album {
	return predicate.Album(func(s *sql.Selector) {
		step := newArtistsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Album) predicate.Album {
	return predicate.Album(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Album) predicate.Album {
	return predicate.Album(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Album) predicate.Album {
	return predicate.Album(sql.NotPredicates(p))
}
