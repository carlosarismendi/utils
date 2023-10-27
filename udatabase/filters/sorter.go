package filters

import (
	"fmt"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
)

func ApplySorter(allowedFields map[string]bool, values []string) (queryResult string, rErr error) {
	sep := byte(' ')
	var sb strings.Builder
	var column, direction string
	for _, v := range values {
		column, direction, rErr = SortFieldAndDirection(allowedFields, v)
		if rErr != nil {
			return "", rErr
		}
		sb.Grow(2 + len(column) + len(direction))
		sb.WriteByte(sep)
		sb.WriteString(column)
		sb.WriteByte(' ')

		if len(direction) > 0 {
			sb.WriteString(direction)
		}
		sep = ','
	}

	return sb.String(), nil
}

func SortFieldAndDirection(allowedFields map[string]bool, value string) (col, dir string, err error) {
	if value[0] == '-' {
		col = value[1:]
		dir = "DESC"
	} else {
		col = value
	}

	if !allowedFields[col] {
		return "", "", uerr.NewError(uerr.WrongInputParameterError, fmt.Sprintf("Invalid sort field %q.", col))
	}

	return col, dir, nil
}
