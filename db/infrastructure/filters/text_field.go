package filters

import (
	"github.com/carlosarismendi/utils/db/domain"
	"gorm.io/gorm"
)

type TextFieldFilter struct {
	field string
}

func TextField(field string) *TextFieldFilter {
	filtered := domain.RemoveSpecialCharacters(field)
	return &TextFieldFilter{
		field: filtered + " = ?",
	}
}

func (f *TextFieldFilter) Apply(db *gorm.DB, values []string) (*gorm.DB, error) {
	err := domain.CheckEmptyValue(f.field, values[0])
	if err != nil {
		return nil, err
	}

	return db.Where(f.field, values[0]), nil
}
