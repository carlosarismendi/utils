package usql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	usqlFilters "github.com/carlosarismendi/utils/udatabase/usql/filters"
	"github.com/carlosarismendi/utils/uerr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ctxk string

const transactionName string = "dbtx"

// DBrepository is built on top of sqlx to provide easier transaction management as
// well as methods for error handling.
type DBrepository[T any] struct {
	db *DBHolder

	filters map[string]usqlFilters.Filter
	sorters map[string]usqlFilters.Sorter
}

// NewDBRepository returns a DBrepository.
// requires a that map will be used in the method Find(context.Context, url.values) to use the filters
// and sorters provided in the url.values{} parameter. In case the url.values contains a filter
// that it is not in the filters map, it will return an error.
func NewDBRepository[T any](dbHolder *DBHolder, filtersMap map[string]usqlFilters.Filter,
	sorters map[string]usqlFilters.Sorter) *DBrepository[T] {
	return &DBrepository[T]{
		db:      dbHolder,
		filters: filtersMap,
		sorters: sorters,
	}
}

// Begin opens a new transaction.
// NOTE: Nested transactions not supported.
func (r *DBrepository[T]) Begin(ctx context.Context) (context.Context, error) {
	txFromCtx := ctx.Value(transactionName)
	if txFromCtx != nil {
		return ctx, nil
	}

	tx, err := r.db.db.BeginTxx(ctx, nil)
	if err != nil {
		tErr := uerr.NewError(uerr.GenericError, "Error beginning transaction.").WithCause(err)
		return nil, tErr
	}

	ctx = context.WithValue(ctx, ctxk(transactionName), tx)
	return ctx, nil
}

// Commit closes and confirms the current transaction.
func (r *DBrepository[T]) Commit(ctx context.Context) error {
	tx := r.GetTransaction(ctx)
	if tx == nil {
		tErr := uerr.NewError(uerr.GenericError, "Missing transaction when doing Commit.")
		return tErr
	}
	return tx.Commit()
}

// Rollback cancels the current transaction.
func (r *DBrepository[T]) Rollback(ctx context.Context) {
	tx := r.GetTransaction(ctx)
	if tx == nil {
		return
	}
	_ = tx.Rollback()
}

// IsResourceNotFound in case of running SELECT queries using *sqlx.DB/*sqlx.Tx, this method
// provides an easy way of checking if the error returned is a NotFound or other type.
func (r *DBrepository[T]) HandleSearchError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "no rows in result") {
		return uerr.NewError(uerr.ResourceNotFoundError, "Resource not found.").WithCause(err)
	}

	return uerr.NewError(uerr.GenericError, "Error searching resource.").WithCause(err)
}

// HandleSaveOrUpdateError in case of running an INSERT/UPDATE query, this method provides
// an easy way of checking if the returned error is nil or if it violates a PRIMARY KEY/UNIQUE constraint.
func (r *DBrepository[T]) HandleSaveOrUpdateError(res sql.Result, err error) error {
	if err == nil {
		n, rErr := res.RowsAffected()
		if rErr != nil {
			return uerr.NewError(uerr.GenericError, "Error saving or updating resource.").WithCause(rErr)
		}

		if n <= 0 {
			return uerr.NewError(uerr.ResourceNotFoundError, "Resource(s) not found.")
		}

		return nil
	}

	if pqErr, ok := err.(*pq.Error); ok {
		if rErr, ok := udatabase.PqErrors[pqErr.Code.Name()]; ok {
			return rErr.WithCause(err)
		}
	}

	return uerr.NewError(uerr.GenericError, "Error saving or updating resource.").WithCause(err)
}

func (r *DBrepository[T]) GetDBInstance() *sqlx.DB {
	return r.db.GetDBInstance()
}

func (r *DBrepository[T]) GetTransaction(ctx context.Context) *sqlx.Tx {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		return nil
	}
	return txFromCtx.(*sqlx.Tx)
}

type Querier interface {
	Rebind(query string) string
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

func (r *DBrepository[T]) GetContext(ctx context.Context, db Querier, dst T, query string, v url.Values) (T, error) {
	v.Del("limit")
	v.Add("limit", "1")
	query, args, _, _, err := r.ApplyFilters(db, query, v)
	if err != nil {
		return dst, err
	}
	err = db.GetContext(ctx, dst, query, args...)
	return dst, r.HandleSearchError(err)
}

func (r *DBrepository[T]) SelectContext(ctx context.Context, db Querier, query string,
	v url.Values) (rp *udatabase.ResourcePage[T], err error) {
	query, args, limit, offset, err := r.ApplyFilters(db, query, v)
	if err != nil {
		return nil, err
	}

	var dst []T
	err = db.SelectContext(ctx, &dst, query, args...)
	if err != nil {
		return nil, r.HandleSearchError(err)
	}

	rp = &udatabase.ResourcePage[T]{
		Total:     int64(len(dst)),
		Limit:     limit,
		Offset:    offset,
		Resources: dst,
	}

	return rp, nil
}

func (r *DBrepository[T]) ApplyFilters(db Querier, query string, v url.Values) (queryResult string, args []any,
	limit, offset int64, err error) {
	limitQ, limit, err := r.applyLimit(v)
	if err != nil {
		return "", nil, 0, 0, err
	}
	offsetQ, offset, err := r.applyOffset(v)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// Apply filters
	conds, args, unknownFilters, err := r.applyFilters(v)
	if err != nil {
		return "", nil, 0, 0, err
	}

	err = r.processUnknownFilters(unknownFilters)
	if err != nil {
		return "", nil, 0, 0, err
	}

	var sb strings.Builder
	sb.Grow(len(query) + len(conds) + len(limitQ) + len(offsetQ))
	sb.WriteString(query)
	sb.WriteString(conds)
	sb.WriteString(limitQ)

	if offset > 0 {
		sb.WriteString(offsetQ)
	}

	query = sb.String()
	if len(args) > 0 {
		query = db.Rebind(query)
	}
	return query, args, limit, offset, nil
}

func (r *DBrepository[T]) applyFilters(v url.Values) (conds string, args []any, unknown []string,
	err error) {
	args = make([]any, 0, len(v))
	var sbConds, sbSorts strings.Builder
	var cSep, sSep string
	for key, values := range v {
		if len(values) == 0 {
			continue
		}

		filter, ok := r.filters[key]
		// If filter not found => try to apply sorter
		if !ok {
			sorter, ok := r.sorters[key]
			if !ok {
				unknown = append(unknown, key)
				continue
			}

			sort, err := sorter.Apply(values)
			if err != nil {
				return "", nil, nil, err
			}

			sbSorts.WriteString(sSep)
			sbSorts.WriteString(sort)
			sSep = ", "
			continue
		}

		cond, fArgs, err := filter.Apply(values)
		if err != nil {
			return "", nil, nil, err
		}

		sbConds.WriteString(cSep)
		sbConds.WriteString(cond)
		cSep = " AND "

		args = append(args, fArgs...)
	}

	var sb strings.Builder
	sb.Grow(7 + sbConds.Len() + 10 + sbSorts.Len())
	if cSep != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(sbConds.String())
	}

	if sSep != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(sbSorts.String())
	}

	return sb.String(), args, unknown, nil
}

func (r *DBrepository[T]) applyLimit(v url.Values) (limitQ string, limitNum int64, rErr error) {
	values, ok := v["limit"]
	var limit string
	if ok {
		limit = values[0]
		v.Del("limit")
	} else {
		return filters.DefaultLimitStr, filters.DefaultLimit, nil
	}

	return filters.ApplyLimit("", limit)
}

func (r *DBrepository[T]) applyOffset(v url.Values) (offsetQ string, offsetNum int64, rErr error) {
	values, ok := v["offset"]
	var offset string
	if !ok {
		return "", 0, nil
	}

	offset = values[0]
	v.Del("offset")
	return filters.ApplyOffset("", offset)
}

func (r *DBrepository[T]) processUnknownFilters(unknown []string) error {
	if len(unknown) == 0 {
		return nil
	}

	var msg, sep string
	for _, value := range unknown {
		msg += fmt.Sprintf("%s Invalid filter %q", sep, value)
		sep = "; "
	}

	return uerr.NewError(uerr.WrongInputParameterError, msg)
}
