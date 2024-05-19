package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
	"strconv"
)

func Limit() Filter {
	return func(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		if len(values) < 1 {
			rp.Limit = filters.DefaultLimit
			return db.Limit(filters.DefaultLimit), nil
		}

		return limit(db, values[0], rp)
	}
}

func LimitWithValue(value string) ValuedFilter {
	return func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		if value == "" {
			rp.Limit = filters.DefaultLimit
			return db.Limit(filters.DefaultLimit), nil
		}

		return limit(db, value, rp)
	}
}

func limit(db *gorm.DB, value string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
	num, err := strconv.Atoi(value)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be a number.`).WithCause(err)
		return nil, rErr
	}

	if num < 1 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be greater than 0.`)
		return nil, rErr
	}

	rp.Limit = int64(num)
	return db.Limit(num), nil
}
