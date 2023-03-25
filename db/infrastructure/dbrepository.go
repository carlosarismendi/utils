package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
	"github.com/carlosarismendi/utils/shared/utilerror"
	"gorm.io/gorm"
)

const defaultLimitValue = "10"

type ctxk string

const transactionName string = "dbtx"

type DBrepository struct {
	db      *gorm.DB
	filters map[string]filters.Filter
}

func NewDBRepository(dbHolder *DBHolder, filters map[string]filters.Filter) *DBrepository {
	return &DBrepository{
		db:      dbHolder.GetDBInstance(),
		filters: filters,
	}
}

func (r *DBrepository) Begin(ctx context.Context) (context.Context, error) {
	txFromCtx := ctx.Value(transactionName)
	if txFromCtx != nil {
		return ctx, nil
	}

	tx := r.db.Begin()

	if tx.Error != nil {
		tErr := utilerror.NewError(utilerror.GenericError, "Error begining transaction.").WithCause(tx.Error)
		return nil, tErr
	}

	ctx = context.WithValue(ctx, ctxk(transactionName), tx)
	return ctx, nil
}

func (r *DBrepository) Commit(ctx context.Context) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		tErr := utilerror.NewError(utilerror.GenericError, "Missing transaction when doing Commit.")
		return tErr
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Commit().Error
}

func (r *DBrepository) Rollback(ctx context.Context) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		tErr := utilerror.NewError(utilerror.GenericError, "Missing transaction when doing Rollback.")
		return tErr
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Rollback().Error
}

func (r *DBrepository) Save(ctx context.Context, value interface{}) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		tErr := utilerror.NewError(utilerror.GenericError, "Missing transaction when doing Save.")
		return tErr
	}
	tx := txFromCtx.(*gorm.DB)
	err := tx.Create(value).Error

	return r.HandleSaveOrUpdateError(err)
}

func (r *DBrepository) FindByID(ctx context.Context, id string, dest interface{}) error {
	err := r.db.Where("id = ?", id).First(dest).Error

	var tErr error
	if err != nil {
		if r.IsResourceNotFound(err) {
			tErr = utilerror.NewError(utilerror.ResourceNotFoundError, "Resource not found.")
		} else {
			tErr = utilerror.NewError(utilerror.GenericError, "Error finding resource by id.").WithCause(err)
		}
	}

	return tErr
}

func (r *DBrepository) Find(ctx context.Context, v url.Values, dst interface{}) (*domain.ResourcePage, error) {
	db, limit, err := r.applyLimit(r.db, &v)
	if err != nil {
		return nil, err
	}

	for key, values := range v {
		if len(values) == 0 {
			continue
		}

		filter, ok := r.filters[key]
		if !ok {
			rErr := utilerror.NewError(utilerror.WrongInputParameterError, fmt.Sprintf("Invalid filter %q.", key))
			return nil, rErr
		}

		var err error
		db, err = filter.Apply(db, values[0])
		if err != nil {
			return nil, err
		}
	}

	result := db.Find(dst)
	if result.Error != nil {
		rErr := utilerror.NewError(utilerror.GenericError, "Error finding resources.").WithCause(result.Error)
		return nil, rErr
	}

	rp := &domain.ResourcePage{
		Total:     result.RowsAffected,
		Limit:     int64(limit),
		Resources: dst,
	}

	return rp, nil
}

func (r *DBrepository) IsResourceNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *DBrepository) HandleSaveOrUpdateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return utilerror.NewError(utilerror.ResourceAlreadyExistsError, "Resource already exists.").WithCause(err)
	} else {
		return utilerror.NewError(utilerror.GenericError, "Error saving or updating resource.").WithCause(err)
	}
}

func (r *DBrepository) applyLimit(db *gorm.DB, v *url.Values) (*gorm.DB, int64, error) {
	values, ok := (*v)["limit"]
	var limit string
	if ok {
		limit = values[0]
	} else {
		limit = defaultLimitValue
	}

	v.Del("limit")
	return applyLimit(db, limit)
}
