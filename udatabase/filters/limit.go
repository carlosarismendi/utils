package filters

import (
	"strconv"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
)

const DefaultLimit = 10
const DefaultLimitStr = " LIMIT 10"

func ApplyLimit(query, value string) (queryResult string, limit int64, rErr error) {
	if value == "" {
		return query, 0, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be a number.`).WithCause(err)
		return "", 0, rErr
	}

	if num < 1 {
		rErr := uerr.NewError(uerr.WrongInputParameterError,
			`Invalid value for "limit". It must be greater than 0.`)
		return "", 0, rErr
	}

	var sb strings.Builder
	sb.WriteString(query)
	sb.WriteString(" LIMIT ")
	sb.WriteString(value)
	return sb.String(), int64(num), nil
}
