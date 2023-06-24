package filters

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(db *gorm.DB, values []string) (*gorm.DB, error)
}
