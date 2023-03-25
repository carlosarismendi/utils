package domain

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/ansel1/merry"
)

func RemoveSpecialCharacters(str string) string {
	if str == "" {
		panic("Invalid value for field filter parameter. It can not be empty.")
	}
	return regexp.MustCompile(NotAlphaNumericRegex).ReplaceAllString(str, "")
}

func CheckEmptyValue(field, value string) error {
	if value == "" {
		errMsg := fmt.Sprintf("Invalid value for filter %q. It can not be empty.", field)
		return merry.New(errMsg).WithHTTPCode(http.StatusUnprocessableEntity)
	}
	return nil
}
