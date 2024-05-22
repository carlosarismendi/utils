package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type SorterFilter[T any] struct {
	allowedFields map[string]bool
}

func Sorter[T any](allowedFields ...string) *SorterFilter[T] {
	fields := make(map[string]bool)
	for _, f := range allowedFields {
		fields[f] = true
	}

	return &SorterFilter[T]{
		allowedFields: fields,
	}
}

func (f *SorterFilter[T]) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
	return f.sorter(db, values...)
}

func (f *SorterFilter[T]) ValuedFilterFunc(values ...string) ValuedFilter[T] {
	return func(db *gorm.DB, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
		return f.sorter(db, values...)
	}
}

func (f *SorterFilter[T]) sorter(db *gorm.DB, values ...string) (*gorm.DB, error) {
	for _, v := range values {
		column, direction, err := filters.SortFieldAndDirection(f.allowedFields, v)
		if err != nil {
			return nil, err
		}

		sort := column + " " + direction
		db = db.Order(sort)
	}

	return db, nil
}
