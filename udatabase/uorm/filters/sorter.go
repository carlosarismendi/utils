package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type SorterFilter struct {
	allowedFields map[string]bool
}

func Sorter(allowedFields ...string) *SorterFilter {
	fields := make(map[string]bool)
	for _, f := range allowedFields {
		fields[f] = true
	}

	return &SorterFilter{
		allowedFields: fields,
	}
}

func (f *SorterFilter) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
	return f.sorter(db, values...)
}

func (f *SorterFilter) ValuedFilterFunc(values ...string) ValuedFilter {
	return func(db *gorm.DB, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		return f.sorter(db, values...)
	}
}

func (f *SorterFilter) sorter(db *gorm.DB, values ...string) (*gorm.DB, error) {
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
