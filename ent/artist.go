// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/m-nny/universe/ent/artist"
)

// Artist is the model entity for the Artist schema.
type Artist struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SpotifyId holds the value of the "spotifyId" field.
	SpotifyId string `json:"spotifyId,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ArtistQuery when eager-loading is set.
	Edges        ArtistEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ArtistEdges holds the relations/edges for other nodes in the graph.
type ArtistEdges struct {
	// Tracks holds the value of the tracks edge.
	Tracks []*Track `json:"tracks,omitempty"`
	// Albums holds the value of the albums edge.
	Albums []*Album `json:"albums,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TracksOrErr returns the Tracks value or an error if the edge
// was not loaded in eager-loading.
func (e ArtistEdges) TracksOrErr() ([]*Track, error) {
	if e.loadedTypes[0] {
		return e.Tracks, nil
	}
	return nil, &NotLoadedError{edge: "tracks"}
}

// AlbumsOrErr returns the Albums value or an error if the edge
// was not loaded in eager-loading.
func (e ArtistEdges) AlbumsOrErr() ([]*Album, error) {
	if e.loadedTypes[1] {
		return e.Albums, nil
	}
	return nil, &NotLoadedError{edge: "albums"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Artist) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case artist.FieldID:
			values[i] = new(sql.NullInt64)
		case artist.FieldSpotifyId, artist.FieldName:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Artist fields.
func (a *Artist) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case artist.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = int(value.Int64)
		case artist.FieldSpotifyId:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field spotifyId", values[i])
			} else if value.Valid {
				a.SpotifyId = value.String
			}
		case artist.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Artist.
// This includes values selected through modifiers, order, etc.
func (a *Artist) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryTracks queries the "tracks" edge of the Artist entity.
func (a *Artist) QueryTracks() *TrackQuery {
	return NewArtistClient(a.config).QueryTracks(a)
}

// QueryAlbums queries the "albums" edge of the Artist entity.
func (a *Artist) QueryAlbums() *AlbumQuery {
	return NewArtistClient(a.config).QueryAlbums(a)
}

// Update returns a builder for updating this Artist.
// Note that you need to call Artist.Unwrap() before calling this method if this Artist
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Artist) Update() *ArtistUpdateOne {
	return NewArtistClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Artist entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Artist) Unwrap() *Artist {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Artist is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Artist) String() string {
	var builder strings.Builder
	builder.WriteString("Artist(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("spotifyId=")
	builder.WriteString(a.SpotifyId)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteByte(')')
	return builder.String()
}

// Artists is a parsable slice of Artist.
type Artists []*Artist
