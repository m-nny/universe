// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/m-nny/universe/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/ent/playlist"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/ent/user"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Album is the client for interacting with the Album builders.
	Album *AlbumClient
	// Artist is the client for interacting with the Artist builders.
	Artist *ArtistClient
	// Playlist is the client for interacting with the Playlist builders.
	Playlist *PlaylistClient
	// Track is the client for interacting with the Track builders.
	Track *TrackClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Album = NewAlbumClient(c.config)
	c.Artist = NewArtistClient(c.config)
	c.Playlist = NewPlaylistClient(c.config)
	c.Track = NewTrackClient(c.config)
	c.User = NewUserClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:      ctx,
		config:   cfg,
		Album:    NewAlbumClient(cfg),
		Artist:   NewArtistClient(cfg),
		Playlist: NewPlaylistClient(cfg),
		Track:    NewTrackClient(cfg),
		User:     NewUserClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:      ctx,
		config:   cfg,
		Album:    NewAlbumClient(cfg),
		Artist:   NewArtistClient(cfg),
		Playlist: NewPlaylistClient(cfg),
		Track:    NewTrackClient(cfg),
		User:     NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Album.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Album.Use(hooks...)
	c.Artist.Use(hooks...)
	c.Playlist.Use(hooks...)
	c.Track.Use(hooks...)
	c.User.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Album.Intercept(interceptors...)
	c.Artist.Intercept(interceptors...)
	c.Playlist.Intercept(interceptors...)
	c.Track.Intercept(interceptors...)
	c.User.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *AlbumMutation:
		return c.Album.mutate(ctx, m)
	case *ArtistMutation:
		return c.Artist.mutate(ctx, m)
	case *PlaylistMutation:
		return c.Playlist.mutate(ctx, m)
	case *TrackMutation:
		return c.Track.mutate(ctx, m)
	case *UserMutation:
		return c.User.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// AlbumClient is a client for the Album schema.
type AlbumClient struct {
	config
}

// NewAlbumClient returns a client for the Album from the given config.
func NewAlbumClient(c config) *AlbumClient {
	return &AlbumClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `album.Hooks(f(g(h())))`.
func (c *AlbumClient) Use(hooks ...Hook) {
	c.hooks.Album = append(c.hooks.Album, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `album.Intercept(f(g(h())))`.
func (c *AlbumClient) Intercept(interceptors ...Interceptor) {
	c.inters.Album = append(c.inters.Album, interceptors...)
}

// Create returns a builder for creating a Album entity.
func (c *AlbumClient) Create() *AlbumCreate {
	mutation := newAlbumMutation(c.config, OpCreate)
	return &AlbumCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Album entities.
func (c *AlbumClient) CreateBulk(builders ...*AlbumCreate) *AlbumCreateBulk {
	return &AlbumCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *AlbumClient) MapCreateBulk(slice any, setFunc func(*AlbumCreate, int)) *AlbumCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &AlbumCreateBulk{err: fmt.Errorf("calling to AlbumClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*AlbumCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &AlbumCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Album.
func (c *AlbumClient) Update() *AlbumUpdate {
	mutation := newAlbumMutation(c.config, OpUpdate)
	return &AlbumUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *AlbumClient) UpdateOne(a *Album) *AlbumUpdateOne {
	mutation := newAlbumMutation(c.config, OpUpdateOne, withAlbum(a))
	return &AlbumUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *AlbumClient) UpdateOneID(id int) *AlbumUpdateOne {
	mutation := newAlbumMutation(c.config, OpUpdateOne, withAlbumID(id))
	return &AlbumUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Album.
func (c *AlbumClient) Delete() *AlbumDelete {
	mutation := newAlbumMutation(c.config, OpDelete)
	return &AlbumDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *AlbumClient) DeleteOne(a *Album) *AlbumDeleteOne {
	return c.DeleteOneID(a.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *AlbumClient) DeleteOneID(id int) *AlbumDeleteOne {
	builder := c.Delete().Where(album.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &AlbumDeleteOne{builder}
}

// Query returns a query builder for Album.
func (c *AlbumClient) Query() *AlbumQuery {
	return &AlbumQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeAlbum},
		inters: c.Interceptors(),
	}
}

// Get returns a Album entity by its id.
func (c *AlbumClient) Get(ctx context.Context, id int) (*Album, error) {
	return c.Query().Where(album.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *AlbumClient) GetX(ctx context.Context, id int) *Album {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTracks queries the tracks edge of a Album.
func (c *AlbumClient) QueryTracks(a *Album) *TrackQuery {
	query := (&TrackClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(album.Table, album.FieldID, id),
			sqlgraph.To(track.Table, track.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, album.TracksTable, album.TracksColumn),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryArtists queries the artists edge of a Album.
func (c *AlbumClient) QueryArtists(a *Album) *ArtistQuery {
	query := (&ArtistClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(album.Table, album.FieldID, id),
			sqlgraph.To(artist.Table, artist.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, album.ArtistsTable, album.ArtistsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *AlbumClient) Hooks() []Hook {
	return c.hooks.Album
}

// Interceptors returns the client interceptors.
func (c *AlbumClient) Interceptors() []Interceptor {
	return c.inters.Album
}

func (c *AlbumClient) mutate(ctx context.Context, m *AlbumMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&AlbumCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&AlbumUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&AlbumUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&AlbumDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Album mutation op: %q", m.Op())
	}
}

// ArtistClient is a client for the Artist schema.
type ArtistClient struct {
	config
}

// NewArtistClient returns a client for the Artist from the given config.
func NewArtistClient(c config) *ArtistClient {
	return &ArtistClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `artist.Hooks(f(g(h())))`.
func (c *ArtistClient) Use(hooks ...Hook) {
	c.hooks.Artist = append(c.hooks.Artist, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `artist.Intercept(f(g(h())))`.
func (c *ArtistClient) Intercept(interceptors ...Interceptor) {
	c.inters.Artist = append(c.inters.Artist, interceptors...)
}

// Create returns a builder for creating a Artist entity.
func (c *ArtistClient) Create() *ArtistCreate {
	mutation := newArtistMutation(c.config, OpCreate)
	return &ArtistCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Artist entities.
func (c *ArtistClient) CreateBulk(builders ...*ArtistCreate) *ArtistCreateBulk {
	return &ArtistCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *ArtistClient) MapCreateBulk(slice any, setFunc func(*ArtistCreate, int)) *ArtistCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &ArtistCreateBulk{err: fmt.Errorf("calling to ArtistClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*ArtistCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &ArtistCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Artist.
func (c *ArtistClient) Update() *ArtistUpdate {
	mutation := newArtistMutation(c.config, OpUpdate)
	return &ArtistUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ArtistClient) UpdateOne(a *Artist) *ArtistUpdateOne {
	mutation := newArtistMutation(c.config, OpUpdateOne, withArtist(a))
	return &ArtistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ArtistClient) UpdateOneID(id int) *ArtistUpdateOne {
	mutation := newArtistMutation(c.config, OpUpdateOne, withArtistID(id))
	return &ArtistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Artist.
func (c *ArtistClient) Delete() *ArtistDelete {
	mutation := newArtistMutation(c.config, OpDelete)
	return &ArtistDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ArtistClient) DeleteOne(a *Artist) *ArtistDeleteOne {
	return c.DeleteOneID(a.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ArtistClient) DeleteOneID(id int) *ArtistDeleteOne {
	builder := c.Delete().Where(artist.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ArtistDeleteOne{builder}
}

// Query returns a query builder for Artist.
func (c *ArtistClient) Query() *ArtistQuery {
	return &ArtistQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeArtist},
		inters: c.Interceptors(),
	}
}

// Get returns a Artist entity by its id.
func (c *ArtistClient) Get(ctx context.Context, id int) (*Artist, error) {
	return c.Query().Where(artist.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ArtistClient) GetX(ctx context.Context, id int) *Artist {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTracks queries the tracks edge of a Artist.
func (c *ArtistClient) QueryTracks(a *Artist) *TrackQuery {
	query := (&TrackClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(artist.Table, artist.FieldID, id),
			sqlgraph.To(track.Table, track.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, artist.TracksTable, artist.TracksPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryAlbums queries the albums edge of a Artist.
func (c *ArtistClient) QueryAlbums(a *Artist) *AlbumQuery {
	query := (&AlbumClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(artist.Table, artist.FieldID, id),
			sqlgraph.To(album.Table, album.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, artist.AlbumsTable, artist.AlbumsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ArtistClient) Hooks() []Hook {
	return c.hooks.Artist
}

// Interceptors returns the client interceptors.
func (c *ArtistClient) Interceptors() []Interceptor {
	return c.inters.Artist
}

func (c *ArtistClient) mutate(ctx context.Context, m *ArtistMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ArtistCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ArtistUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ArtistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ArtistDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Artist mutation op: %q", m.Op())
	}
}

// PlaylistClient is a client for the Playlist schema.
type PlaylistClient struct {
	config
}

// NewPlaylistClient returns a client for the Playlist from the given config.
func NewPlaylistClient(c config) *PlaylistClient {
	return &PlaylistClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `playlist.Hooks(f(g(h())))`.
func (c *PlaylistClient) Use(hooks ...Hook) {
	c.hooks.Playlist = append(c.hooks.Playlist, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `playlist.Intercept(f(g(h())))`.
func (c *PlaylistClient) Intercept(interceptors ...Interceptor) {
	c.inters.Playlist = append(c.inters.Playlist, interceptors...)
}

// Create returns a builder for creating a Playlist entity.
func (c *PlaylistClient) Create() *PlaylistCreate {
	mutation := newPlaylistMutation(c.config, OpCreate)
	return &PlaylistCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Playlist entities.
func (c *PlaylistClient) CreateBulk(builders ...*PlaylistCreate) *PlaylistCreateBulk {
	return &PlaylistCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *PlaylistClient) MapCreateBulk(slice any, setFunc func(*PlaylistCreate, int)) *PlaylistCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &PlaylistCreateBulk{err: fmt.Errorf("calling to PlaylistClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*PlaylistCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &PlaylistCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Playlist.
func (c *PlaylistClient) Update() *PlaylistUpdate {
	mutation := newPlaylistMutation(c.config, OpUpdate)
	return &PlaylistUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PlaylistClient) UpdateOne(pl *Playlist) *PlaylistUpdateOne {
	mutation := newPlaylistMutation(c.config, OpUpdateOne, withPlaylist(pl))
	return &PlaylistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PlaylistClient) UpdateOneID(id string) *PlaylistUpdateOne {
	mutation := newPlaylistMutation(c.config, OpUpdateOne, withPlaylistID(id))
	return &PlaylistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Playlist.
func (c *PlaylistClient) Delete() *PlaylistDelete {
	mutation := newPlaylistMutation(c.config, OpDelete)
	return &PlaylistDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PlaylistClient) DeleteOne(pl *Playlist) *PlaylistDeleteOne {
	return c.DeleteOneID(pl.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PlaylistClient) DeleteOneID(id string) *PlaylistDeleteOne {
	builder := c.Delete().Where(playlist.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PlaylistDeleteOne{builder}
}

// Query returns a query builder for Playlist.
func (c *PlaylistClient) Query() *PlaylistQuery {
	return &PlaylistQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePlaylist},
		inters: c.Interceptors(),
	}
}

// Get returns a Playlist entity by its id.
func (c *PlaylistClient) Get(ctx context.Context, id string) (*Playlist, error) {
	return c.Query().Where(playlist.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PlaylistClient) GetX(ctx context.Context, id string) *Playlist {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryOwner queries the owner edge of a Playlist.
func (c *PlaylistClient) QueryOwner(pl *Playlist) *UserQuery {
	query := (&UserClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pl.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(playlist.Table, playlist.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, playlist.OwnerTable, playlist.OwnerColumn),
		)
		fromV = sqlgraph.Neighbors(pl.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PlaylistClient) Hooks() []Hook {
	return c.hooks.Playlist
}

// Interceptors returns the client interceptors.
func (c *PlaylistClient) Interceptors() []Interceptor {
	return c.inters.Playlist
}

func (c *PlaylistClient) mutate(ctx context.Context, m *PlaylistMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PlaylistCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PlaylistUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PlaylistUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PlaylistDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Playlist mutation op: %q", m.Op())
	}
}

// TrackClient is a client for the Track schema.
type TrackClient struct {
	config
}

// NewTrackClient returns a client for the Track from the given config.
func NewTrackClient(c config) *TrackClient {
	return &TrackClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `track.Hooks(f(g(h())))`.
func (c *TrackClient) Use(hooks ...Hook) {
	c.hooks.Track = append(c.hooks.Track, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `track.Intercept(f(g(h())))`.
func (c *TrackClient) Intercept(interceptors ...Interceptor) {
	c.inters.Track = append(c.inters.Track, interceptors...)
}

// Create returns a builder for creating a Track entity.
func (c *TrackClient) Create() *TrackCreate {
	mutation := newTrackMutation(c.config, OpCreate)
	return &TrackCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Track entities.
func (c *TrackClient) CreateBulk(builders ...*TrackCreate) *TrackCreateBulk {
	return &TrackCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TrackClient) MapCreateBulk(slice any, setFunc func(*TrackCreate, int)) *TrackCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TrackCreateBulk{err: fmt.Errorf("calling to TrackClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TrackCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TrackCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Track.
func (c *TrackClient) Update() *TrackUpdate {
	mutation := newTrackMutation(c.config, OpUpdate)
	return &TrackUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TrackClient) UpdateOne(t *Track) *TrackUpdateOne {
	mutation := newTrackMutation(c.config, OpUpdateOne, withTrack(t))
	return &TrackUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TrackClient) UpdateOneID(id int) *TrackUpdateOne {
	mutation := newTrackMutation(c.config, OpUpdateOne, withTrackID(id))
	return &TrackUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Track.
func (c *TrackClient) Delete() *TrackDelete {
	mutation := newTrackMutation(c.config, OpDelete)
	return &TrackDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TrackClient) DeleteOne(t *Track) *TrackDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TrackClient) DeleteOneID(id int) *TrackDeleteOne {
	builder := c.Delete().Where(track.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TrackDeleteOne{builder}
}

// Query returns a query builder for Track.
func (c *TrackClient) Query() *TrackQuery {
	return &TrackQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTrack},
		inters: c.Interceptors(),
	}
}

// Get returns a Track entity by its id.
func (c *TrackClient) Get(ctx context.Context, id int) (*Track, error) {
	return c.Query().Where(track.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TrackClient) GetX(ctx context.Context, id int) *Track {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySavedBy queries the savedBy edge of a Track.
func (c *TrackClient) QuerySavedBy(t *Track) *UserQuery {
	query := (&UserClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(track.Table, track.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, track.SavedByTable, track.SavedByPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryAlbum queries the album edge of a Track.
func (c *TrackClient) QueryAlbum(t *Track) *AlbumQuery {
	query := (&AlbumClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(track.Table, track.FieldID, id),
			sqlgraph.To(album.Table, album.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, track.AlbumTable, track.AlbumColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryArtists queries the artists edge of a Track.
func (c *TrackClient) QueryArtists(t *Track) *ArtistQuery {
	query := (&ArtistClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(track.Table, track.FieldID, id),
			sqlgraph.To(artist.Table, artist.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, track.ArtistsTable, track.ArtistsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TrackClient) Hooks() []Hook {
	return c.hooks.Track
}

// Interceptors returns the client interceptors.
func (c *TrackClient) Interceptors() []Interceptor {
	return c.inters.Track
}

func (c *TrackClient) mutate(ctx context.Context, m *TrackMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TrackCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TrackUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TrackUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TrackDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Track mutation op: %q", m.Op())
	}
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `user.Intercept(f(g(h())))`.
func (c *UserClient) Intercept(interceptors ...Interceptor) {
	c.inters.User = append(c.inters.User, interceptors...)
}

// Create returns a builder for creating a User entity.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *UserClient) MapCreateBulk(slice any, setFunc func(*UserCreate, int)) *UserCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &UserCreateBulk{err: fmt.Errorf("calling to UserClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*UserCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id string) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *UserClient) DeleteOneID(id string) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeUser},
		inters: c.Interceptors(),
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id string) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id string) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPlaylists queries the playlists edge of a User.
func (c *UserClient) QueryPlaylists(u *User) *PlaylistQuery {
	query := (&PlaylistClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(playlist.Table, playlist.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.PlaylistsTable, user.PlaylistsColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySavedTracks queries the savedTracks edge of a User.
func (c *UserClient) QuerySavedTracks(u *User) *TrackQuery {
	query := (&TrackClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(track.Table, track.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, user.SavedTracksTable, user.SavedTracksPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}

// Interceptors returns the client interceptors.
func (c *UserClient) Interceptors() []Interceptor {
	return c.inters.User
}

func (c *UserClient) mutate(ctx context.Context, m *UserMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&UserCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&UserUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&UserDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown User mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Album, Artist, Playlist, Track, User []ent.Hook
	}
	inters struct {
		Album, Artist, Playlist, Track, User []ent.Interceptor
	}
)
