package album

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/m-nny/universe/ent/predicate"
)

func SpotifyIdContains(spotifyId string) predicate.Album {
	return predicate.Album(func(s *sql.Selector) {
		s.Where(sqljson.ValueContains(FieldSpotifyIds, spotifyId))
	})
}

func Similar(spotifyId string, simplifiedName string) predicate.Album {
	return Or(
		SpotifyIdContains(spotifyId),
		SimplifiedName(simplifiedName),
	)
}
