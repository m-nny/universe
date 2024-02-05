// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/ent/user"
)

// TrackCreate is the builder for creating a Track entity.
type TrackCreate struct {
	config
	mutation *TrackMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (tc *TrackCreate) SetName(s string) *TrackCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetID sets the "id" field.
func (tc *TrackCreate) SetID(s string) *TrackCreate {
	tc.mutation.SetID(s)
	return tc
}

// AddSavedByIDs adds the "savedBy" edge to the User entity by IDs.
func (tc *TrackCreate) AddSavedByIDs(ids ...string) *TrackCreate {
	tc.mutation.AddSavedByIDs(ids...)
	return tc
}

// AddSavedBy adds the "savedBy" edges to the User entity.
func (tc *TrackCreate) AddSavedBy(u ...*User) *TrackCreate {
	ids := make([]string, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tc.AddSavedByIDs(ids...)
}

// SetAlbumID sets the "album" edge to the Album entity by ID.
func (tc *TrackCreate) SetAlbumID(id int) *TrackCreate {
	tc.mutation.SetAlbumID(id)
	return tc
}

// SetNillableAlbumID sets the "album" edge to the Album entity by ID if the given value is not nil.
func (tc *TrackCreate) SetNillableAlbumID(id *int) *TrackCreate {
	if id != nil {
		tc = tc.SetAlbumID(*id)
	}
	return tc
}

// SetAlbum sets the "album" edge to the Album entity.
func (tc *TrackCreate) SetAlbum(a *Album) *TrackCreate {
	return tc.SetAlbumID(a.ID)
}

// AddArtistIDs adds the "artists" edge to the Artist entity by IDs.
func (tc *TrackCreate) AddArtistIDs(ids ...int) *TrackCreate {
	tc.mutation.AddArtistIDs(ids...)
	return tc
}

// AddArtists adds the "artists" edges to the Artist entity.
func (tc *TrackCreate) AddArtists(a ...*Artist) *TrackCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tc.AddArtistIDs(ids...)
}

// Mutation returns the TrackMutation object of the builder.
func (tc *TrackCreate) Mutation() *TrackMutation {
	return tc.mutation
}

// Save creates the Track in the database.
func (tc *TrackCreate) Save(ctx context.Context) (*Track, error) {
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TrackCreate) SaveX(ctx context.Context) *Track {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TrackCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TrackCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TrackCreate) check() error {
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Track.name"`)}
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := track.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Track.name": %w`, err)}
		}
	}
	if v, ok := tc.mutation.ID(); ok {
		if err := track.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Track.id": %w`, err)}
		}
	}
	return nil
}

func (tc *TrackCreate) sqlSave(ctx context.Context) (*Track, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Track.ID type: %T", _spec.ID.Value)
		}
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TrackCreate) createSpec() (*Track, *sqlgraph.CreateSpec) {
	var (
		_node = &Track{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(track.Table, sqlgraph.NewFieldSpec(track.FieldID, field.TypeString))
	)
	_spec.OnConflict = tc.conflict
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(track.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if nodes := tc.mutation.SavedByIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   track.SavedByTable,
			Columns: track.SavedByPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.AlbumIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   track.AlbumTable,
			Columns: []string{track.AlbumColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(album.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.album_tracks = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ArtistsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   track.ArtistsTable,
			Columns: track.ArtistsPrimaryKey,
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
//	client.Track.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TrackUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (tc *TrackCreate) OnConflict(opts ...sql.ConflictOption) *TrackUpsertOne {
	tc.conflict = opts
	return &TrackUpsertOne{
		create: tc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Track.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tc *TrackCreate) OnConflictColumns(columns ...string) *TrackUpsertOne {
	tc.conflict = append(tc.conflict, sql.ConflictColumns(columns...))
	return &TrackUpsertOne{
		create: tc,
	}
}

type (
	// TrackUpsertOne is the builder for "upsert"-ing
	//  one Track node.
	TrackUpsertOne struct {
		create *TrackCreate
	}

	// TrackUpsert is the "OnConflict" setter.
	TrackUpsert struct {
		*sql.UpdateSet
	}
)

// SetName sets the "name" field.
func (u *TrackUpsert) SetName(v string) *TrackUpsert {
	u.Set(track.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TrackUpsert) UpdateName() *TrackUpsert {
	u.SetExcluded(track.FieldName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Track.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(track.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TrackUpsertOne) UpdateNewValues() *TrackUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(track.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Track.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *TrackUpsertOne) Ignore() *TrackUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TrackUpsertOne) DoNothing() *TrackUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TrackCreate.OnConflict
// documentation for more info.
func (u *TrackUpsertOne) Update(set func(*TrackUpsert)) *TrackUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TrackUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *TrackUpsertOne) SetName(v string) *TrackUpsertOne {
	return u.Update(func(s *TrackUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TrackUpsertOne) UpdateName() *TrackUpsertOne {
	return u.Update(func(s *TrackUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *TrackUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TrackCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TrackUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *TrackUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: TrackUpsertOne.ID is not supported by MySQL driver. Use TrackUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *TrackUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// TrackCreateBulk is the builder for creating many Track entities in bulk.
type TrackCreateBulk struct {
	config
	err      error
	builders []*TrackCreate
	conflict []sql.ConflictOption
}

// Save creates the Track entities in the database.
func (tcb *TrackCreateBulk) Save(ctx context.Context) ([]*Track, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Track, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TrackMutation)
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
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = tcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TrackCreateBulk) SaveX(ctx context.Context) []*Track {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TrackCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TrackCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Track.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TrackUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (tcb *TrackCreateBulk) OnConflict(opts ...sql.ConflictOption) *TrackUpsertBulk {
	tcb.conflict = opts
	return &TrackUpsertBulk{
		create: tcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Track.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tcb *TrackCreateBulk) OnConflictColumns(columns ...string) *TrackUpsertBulk {
	tcb.conflict = append(tcb.conflict, sql.ConflictColumns(columns...))
	return &TrackUpsertBulk{
		create: tcb,
	}
}

// TrackUpsertBulk is the builder for "upsert"-ing
// a bulk of Track nodes.
type TrackUpsertBulk struct {
	create *TrackCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Track.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(track.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TrackUpsertBulk) UpdateNewValues() *TrackUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(track.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Track.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *TrackUpsertBulk) Ignore() *TrackUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TrackUpsertBulk) DoNothing() *TrackUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TrackCreateBulk.OnConflict
// documentation for more info.
func (u *TrackUpsertBulk) Update(set func(*TrackUpsert)) *TrackUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TrackUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *TrackUpsertBulk) SetName(v string) *TrackUpsertBulk {
	return u.Update(func(s *TrackUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TrackUpsertBulk) UpdateName() *TrackUpsertBulk {
	return u.Update(func(s *TrackUpsert) {
		s.UpdateName()
	})
}

// Exec executes the query.
func (u *TrackUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the TrackCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TrackCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TrackUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
