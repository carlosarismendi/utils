package infrastructure

import (
	"gorm.io/gorm"
)

type Sorter interface {
	Apply(db *gorm.DB, value string) (*gorm.DB, error)
}
