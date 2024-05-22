package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type NumFieldFilter[T any] struct {
	field string
}

func NumField[T any](field string) *NumFieldFilter[T] {
	return &NumFieldFilter[T]{
		field: field,
	}
}

func (f *NumFieldFilter[T]) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
	return f.numField(db, values...)
}

func (f *NumFieldFilter[T]) ValuedFilterFunc(values ...string) ValuedFilter[T] {
	return func(db *gorm.DB, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
		return f.numField(db, values...)
	}
}

func (f *NumFieldFilter[T]) numField(db *gorm.DB, values ...string) (*gorm.DB, error) {
	query, args, err := filters.ApplyNumField(f.field, values...)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
