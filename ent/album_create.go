// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/ent/track"
)

// AlbumCreate is the builder for creating a Album entity.
type AlbumCreate struct {
	config
	mutation *AlbumMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetSpotifyIds sets the "spotifyIds" field.
func (ac *AlbumCreate) SetSpotifyIds(s []string) *AlbumCreate {
	ac.mutation.SetSpotifyIds(s)
	return ac
}

// SetName sets the "name" field.
func (ac *AlbumCreate) SetName(s string) *AlbumCreate {
	ac.mutation.SetName(s)
	return ac
}

// SetSimplifiedName sets the "simplifiedName" field.
func (ac *AlbumCreate) SetSimplifiedName(s string) *AlbumCreate {
	ac.mutation.SetSimplifiedName(s)
	return ac
}

// AddTrackIDs adds the "tracks" edge to the Track entity by IDs.
func (ac *AlbumCreate) AddTrackIDs(ids ...string) *AlbumCreate {
	ac.mutation.AddTrackIDs(ids...)
	return ac
}

// AddTracks adds the "tracks" edges to the Track entity.
func (ac *AlbumCreate) AddTracks(t ...*Track) *AlbumCreate {
	ids := make([]string, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ac.AddTrackIDs(ids...)
}

// AddArtistIDs adds the "artists" edge to the Artist entity by IDs.
func (ac *AlbumCreate) AddArtistIDs(ids ...int) *AlbumCreate {
	ac.mutation.AddArtistIDs(ids...)
	return ac
}

// AddArtists adds the "artists" edges to the Artist entity.
func (ac *AlbumCreate) AddArtists(a ...*Artist) *AlbumCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ac.AddArtistIDs(ids...)
}

// Mutation returns the AlbumMutation object of the builder.
func (ac *AlbumCreate) Mutation() *AlbumMutation {
	return ac.mutation
}

// Save creates the Album in the database.
func (ac *AlbumCreate) Save(ctx context.Context) (*Album, error) {
	return withHooks(ctx, ac.sqlSave, ac.mutation, ac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AlbumCreate) SaveX(ctx context.Context) *Album {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AlbumCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AlbumCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AlbumCreate) check() error {
	if _, ok := ac.mutation.SpotifyIds(); !ok {
		return &ValidationError{Name: "spotifyIds", err: errors.New(`ent: missing required field "Album.spotifyIds"`)}
	}
	if _, ok := ac.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Album.name"`)}
	}
	if v, ok := ac.mutation.Name(); ok {
		if err := album.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Album.name": %w`, err)}
		}
	}
	if _, ok := ac.mutation.SimplifiedName(); !ok {
		return &ValidationError{Name: "simplifiedName", err: errors.New(`ent: missing required field "Album.simplifiedName"`)}
	}
	if v, ok := ac.mutation.SimplifiedName(); ok {
		if err := album.SimplifiedNameValidator(v); err != nil {
			return &ValidationError{Name: "simplifiedName", err: fmt.Errorf(`ent: validator failed for field "Album.simplifiedName": %w`, err)}
		}
	}
	return nil
}

func (ac *AlbumCreate) sqlSave(ctx context.Context) (*Album, error) {
	if err := ac.check(); err != nil {
		return nil, err
	}
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ac.mutation.id = &_node.ID
	ac.mutation.done = true
	return _node, nil
}

func (ac *AlbumCreate) createSpec() (*Album, *sqlgraph.CreateSpec) {
	var (
		_node = &Album{config: ac.config}
		_spec = sqlgraph.NewCreateSpec(album.Table, sqlgraph.NewFieldSpec(album.FieldID, field.TypeInt))
	)
	_spec.OnConflict = ac.conflict
	if value, ok := ac.mutation.SpotifyIds(); ok {
		_spec.SetField(album.FieldSpotifyIds, field.TypeJSON, value)
		_node.SpotifyIds = value
	}
	if value, ok := ac.mutation.Name(); ok {
		_spec.SetField(album.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ac.mutation.SimplifiedName(); ok {
		_spec.SetField(album.FieldSimplifiedName, field.TypeString, value)
		_node.SimplifiedName = value
	}
	if nodes := ac.mutation.TracksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   album.TracksTable,
			Columns: []string{album.TracksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(track.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.ArtistsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   album.ArtistsTable,
			Columns: album.ArtistsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(artist.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Album.Create().
//		SetSpotifyIds(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AlbumUpsert) {
//			SetSpotifyIds(v+v).
//		}).
//		Exec(ctx)
func (ac *AlbumCreate) OnConflict(opts ...sql.ConflictOption) *AlbumUpsertOne {
	ac.conflict = opts
	return &AlbumUpsertOne{
		create: ac,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Album.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ac *AlbumCreate) OnConflictColumns(columns ...string) *AlbumUpsertOne {
	ac.conflict = append(ac.conflict, sql.ConflictColumns(columns...))
	return &AlbumUpsertOne{
		create: ac,
	}
}

type (
	// AlbumUpsertOne is the builder for "upsert"-ing
	//  one Album node.
	AlbumUpsertOne struct {
		create *AlbumCreate
	}

	// AlbumUpsert is the "OnConflict" setter.
	AlbumUpsert struct {
		*sql.UpdateSet
	}
)

// SetSpotifyIds sets the "spotifyIds" field.
func (u *AlbumUpsert) SetSpotifyIds(v []string) *AlbumUpsert {
	u.Set(album.FieldSpotifyIds, v)
	return u
}

// UpdateSpotifyIds sets the "spotifyIds" field to the value that was provided on create.
func (u *AlbumUpsert) UpdateSpotifyIds() *AlbumUpsert {
	u.SetExcluded(album.FieldSpotifyIds)
	return u
}

// SetName sets the "name" field.
func (u *AlbumUpsert) SetName(v string) *AlbumUpsert {
	u.Set(album.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *AlbumUpsert) UpdateName() *AlbumUpsert {
	u.SetExcluded(album.FieldName)
	return u
}

// SetSimplifiedName sets the "simplifiedName" field.
func (u *AlbumUpsert) SetSimplifiedName(v string) *AlbumUpsert {
	u.Set(album.FieldSimplifiedName, v)
	return u
}

// UpdateSimplifiedName sets the "simplifiedName" field to the value that was provided on create.
func (u *AlbumUpsert) UpdateSimplifiedName() *AlbumUpsert {
	u.SetExcluded(album.FieldSimplifiedName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Album.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AlbumUpsertOne) UpdateNewValues() *AlbumUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Album.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AlbumUpsertOne) Ignore() *AlbumUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AlbumUpsertOne) DoNothing() *AlbumUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AlbumCreate.OnConflict
// documentation for more info.
func (u *AlbumUpsertOne) Update(set func(*AlbumUpsert)) *AlbumUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AlbumUpsert{UpdateSet: update})
	}))
	return u
}

// SetSpotifyIds sets the "spotifyIds" field.
func (u *AlbumUpsertOne) SetSpotifyIds(v []string) *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.SetSpotifyIds(v)
	})
}

// UpdateSpotifyIds sets the "spotifyIds" field to the value that was provided on create.
func (u *AlbumUpsertOne) UpdateSpotifyIds() *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateSpotifyIds()
	})
}

// SetName sets the "name" field.
func (u *AlbumUpsertOne) SetName(v string) *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *AlbumUpsertOne) UpdateName() *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateName()
	})
}

// SetSimplifiedName sets the "simplifiedName" field.
func (u *AlbumUpsertOne) SetSimplifiedName(v string) *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.SetSimplifiedName(v)
	})
}

// UpdateSimplifiedName sets the "simplifiedName" field to the value that was provided on create.
func (u *AlbumUpsertOne) UpdateSimplifiedName() *AlbumUpsertOne {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateSimplifiedName()
	})
}

// Exec executes the query.
func (u *AlbumUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AlbumCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AlbumUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AlbumUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AlbumUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AlbumCreateBulk is the builder for creating many Album entities in bulk.
type AlbumCreateBulk struct {
	config
	err      error
	builders []*AlbumCreate
	conflict []sql.ConflictOption
}

// Save creates the Album entities in the database.
func (acb *AlbumCreateBulk) Save(ctx context.Context) ([]*Album, error) {
	if acb.err != nil {
		return nil, acb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Album, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AlbumMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = acb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AlbumCreateBulk) SaveX(ctx context.Context) []*Album {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AlbumCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AlbumCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Album.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AlbumUpsert) {
//			SetSpotifyIds(v+v).
//		}).
//		Exec(ctx)
func (acb *AlbumCreateBulk) OnConflict(opts ...sql.ConflictOption) *AlbumUpsertBulk {
	acb.conflict = opts
	return &AlbumUpsertBulk{
		create: acb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Album.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (acb *AlbumCreateBulk) OnConflictColumns(columns ...string) *AlbumUpsertBulk {
	acb.conflict = append(acb.conflict, sql.ConflictColumns(columns...))
	return &AlbumUpsertBulk{
		create: acb,
	}
}

// AlbumUpsertBulk is the builder for "upsert"-ing
// a bulk of Album nodes.
type AlbumUpsertBulk struct {
	create *AlbumCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Album.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AlbumUpsertBulk) UpdateNewValues() *AlbumUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Album.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AlbumUpsertBulk) Ignore() *AlbumUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AlbumUpsertBulk) DoNothing() *AlbumUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AlbumCreateBulk.OnConflict
// documentation for more info.
func (u *AlbumUpsertBulk) Update(set func(*AlbumUpsert)) *AlbumUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AlbumUpsert{UpdateSet: update})
	}))
	return u
}

// SetSpotifyIds sets the "spotifyIds" field.
func (u *AlbumUpsertBulk) SetSpotifyIds(v []string) *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.SetSpotifyIds(v)
	})
}

// UpdateSpotifyIds sets the "spotifyIds" field to the value that was provided on create.
func (u *AlbumUpsertBulk) UpdateSpotifyIds() *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateSpotifyIds()
	})
}

// SetName sets the "name" field.
func (u *AlbumUpsertBulk) SetName(v string) *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *AlbumUpsertBulk) UpdateName() *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateName()
	})
}

// SetSimplifiedName sets the "simplifiedName" field.
func (u *AlbumUpsertBulk) SetSimplifiedName(v string) *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.SetSimplifiedName(v)
	})
}

// UpdateSimplifiedName sets the "simplifiedName" field to the value that was provided on create.
func (u *AlbumUpsertBulk) UpdateSimplifiedName() *AlbumUpsertBulk {
	return u.Update(func(s *AlbumUpsert) {
		s.UpdateSimplifiedName()
	})
}

// Exec executes the query.
func (u *AlbumUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AlbumCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AlbumCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AlbumUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
