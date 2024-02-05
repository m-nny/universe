// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/ent/schema"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	albumFields := schema.Album{}.Fields()
	_ = albumFields
	// albumDescName is the schema descriptor for name field.
	albumDescName := albumFields[1].Descriptor()
	// album.NameValidator is a validator for the "name" field. It is called by the builders before save.
	album.NameValidator = albumDescName.Validators[0].(func(string) error)
	// albumDescSimplifiedName is the schema descriptor for simplifiedName field.
	albumDescSimplifiedName := albumFields[2].Descriptor()
	// album.SimplifiedNameValidator is a validator for the "simplifiedName" field. It is called by the builders before save.
	album.SimplifiedNameValidator = albumDescSimplifiedName.Validators[0].(func(string) error)
	artistFields := schema.Artist{}.Fields()
	_ = artistFields
	// artistDescSpotifyId is the schema descriptor for spotifyId field.
	artistDescSpotifyId := artistFields[0].Descriptor()
	// artist.SpotifyIdValidator is a validator for the "spotifyId" field. It is called by the builders before save.
	artist.SpotifyIdValidator = artistDescSpotifyId.Validators[0].(func(string) error)
	// artistDescName is the schema descriptor for name field.
	artistDescName := artistFields[1].Descriptor()
	// artist.NameValidator is a validator for the "name" field. It is called by the builders before save.
	artist.NameValidator = artistDescName.Validators[0].(func(string) error)
	trackFields := schema.Track{}.Fields()
	_ = trackFields
	// trackDescName is the schema descriptor for name field.
	trackDescName := trackFields[1].Descriptor()
	// track.NameValidator is a validator for the "name" field. It is called by the builders before save.
	track.NameValidator = trackDescName.Validators[0].(func(string) error)
	// trackDescTrackNumber is the schema descriptor for trackNumber field.
	trackDescTrackNumber := trackFields[2].Descriptor()
	// track.TrackNumberValidator is a validator for the "trackNumber" field. It is called by the builders before save.
	track.TrackNumberValidator = trackDescTrackNumber.Validators[0].(func(int) error)
	// trackDescSimplifiedName is the schema descriptor for simplifiedName field.
	trackDescSimplifiedName := trackFields[3].Descriptor()
	// track.SimplifiedNameValidator is a validator for the "simplifiedName" field. It is called by the builders before save.
	track.SimplifiedNameValidator = trackDescSimplifiedName.Validators[0].(func(string) error)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.IDValidator is a validator for the "id" field. It is called by the builders before save.
	user.IDValidator = userDescID.Validators[0].(func(string) error)
}
