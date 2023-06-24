package filters

import (
	"fmt"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/utilerror"
	"gorm.io/gorm"
)

type SortFilter struct {
	field string

	allowedFields map[string]bool
}

func Sorter(fields ...string) *SortFilter {
	allowedFields := make(map[string]bool)
	for _, f := range fields {
		allowedFields[f] = true
	}
	return &SortFilter{
		allowedFields: allowedFields,
	}
}

func (f *SortFilter) Apply(db *gorm.DB, values []string) (*gorm.DB, error) {
	newDB := db
	for _, v := range values {
		filteredValue := domain.RemoveSpecialCharacters(v)
		err := domain.CheckEmptyValue(f.field, filteredValue)
		if err != nil {
			return nil, err
		}

		var field string
		var direction string
		if filteredValue[0] == '-' {
			field = filteredValue[1:]
			direction = "desc"
		} else {
			field = filteredValue
			direction = "asc"
		}

		if _, ok := f.allowedFields[field]; !ok {
			return nil, utilerror.NewError(utilerror.WrongInputParameterError, fmt.Sprintf("Invalid sort field %q.", field))
		}

		sort := fmt.Sprintf("%s %s", field, direction)
		newDB = newDB.Order(sort)
	}

	return newDB, nil
}
