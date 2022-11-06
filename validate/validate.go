package validate

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New()

func Validate(v interface{}) error {
	err := validate.Struct(v)
	if err != nil {
		var mErr merry.Error
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.ActualTag()
			if err.Param() != "" {
				tag = fmt.Sprintf("%s=%s", tag, err.Param())
			}

			errMsg := fmt.Sprintf("Invalid field %s: the value must be '%s'. The value received is '%v'", err.Field(), tag, err.Value())
			if mErr == nil {
				mErr = merry.New(errMsg)
			} else {
				mErr = merry.Append(mErr, fmt.Sprintf("\n%s", errMsg))
			}
		}

		return mErr.WithHTTPCode(http.StatusUnprocessableEntity)
	}

	return nil
}
