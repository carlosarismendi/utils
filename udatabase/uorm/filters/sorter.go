package filters

import (
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type SortFilter struct {
	allowedFields map[string]bool
}

func Sorter(allowByFields ...string) *SortFilter {
	allowedFields := make(map[string]bool)
	for _, f := range allowByFields {
		allowedFields[f] = true
	}
	return &SortFilter{
		allowedFields: allowedFields,
	}
}

func (f *SortFilter) Apply(db *gorm.DB, values []string) (*gorm.DB, error) {
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
