package brain

import "github.com/jmoiron/sqlx"

const _INIT_DB_QUERY = `
CREATE TABLE IF NOT EXISTS artists (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	spotify_id text NOT NULL,
	name text NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	username text PRIMARY KEY NOT NULL,
	spotify_token_str blob
);

CREATE TABLE IF NOT EXISTS meta_albums (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	simplified_name text NOT NULL,
	any_name text NOT NULL
);

CREATE TABLE IF NOT EXISTS meta_album_artists (
	meta_album_id integer NOT NULL,
	artist_id integer NOT NULL,
	PRIMARY KEY (meta_album_id, artist_id),
	CONSTRAINT fk_meta_album_artists_meta_album
		FOREIGN KEY (meta_album_id) REFERENCES meta_albums(id),
	CONSTRAINT fk_meta_album_artists_artist
		FOREIGN KEY (artist_id) REFERENCES artists(id)
);
CREATE TABLE IF NOT EXISTS spotify_albums (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	spotify_id text NOT NULL,
	name text NOT NULL,
	meta_album_id integer NOT NULL,
	CONSTRAINT fk_spotify_albums_meta_album
		FOREIGN KEY (meta_album_id) REFERENCES meta_albums(id)
);

CREATE TABLE IF NOT EXISTS spotify_album_artists (
	spotify_album_id integer NOT NULL,
	artist_id integer NOT NULL,
	PRIMARY KEY (spotify_album_id, artist_id),
	CONSTRAINT fk_spotify_album_artists_spotify_album
		FOREIGN KEY (spotify_album_id) REFERENCES spotify_albums(id),
	CONSTRAINT fk_spotify_album_artists_artist
		FOREIGN KEY (artist_id) REFERENCES artists(id)
);

CREATE TABLE IF NOT EXISTS meta_track_artists (
	meta_track_id integer NOT NULL,
	artist_id integer NOT NULL,
	PRIMARY KEY (meta_track_id,artist_id),
	CONSTRAINT fk_meta_track_artists_meta_track
		FOREIGN KEY (meta_track_id) REFERENCES meta_tracks(id),
	CONSTRAINT fk_meta_track_artists_artist
		FOREIGN KEY (artist_id) REFERENCES artists(id)
);

CREATE TABLE IF NOT EXISTS meta_tracks (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	simplified_name text NOT NULL,
	meta_album_id integer NOT NULL,
	CONSTRAINT fk_meta_tracks_meta_album
		FOREIGN KEY (meta_album_id) REFERENCES meta_albums(id)
);

CREATE TABLE IF NOT EXISTS spotify_tracks (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	spotify_id text NOT NULL,
	name text NOT NULL,
	spotify_album_id integer NOT NULL,
	meta_track_id integer NOT NULL,
	CONSTRAINT fk_spotify_tracks_spotify_album
		FOREIGN KEY (spotify_album_id) REFERENCES spotify_albums(id),
	CONSTRAINT fk_spotify_tracks_meta_track
		FOREIGN KEY (meta_track_id) REFERENCES meta_tracks(id)
);

CREATE TABLE IF NOT EXISTS spotify_track_artists (
	spotify_track_id integer NOT NULL,
	artist_id integer NOT NULL,
	PRIMARY KEY (spotify_track_id,artist_id),
	CONSTRAINT fk_spotify_track_artists_spotify_track
		FOREIGN KEY (spotify_track_id) REFERENCES spotify_tracks(id),
	CONSTRAINT fk_spotify_track_artists_artist
		FOREIGN KEY (artist_id) REFERENCES artists(id)
);

CREATE TABLE IF NOT EXISTS user_saved_tracks (
	user_username text NOT NULL,
	meta_track_id integer NOT NULL,
	PRIMARY KEY (user_username,meta_track_id),
	CONSTRAINT fk_user_saved_tracks_user
		FOREIGN KEY (user_username) REFERENCES users(username),
	CONSTRAINT fk_user_saved_tracks_meta_track
		FOREIGN KEY (meta_track_id) REFERENCES meta_tracks(id)
);

CREATE TABLE IF NOT EXISTS discogs_releases (
	id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	discogs_id integer NOT NULL,
	name text NOT NULL,
	artist_name text NOT NULL,
	format text NOT NULL,
	searched_meta_album integer DEFAULT FALSE NOT NULL,
	meta_album_score integer DEFAULT 0 NOT NULL,
	meta_album_id integer,
	CONSTRAINT fk_discogs_releases_meta_album
		FOREIGN KEY (meta_album_id) REFERENCES meta_albums(id)
);

CREATE TABLE IF NOT EXISTS discogs_sellers (
	username text PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS discogs_seller_selling_releases (
	discogs_seller_username text NOT NULL,
	discogs_release_id integer NOT NULL,
	PRIMARY KEY (discogs_seller_username,discogs_release_id),
	CONSTRAINT fk_discogs_seller_selling_releases_discogs_seller
		FOREIGN KEY (discogs_seller_username) REFERENCES discogs_sellers(username),
	CONSTRAINT fk_discogs_seller_selling_releases_discogs_release
		FOREIGN KEY (discogs_release_id) REFERENCES discogs_releases(id)
);

`

func initSqlx(db *sqlx.DB) error {
	if _, err := db.Exec(_INIT_DB_QUERY); err != nil {
		return err
	}
	return nil
}
