package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

func NumField(field string) Filter {
	return func(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		query, args, err := filters.ApplyNumField(field, values...)
		if err != nil {
			return nil, err
		}

		return db.Where(query, args...), nil
	}
}

func NumFieldWithValue(field, value string) ValuedFilter {
	return func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		query, args, err := filters.ApplyNumField(field, value)
		if err != nil {
			return nil, err
		}

		return db.Where(query, args...), nil
	}
}
