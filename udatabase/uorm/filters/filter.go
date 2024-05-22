package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"gorm.io/gorm"
)

type ValuedFilter[T any] func(db *gorm.DB, rp *udatabase.ResourcePage[T]) (*gorm.DB, error)

type Filter[T any] interface {
	Apply(db *gorm.DB, values []string, rp *udatabase.ResourcePage[T]) (*gorm.DB, error)
	ValuedFilterFunc(values ...string) ValuedFilter[T]
}
