package uorm

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	uormFilters "github.com/carlosarismendi/utils/udatabase/uorm/filters"
	"github.com/carlosarismendi/utils/uerr"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ctxk string

const transactionName string = "dbtx"

// DBrepository is built on top of GORM to provide easier transaction management as
// well as methods like Save or Find.
type DBrepository struct {
	db      *DBHolder
	filters map[string]uormFilters.Filter
}

// NewDBRepository returns a DBrepository.
// requires a that map will be used in the method Find(context.Context, url.values) to use the filters
// and sorters provided in the url.values{} parameter. In case the url.values contains a filter
// that it is not in the filters map, it will return an error.
func NewDBRepository(dbHolder *DBHolder, filtersMap map[string]uormFilters.Filter) *DBrepository {
	return &DBrepository{
		db:      dbHolder,
		filters: filtersMap,
	}
}

// Begin opens a new transaction.
// NOTE: Nested transactions not supported.
func (r *DBrepository) Begin(ctx context.Context) (context.Context, error) {
	txFromCtx := ctx.Value(transactionName)
	if txFromCtx != nil {
		return ctx, nil
	}

	tx := r.db.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		tErr := uerr.NewError(uerr.GenericError, "Error beginning transaction.").WithCause(tx.Error)
		return nil, tErr
	}

	ctx = context.WithValue(ctx, ctxk(transactionName), tx)
	return ctx, nil
}

// Commit closes and confirms the current transaction.
func (r *DBrepository) Commit(ctx context.Context) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		tErr := uerr.NewError(uerr.GenericError, "Missing transaction when doing Commit.")
		return tErr
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Commit().Error
}

// Rollback cancels the current transaction.
func (r *DBrepository) Rollback(ctx context.Context) {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		return
	}
	tx := txFromCtx.(*gorm.DB)
	_ = tx.Rollback().Error
}

// Save is a combination function. If save value does not contain primary key,
// it will execute Create, otherwise it will execute Update (with all fields).
func (r *DBrepository) Save(ctx context.Context, value interface{}) error {
	db := r.GetDBInstance(ctx)
	err := db.Save(value).Error

	return r.HandleSaveOrUpdateError(err)
}

// Create is a function that creates the resource in the database.
func (r *DBrepository) Create(ctx context.Context, value interface{}) error {
	db := r.GetDBInstance(ctx)
	err := db.Create(value).Error

	return r.HandleSaveOrUpdateError(err)
}

// FindByID returns the resource found in the variable dst.
// Usage:
//
//	type Resource struct {...}
//	var obj Resource
//	repository.FindByID(ctx, "an_ID", &obj)
func (r *DBrepository) FindByID(ctx context.Context, id string, dest interface{}) error {
	db := r.GetDBInstance(ctx)
	err := db.Where("id = ?", id).First(dest).Error

	var tErr error
	if err != nil {
		if r.IsResourceNotFound(err) {
			tErr = uerr.NewError(uerr.ResourceNotFoundError, "Resource not found.")
		} else {
			tErr = uerr.NewError(uerr.GenericError, "Error finding resource by id.").WithCause(err)
		}
	}

	return tErr
}

// Find returns a list of elements matching the provided filters.
// Usage:
//
//	type Resource struct {...}
//	var list []*Resource
//	repository.FindByID(ctx, url.values{}, list)
//
// It is necessary to pass the list parameter so
// internally can infer the type and table to use to
// request the data.
// resourcePage is of type:
//
//	type ResourcePage struct {
//		   Total  int64 `json:"total"`
//		   Limit  int64 `json:"limit"`
//		   Offset int64 `json:"offset"`
//
//	    // Resource will be a pointer to the type passed as
//	    // dst parameter in Find method. In this example,
//	    // *[]*Resource.
//	    Resources interface{} `json:"resources"`
//	}
//
// Filter:
//
//	v := url.values{}
//	v.Add("field", "value to use to filter")
//	v.Add("sort", "field")  // sort in ascending order
//	v.Add("sort", "-field") // sort in descending order
func (r *DBrepository) Find(ctx context.Context, v url.Values, dst interface{}) (*udatabase.ResourcePage, error) {
	db := r.GetDBInstance(ctx)
	db, limit, err := r.applyLimit(db, v)
	if err != nil {
		return nil, err
	}

	db, offset, err := r.applyOffset(db, v)
	if err != nil {
		return nil, err
	}

	for key, values := range v {
		if len(values) == 0 {
			continue
		}

		filter, ok := r.filters[key]
		if !ok {
			rErr := uerr.NewError(uerr.WrongInputParameterError, fmt.Sprintf("Invalid filter %q.", key))
			return nil, rErr
		}

		var err error
		db, err = filter.Apply(db, values)
		if err != nil {
			return nil, err
		}
	}

	result := db.Find(dst)
	if result.Error != nil {
		rErr := uerr.NewError(uerr.GenericError, "Error finding resources.").WithCause(result.Error)
		return nil, rErr
	}

	rp := &udatabase.ResourcePage{
		Total:     result.RowsAffected,
		Limit:     limit,
		Offset:    offset,
		Resources: dst,
	}

	return rp, nil
}

// IsResourceNotFound in case of running custom SELECT queries using *gorm.DB, this method
// provides an easy way of checking if the error returned is a NotFound or other type.
func (r *DBrepository) IsResourceNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// HandleSaveOrUpdateError in case of running an INSERT/UPDATE query, this method provides
// an easy way of checking if the returned error is nil or if it violates a PRIMARY KEY/UNIQUE constraint.
func (r *DBrepository) HandleSaveOrUpdateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return uerr.NewError(uerr.ResourceAlreadyExistsError, "Resource already exists.").WithCause(err)
	}

	if pqErr, ok := err.(*pq.Error); ok {
		if rErr, ok := udatabase.PqErrors[pqErr.Code.Name()]; ok {
			return rErr.WithCause(err)
		}
	} else if pgErr, ok := err.(*pgconn.PgError); ok {
		pqErr := pq.Error{Code: pq.ErrorCode(pgErr.Code)}
		if rErr, ok := udatabase.PqErrors[pqErr.Code.Name()]; ok {
			return rErr.WithCause(err)
		}
	}

	return uerr.NewError(uerr.GenericError, "Error saving or updating resource.").WithCause(err)
}

func (r *DBrepository) applyLimit(db *gorm.DB, v url.Values) (*gorm.DB, int64, error) {
	values, ok := v["limit"]
	var limit string
	if ok {
		limit = values[0]
		v.Del("limit")
	} else {
		return db.Limit(filters.DefaultLimit), filters.DefaultLimit, nil
	}

	return applyLimit(db, limit)
}

func (r *DBrepository) applyOffset(db *gorm.DB, v url.Values) (*gorm.DB, int64, error) {
	values, ok := v["offset"]
	if !ok {
		return db, 0, nil
	}

	offset := values[0]

	v.Del("offset")
	return applyOffset(db, offset)
}

func (r *DBrepository) GetDBInstance(ctx context.Context) *gorm.DB {
	return r.db.GetDBInstance(ctx)
}
