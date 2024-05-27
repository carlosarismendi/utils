package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type BoolFieldFilter[T any] struct {
	field string
}

func BoolField[T any](field string) *BoolFieldFilter[T] {
	return &BoolFieldFilter[T]{
		field: field,
	}
}

func (f *BoolFieldFilter[T]) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
	return f.BoolField(db, values...)
}

func (f *BoolFieldFilter[T]) ValuedFilterFunc(values ...string) ValuedFilter[T] {
	return func(db *gorm.DB, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
		return f.BoolField(db, values...)
	}
}

func (f *BoolFieldFilter[T]) BoolField(db *gorm.DB, values ...string) (*gorm.DB, error) {
	query, args, err := filters.ApplyBoolField(f.field, values...)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
