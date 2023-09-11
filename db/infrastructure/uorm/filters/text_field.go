package filters

import (
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
	"gorm.io/gorm"
)

type TextFieldFilter struct {
	field string
}

func TextField(field string) *TextFieldFilter {
	return &TextFieldFilter{
		field: field,
	}
}

func (f *TextFieldFilter) Apply(db *gorm.DB, values []string) (*gorm.DB, error) {
	query, args, err := filters.ApplyTextField(f.field, values)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
