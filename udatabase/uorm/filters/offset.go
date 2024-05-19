package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
	"strconv"
)

func Offset() Filter {
	return func(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		if len(values) == 0 {
			return db, nil
		}

		return offset(db, values[0], rp)
	}
}

func OffsetWithValue(value string) ValuedFilter {
	return func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		if value == "" {
			return db, nil
		}

		return offset(db, value, rp)
	}
}

func offset(db *gorm.DB, value string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
	num, err := strconv.Atoi(value)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be a number.`).WithCause(err)
		return nil, rErr
	}

	if num < 0 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be greater or equal to 0.`)
		return nil, rErr
	}

	rp.Offset = int64(num)
	return db.Offset(num), nil
}
