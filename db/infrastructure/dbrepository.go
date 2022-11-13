package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type ctxk string

const transactionName string = "dbtx"

type DBrepository struct {
	db *gorm.DB
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
