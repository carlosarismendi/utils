package filters

import (
	"fmt"

	"github.com/carlosarismendi/utils/db/domain"
	"gorm.io/gorm"
)

type SortFilter struct {
	field string
}

func Sorter() *SortFilter {
	return &SortFilter{}
}

func (f *SortFilter) Apply(db *gorm.DB, value string) (*gorm.DB, error) {
	filteredValue := domain.RemoveSpecialCharacters(value)
	err := domain.CheckEmptyValue(f.field, filteredValue)
	if err != nil {
		return nil, err
	}

	var field string
	var direction string
	if filteredValue[0] == '-' {
		field = filteredValue[1 : len(filteredValue)-1]
		direction = "desc"
	} else {
		field = filteredValue
		direction = "asc"
	}

	sort := fmt.Sprintf("%s %s", field, direction)

	return db.Order(sort), nil
}
