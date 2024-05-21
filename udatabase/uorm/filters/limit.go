package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
	"strconv"
)

type LimitFilter struct {
	defaultValue int
}

func Limit(defaultValue int) *LimitFilter {
	if defaultValue < 1 {
		panic("Limit defaultValue must be greater than 0")
	}
	return &LimitFilter{
		defaultValue: defaultValue,
	}
}

func (f *LimitFilter) Apply(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
	return f.limit(db, rp, values...)
}

func (f *LimitFilter) ValuedFilterFunc(values ...string) ValuedFilter {
	return func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		return f.limit(db, rp, values...)
	}
}

func (f *LimitFilter) limit(db *gorm.DB, rp *udatabase.ResourcePage, values ...string) (*gorm.DB, error) {
	if len(values) < 1 || values[0] == "" {
		rp.Limit = int64(f.defaultValue)
		return db.Limit(f.defaultValue), nil
	}

	num, err := strconv.Atoi(values[0])
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
