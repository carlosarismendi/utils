package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

func Sorter(allowedFields ...string) Filter {
	fields := make(map[string]bool)
	for _, f := range allowedFields {
		fields[f] = true
	}

	return func(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		for _, v := range values {
			column, direction, err := filters.SortFieldAndDirection(fields, v)
			if err != nil {
				return nil, err
			}

			sort := column + " " + direction
			db = db.Order(sort)
		}

		return db, nil
	}
}
