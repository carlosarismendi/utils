package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/ansel1/merry"
	"github.com/carlosarismendi/utils/db/domain"
	"gorm.io/gorm"
)

func applyLimit(db *gorm.DB, value string) (*gorm.DB, int64, error) {
	err := domain.CheckEmptyValue("limit", value)
	if err != nil {
		return nil, 0, err
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, 0, merry.New(`Invalid value for "limit". It must be a number.`).WithHTTPCode(http.StatusUnprocessableEntity)
	}

	if num < 1 {
		return nil, 0, merry.New(`Invalid value for "limit". It must be greater than 0.`).WithHTTPCode(http.StatusUnprocessableEntity)
	}

	return db.Limit(int(num)), num, nil
}
