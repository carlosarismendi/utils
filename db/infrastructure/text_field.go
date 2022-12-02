package infrastructure

import (
	"gorm.io/gorm"
)

type TextFieldFilter struct {
	field string
}

func TextField(field string) *TextFieldFilter {
	filtered := removeSpecialCharacters(field)
	return &TextFieldFilter{
		field: filtered + "=",
	}
}

func (f *TextFieldFilter) Apply(db *gorm.DB, value string) (*gorm.DB, error) {
	err := checkEmptyValue(f.field, value)
	if err != nil {
		return nil, err
	}

	return db.Where(f.field, value), nil
}
