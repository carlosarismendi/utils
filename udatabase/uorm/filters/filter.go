package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"gorm.io/gorm"
)

type ValuedFilter func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error)

type Filter interface {
	Apply(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error)
	ValuedFilterFunc(values ...string) ValuedFilter
}
