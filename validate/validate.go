package validate

import (
	"fmt"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New()

func Validate(v interface{}) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var sb strings.Builder
	for i, err := range err.(validator.ValidationErrors) {
		if i > 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString("Invalid field ")
		sb.WriteString(err.Field())
		sb.WriteString(": the value must be '")

		tag := err.ActualTag()
		sb.WriteString(tag)
		if err.Param() != "" {
			sb.WriteByte('=')
			sb.WriteString(err.Param())
		}

		sb.WriteString("'. The value received is ")
		sb.WriteString(fmt.Sprintf("'%v'.", err.Value()))
	}

	return uerr.NewError(uerr.WrongInputParameterError, sb.String())
}
