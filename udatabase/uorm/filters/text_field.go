package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"gorm.io/gorm"
)

type TextFieldFilter[T any] struct {
	field string
}

func TextField[T any](field string) *TextFieldFilter[T] {
	return &TextFieldFilter[T]{
		field: field,
	}
}

func (f *TextFieldFilter[T]) Apply(db *gorm.DB, values []string, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
	return f.textField(db, values)
}

func (f *TextFieldFilter[T]) ValuedFilterFunc(values ...string) ValuedFilter[T] {
	return func(db *gorm.DB, _ *udatabase.ResourcePage[T]) (*gorm.DB, error) {
		return f.textField(db, values)
	}
}

func (f *TextFieldFilter[T]) textField(db *gorm.DB, values []string) (*gorm.DB, error) {
	query, args, err := filters.ApplyTextField(f.field, values...)
	if err != nil {
		return nil, err
	}

	return db.Where(query, args...), nil
}
