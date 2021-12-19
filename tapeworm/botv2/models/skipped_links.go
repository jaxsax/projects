// Code generated by SQLBoiler 4.8.3 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// SkippedLink is an object representing the database table.
type SkippedLink struct {
	ID          int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	ErrorReason string `boil:"error_reason" json:"error_reason" toml:"error_reason" yaml:"error_reason"`
	Link        string `boil:"link" json:"link" toml:"link" yaml:"link"`

	R *skippedLinkR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L skippedLinkL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SkippedLinkColumns = struct {
	ID          string
	ErrorReason string
	Link        string
}{
	ID:          "id",
	ErrorReason: "error_reason",
	Link:        "link",
}

var SkippedLinkTableColumns = struct {
	ID          string
	ErrorReason string
	Link        string
}{
	ID:          "skipped_links.id",
	ErrorReason: "skipped_links.error_reason",
	Link:        "skipped_links.link",
}

// Generated where

var SkippedLinkWhere = struct {
	ID          whereHelperint64
	ErrorReason whereHelperstring
	Link        whereHelperstring
}{
	ID:          whereHelperint64{field: "\"skipped_links\".\"id\""},
	ErrorReason: whereHelperstring{field: "\"skipped_links\".\"error_reason\""},
	Link:        whereHelperstring{field: "\"skipped_links\".\"link\""},
}

// SkippedLinkRels is where relationship names are stored.
var SkippedLinkRels = struct {
}{}

// skippedLinkR is where relationships are stored.
type skippedLinkR struct {
}

// NewStruct creates a new relationship struct
func (*skippedLinkR) NewStruct() *skippedLinkR {
	return &skippedLinkR{}
}

// skippedLinkL is where Load methods for each relationship are stored.
type skippedLinkL struct{}

var (
	skippedLinkAllColumns            = []string{"id", "error_reason", "link"}
	skippedLinkColumnsWithoutDefault = []string{}
	skippedLinkColumnsWithDefault    = []string{"id", "error_reason", "link"}
	skippedLinkPrimaryKeyColumns     = []string{"id"}
)

type (
	// SkippedLinkSlice is an alias for a slice of pointers to SkippedLink.
	// This should almost always be used instead of []SkippedLink.
	SkippedLinkSlice []*SkippedLink
	// SkippedLinkHook is the signature for custom SkippedLink hook methods
	SkippedLinkHook func(context.Context, boil.ContextExecutor, *SkippedLink) error

	skippedLinkQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	skippedLinkType                 = reflect.TypeOf(&SkippedLink{})
	skippedLinkMapping              = queries.MakeStructMapping(skippedLinkType)
	skippedLinkPrimaryKeyMapping, _ = queries.BindMapping(skippedLinkType, skippedLinkMapping, skippedLinkPrimaryKeyColumns)
	skippedLinkInsertCacheMut       sync.RWMutex
	skippedLinkInsertCache          = make(map[string]insertCache)
	skippedLinkUpdateCacheMut       sync.RWMutex
	skippedLinkUpdateCache          = make(map[string]updateCache)
	skippedLinkUpsertCacheMut       sync.RWMutex
	skippedLinkUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var skippedLinkBeforeInsertHooks []SkippedLinkHook
var skippedLinkBeforeUpdateHooks []SkippedLinkHook
var skippedLinkBeforeDeleteHooks []SkippedLinkHook
var skippedLinkBeforeUpsertHooks []SkippedLinkHook

var skippedLinkAfterInsertHooks []SkippedLinkHook
var skippedLinkAfterSelectHooks []SkippedLinkHook
var skippedLinkAfterUpdateHooks []SkippedLinkHook
var skippedLinkAfterDeleteHooks []SkippedLinkHook
var skippedLinkAfterUpsertHooks []SkippedLinkHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SkippedLink) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SkippedLink) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SkippedLink) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SkippedLink) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SkippedLink) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SkippedLink) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SkippedLink) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SkippedLink) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SkippedLink) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range skippedLinkAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSkippedLinkHook registers your hook function for all future operations.
func AddSkippedLinkHook(hookPoint boil.HookPoint, skippedLinkHook SkippedLinkHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		skippedLinkBeforeInsertHooks = append(skippedLinkBeforeInsertHooks, skippedLinkHook)
	case boil.BeforeUpdateHook:
		skippedLinkBeforeUpdateHooks = append(skippedLinkBeforeUpdateHooks, skippedLinkHook)
	case boil.BeforeDeleteHook:
		skippedLinkBeforeDeleteHooks = append(skippedLinkBeforeDeleteHooks, skippedLinkHook)
	case boil.BeforeUpsertHook:
		skippedLinkBeforeUpsertHooks = append(skippedLinkBeforeUpsertHooks, skippedLinkHook)
	case boil.AfterInsertHook:
		skippedLinkAfterInsertHooks = append(skippedLinkAfterInsertHooks, skippedLinkHook)
	case boil.AfterSelectHook:
		skippedLinkAfterSelectHooks = append(skippedLinkAfterSelectHooks, skippedLinkHook)
	case boil.AfterUpdateHook:
		skippedLinkAfterUpdateHooks = append(skippedLinkAfterUpdateHooks, skippedLinkHook)
	case boil.AfterDeleteHook:
		skippedLinkAfterDeleteHooks = append(skippedLinkAfterDeleteHooks, skippedLinkHook)
	case boil.AfterUpsertHook:
		skippedLinkAfterUpsertHooks = append(skippedLinkAfterUpsertHooks, skippedLinkHook)
	}
}

// One returns a single skippedLink record from the query.
func (q skippedLinkQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SkippedLink, error) {
	o := &SkippedLink{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for skipped_links")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SkippedLink records from the query.
func (q skippedLinkQuery) All(ctx context.Context, exec boil.ContextExecutor) (SkippedLinkSlice, error) {
	var o []*SkippedLink

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SkippedLink slice")
	}

	if len(skippedLinkAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SkippedLink records in the query.
func (q skippedLinkQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count skipped_links rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q skippedLinkQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if skipped_links exists")
	}

	return count > 0, nil
}

// SkippedLinks retrieves all the records using an executor.
func SkippedLinks(mods ...qm.QueryMod) skippedLinkQuery {
	mods = append(mods, qm.From("\"skipped_links\""))
	return skippedLinkQuery{NewQuery(mods...)}
}

// FindSkippedLink retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSkippedLink(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*SkippedLink, error) {
	skippedLinkObj := &SkippedLink{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"skipped_links\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, skippedLinkObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from skipped_links")
	}

	if err = skippedLinkObj.doAfterSelectHooks(ctx, exec); err != nil {
		return skippedLinkObj, err
	}

	return skippedLinkObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SkippedLink) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no skipped_links provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(skippedLinkColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	skippedLinkInsertCacheMut.RLock()
	cache, cached := skippedLinkInsertCache[key]
	skippedLinkInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			skippedLinkAllColumns,
			skippedLinkColumnsWithDefault,
			skippedLinkColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(skippedLinkType, skippedLinkMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(skippedLinkType, skippedLinkMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"skipped_links\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"skipped_links\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into skipped_links")
	}

	if !cached {
		skippedLinkInsertCacheMut.Lock()
		skippedLinkInsertCache[key] = cache
		skippedLinkInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SkippedLink.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SkippedLink) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	skippedLinkUpdateCacheMut.RLock()
	cache, cached := skippedLinkUpdateCache[key]
	skippedLinkUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			skippedLinkAllColumns,
			skippedLinkPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update skipped_links, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"skipped_links\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, skippedLinkPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(skippedLinkType, skippedLinkMapping, append(wl, skippedLinkPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update skipped_links row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for skipped_links")
	}

	if !cached {
		skippedLinkUpdateCacheMut.Lock()
		skippedLinkUpdateCache[key] = cache
		skippedLinkUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q skippedLinkQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for skipped_links")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for skipped_links")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SkippedLinkSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), skippedLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"skipped_links\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, skippedLinkPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in skippedLink slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all skippedLink")
	}
	return rowsAff, nil
}

// Delete deletes a single SkippedLink record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SkippedLink) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no SkippedLink provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), skippedLinkPrimaryKeyMapping)
	sql := "DELETE FROM \"skipped_links\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from skipped_links")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for skipped_links")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q skippedLinkQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no skippedLinkQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from skipped_links")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for skipped_links")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SkippedLinkSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(skippedLinkBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), skippedLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"skipped_links\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, skippedLinkPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from skippedLink slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for skipped_links")
	}

	if len(skippedLinkAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SkippedLink) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSkippedLink(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SkippedLinkSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SkippedLinkSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), skippedLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"skipped_links\".* FROM \"skipped_links\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, skippedLinkPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SkippedLinkSlice")
	}

	*o = slice

	return nil
}

// SkippedLinkExists checks if the SkippedLink row exists.
func SkippedLinkExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"skipped_links\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if skipped_links exists")
	}

	return exists, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SkippedLink) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no skipped_links provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(skippedLinkColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	skippedLinkUpsertCacheMut.RLock()
	cache, cached := skippedLinkUpsertCache[key]
	skippedLinkUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			skippedLinkAllColumns,
			skippedLinkColumnsWithDefault,
			skippedLinkColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			skippedLinkAllColumns,
			skippedLinkPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert skipped_links, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(skippedLinkPrimaryKeyColumns))
			copy(conflict, skippedLinkPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"skipped_links\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(skippedLinkType, skippedLinkMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(skippedLinkType, skippedLinkMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert skipped_links")
	}

	if !cached {
		skippedLinkUpsertCacheMut.Lock()
		skippedLinkUpsertCache[key] = cache
		skippedLinkUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
