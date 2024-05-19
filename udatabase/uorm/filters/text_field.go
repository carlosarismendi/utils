package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

func TextField(field string) Filter {
	return func(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		query, args, err := filters.ApplyTextField(field, values...)
		if err != nil {
			return nil, err
		}

		return db.Where(query, args...), nil
	}
}

func TextFieldWithValue(field string, values ...string) ValuedFilter {
	return func(db *gorm.DB, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		query, args, err := filters.ApplyTextField(field, values...)
		if err != nil {
			return nil, err
		}

		return db.Where(query, args...), nil
	}
}
