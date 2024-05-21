package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
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

func (f *TextFieldFilter) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage) (*gorm.DB, error) {
	return f.textField(db, values)
}

func (f *TextFieldFilter) ValuedFilterFunc(values ...string) ValuedFilter {
	return func(db *gorm.DB, _ *udatabase.ResourcePage) (*gorm.DB, error) {
		return f.textField(db, values)
	}
}

func (f *TextFieldFilter) textField(db *gorm.DB, values []string) (*gorm.DB, error) {
	query, args, err := filters.ApplyTextField(f.field, values...)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
