package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type NumFieldFilter struct {
	field string
}

func NumField(field string) *NumFieldFilter {
	return &NumFieldFilter{
		field: field,
	}
}

func (f *NumFieldFilter) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
	return f.numField(db, values...)
}

func (f *NumFieldFilter) ValuedFilterFunc(values ...string) ValuedFilter {
	return func(db *gorm.DB, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		return f.numField(db, values...)
	}
}

func (f *NumFieldFilter) numField(db *gorm.DB, values ...string) (*gorm.DB, error) {
	query, args, err := filters.ApplyNumField(f.field, values...)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
