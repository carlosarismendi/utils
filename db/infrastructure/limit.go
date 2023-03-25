package infrastructure

import (
	"strconv"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/shared/utilerror"
	"gorm.io/gorm"
)

func applyLimit(db *gorm.DB, value string) (*gorm.DB, int64, error) {
	err := domain.CheckEmptyValue("limit", value)
	if err != nil {
		return nil, 0, err
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rErr := utilerror.NewError(utilerror.WrongInputParameterError,
			`Invalid value for "limit". It must be a number.`).WithCause(err)
		return nil, 0, rErr
	}

	if num < 1 {
		rErr := utilerror.NewError(utilerror.WrongInputParameterError, `Invalid value for "limit". It must be greater than 0.`)
		return nil, 0, rErr
	}

	return db.Limit(int(num)), num, nil
}
