// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/ent/playlist"
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
	// albumDescID is the schema descriptor for id field.
	albumDescID := albumFields[0].Descriptor()
	// album.IDValidator is a validator for the "id" field. It is called by the builders before save.
	album.IDValidator = albumDescID.Validators[0].(func(string) error)
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
	playlistFields := schema.Playlist{}.Fields()
	_ = playlistFields
	// playlistDescName is the schema descriptor for name field.
	playlistDescName := playlistFields[1].Descriptor()
	// playlist.NameValidator is a validator for the "name" field. It is called by the builders before save.
	playlist.NameValidator = playlistDescName.Validators[0].(func(string) error)
	// playlistDescSnaphotID is the schema descriptor for snaphot_id field.
	playlistDescSnaphotID := playlistFields[2].Descriptor()
	// playlist.SnaphotIDValidator is a validator for the "snaphot_id" field. It is called by the builders before save.
	playlist.SnaphotIDValidator = playlistDescSnaphotID.Validators[0].(func(string) error)
	// playlistDescID is the schema descriptor for id field.
	playlistDescID := playlistFields[0].Descriptor()
	// playlist.IDValidator is a validator for the "id" field. It is called by the builders before save.
	playlist.IDValidator = playlistDescID.Validators[0].(func(string) error)
	trackFields := schema.Track{}.Fields()
	_ = trackFields
	// trackDescName is the schema descriptor for name field.
	trackDescName := trackFields[1].Descriptor()
	// track.NameValidator is a validator for the "name" field. It is called by the builders before save.
	track.NameValidator = trackDescName.Validators[0].(func(string) error)
	// trackDescID is the schema descriptor for id field.
	trackDescID := trackFields[0].Descriptor()
	// track.IDValidator is a validator for the "id" field. It is called by the builders before save.
	track.IDValidator = trackDescID.Validators[0].(func(string) error)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.IDValidator is a validator for the "id" field. It is called by the builders before save.
	user.IDValidator = userDescID.Validators[0].(func(string) error)
}
