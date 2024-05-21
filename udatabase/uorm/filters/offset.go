package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/uerr"
	"gorm.io/gorm"
	"strconv"
)

type OffsetFilter struct {
}

func Offset() *OffsetFilter {
	return &OffsetFilter{}
}

func (f *OffsetFilter) Apply(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error) {
	return offset(db, rp, values...)
}

func (f *OffsetFilter) ValuedFilterFunc(values ...string) ValuedFilter {
	return func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error) {
		return offset(db, rp, values...)
	}
}

func offset(db *gorm.DB, rp *udatabase.ResourcePage, values ...string) (*gorm.DB, error) {
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
