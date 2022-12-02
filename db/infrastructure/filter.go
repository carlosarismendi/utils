package infrastructure

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/ansel1/merry"
	"github.com/carlosarismendi/dddhelper/db/domain"
	"gorm.io/gorm"
)

type Filter interface {
	Apply(db *gorm.DB, value string) (*gorm.DB, error)
}

func removeSpecialCharacters(str string) string {
	if str == "" {
		panic("Invalid value for field filter parameter. It can not be empty.")
	}
	return regexp.MustCompile(domain.AlphaNumericRegex).ReplaceAllString(str, "")
}

func checkEmptyValue(field, value string) error {
	if value == "" {
		errMsg := fmt.Sprintf("Invalid value for filter %q. It can not be empty.", field)
		return merry.New(errMsg).WithHTTPCode(http.StatusUnprocessableEntity)
	}
	return nil
}
