package domain

import (
	"fmt"
	"regexp"

	"github.com/carlosarismendi/utils/shared/utilerror"
)

func RemoveSpecialCharacters(str string) string {
	if str == "" {
		panic("Invalid value for field filter parameter. It can not be empty.")
	}
	return regexp.MustCompile(NotAlphaNumericRegex).ReplaceAllString(str, "")
}

func CheckEmptyValue(field, value string) error {
	if value == "" {
		return utilerror.NewError(utilerror.WrongInputParameterError,
			fmt.Sprintf("Invalid value for filter %q. It can not be empty.", field))
	}
	return nil
}
