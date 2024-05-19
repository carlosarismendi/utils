package filters

import (
	"github.com/carlosarismendi/utils/udatabase"
	"gorm.io/gorm"
)

type Filter func(db *gorm.DB, values []string, rp *udatabase.ResourcePage) (*gorm.DB, error)

type ValuedFilter func(db *gorm.DB, rp *udatabase.ResourcePage) (*gorm.DB, error)
