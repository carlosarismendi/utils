package infrastructure

import (
	"strconv"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/utilerror"
	"gorm.io/gorm"
)

func applyOffset(db *gorm.DB, value string) (*gorm.DB, int64, error) {
	err := domain.CheckEmptyValue("offset", value)
	if err != nil {
		return nil, 0, err
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rErr := utilerror.NewError(utilerror.WrongInputParameterError,
			`Invalid value for "offset". It must be a number.`).WithCause(err)
		return nil, 0, rErr
	}

	if num < 0 {
		rErr := utilerror.NewError(utilerror.WrongInputParameterError,
			`Invalid value for "offset". It must be greater or equal to 0.`)
		return nil, 0, rErr
	}

	return db.Offset(int(num)), num, nil
}
