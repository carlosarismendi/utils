package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ansel1/merry"
	"github.com/carlosarismendi/dddhelper/db/domain"
	"gorm.io/gorm"
)

const defaultLimitValue = "10"

type ctxk string

const transactionName string = "dbtx"

type DBrepository struct {
	db      *gorm.DB
	filters map[string]Filter
}

func NewDBRepository(dbHolder *DBHolder) *DBrepository {
	return &DBrepository{
		db: dbHolder.GetDBInstance(),
	}
}

func (r *DBrepository) Begin(ctx context.Context) (context.Context, error) {
	txFromCtx := ctx.Value(transactionName)
	if txFromCtx != nil {
		return ctx, nil
	}

	tx := r.db.Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	ctx = context.WithValue(ctx, ctxk(transactionName), tx)
	return ctx, nil
}

func (r *DBrepository) Commit(ctx context.Context) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		return errors.New("NIL TX in Commit")
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Commit().Error
}

func (r *DBrepository) Rollback(ctx context.Context) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		return errors.New("NIL TX in Rollback")
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Rollback().Error
}

func (r *DBrepository) Save(ctx context.Context, value interface{}) error {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		return errors.New("NIL TX in Save")
	}
	tx := txFromCtx.(*gorm.DB)
	return tx.Create(value).Error
}

func (r *DBrepository) FindByID(ctx context.Context, id string, dest interface{}) error {
	return r.db.Where("id = ?", id).First(dest).Error
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
			return nil, merry.New(fmt.Sprintf("Invalid filter %q.", key)).WithHTTPCode(http.StatusUnprocessableEntity)
		}

		var err error
		db, err = filter.Apply(db, values[0])
		if err != nil {
			return nil, err
		}
	}

	result := db.Find(dst)
	if result.Error != nil {
		return nil, result.Error
	}

	rp := &domain.ResourcePage{
		Total:     result.RowsAffected,
		Limit:     int64(limit),
		Resources: r,
	}

	return rp, nil
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
