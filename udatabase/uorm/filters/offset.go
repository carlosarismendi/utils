package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
	"strconv"
)

type OffsetFilter[T any] struct {
}

func Offset[T any]() *OffsetFilter[T] {
	return &OffsetFilter[T]{}
}

func (f *OffsetFilter[T]) Apply(db *gorm.DB, values []string, rp *udatabase.ResourcePage[T]) (*gorm.DB, error) {
	return f.offset(db, rp, values...)
}

func (f *OffsetFilter[T]) ValuedFilterFunc(values ...string) ValuedFilter[T] {
	return func(db *gorm.DB, rp *udatabase.ResourcePage[T]) (*gorm.DB, error) {
		return f.offset(db, rp, values...)
	}
}

func (f *OffsetFilter[T]) offset(db *gorm.DB, rp *udatabase.ResourcePage[T], values ...string) (*gorm.DB, error) {
	if len(values) == 0 || values[0] == "" {
		return db, nil
	}

	num, err := strconv.Atoi(values[0])
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
