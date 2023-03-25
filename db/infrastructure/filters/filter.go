package filters

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(db *gorm.DB, value string) (*gorm.DB, error)
}
