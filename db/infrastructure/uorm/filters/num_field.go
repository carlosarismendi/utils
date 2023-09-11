package filters

import (
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
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

func (f *NumFieldFilter) Apply(db *gorm.DB, values []string) (*gorm.DB, error) {
	query, args, err := filters.ApplyNumField(f.field, values)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
