package uorm

import (
	"strconv"

	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
)

func applyLimit(db *gorm.DB, value string) (*gorm.DB, int64, error) {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be a number.`).WithCause(err)
		return nil, 0, rErr
	}

	if num < 1 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be greater than 0.`)
		return nil, 0, rErr
	}

	return db.Limit(int(num)), num, nil
}
