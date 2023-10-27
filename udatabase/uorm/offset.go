package uorm

import (
	"strconv"

	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
)

func applyOffset(db *gorm.DB, value string) (*gorm.DB, int64, error) {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be a number.`).WithCause(err)
		return nil, 0, rErr
	}

	if num < 0 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be greater or equal to 0.`)
		return nil, 0, rErr
	}

	return db.Offset(int(num)), num, nil
}
