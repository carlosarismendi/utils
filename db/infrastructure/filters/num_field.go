package filters

import (
	"fmt"
	"strconv"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/shared/utilerror"
	"gorm.io/gorm"
)

type NumFieldFilter struct {
	field string
}

func NumField(field string) *NumFieldFilter {
	filtered := domain.RemoveSpecialCharacters(field)
	return &NumFieldFilter{
		field: filtered + " = ?",
	}
}

func (f *NumFieldFilter) Apply(db *gorm.DB, value string) (*gorm.DB, error) {
	err := domain.CheckEmptyValue(f.field, value)
	if err != nil {
		return nil, err
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid value for filter %q. It must be a number.", f.field)
		return nil, utilerror.NewError(utilerror.WrongInputParameterError, errMsg).WithCause(err)
	}
	return db.Where(f.field, num), nil
}
