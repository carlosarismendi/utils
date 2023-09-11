package filters

import (
	"strconv"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
)

func ApplyOffset(query, value string) (queryResult string, offset int64, rErr error) {
	if value == "" {
		return query, 0, nil
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be a number.`).WithCause(err)
		return "", 0, rErr
	}

	if num < 1 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "offset". It must be greater than 0.`)
		return "", 0, rErr
	}

	var sb strings.Builder
	sb.WriteString(query)
	sb.WriteString(" OFFSET ")
	sb.WriteString(value)
	return sb.String(), num, nil
}
