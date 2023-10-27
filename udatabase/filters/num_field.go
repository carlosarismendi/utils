package filters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
)

func ApplyNumField(fieldName string, values []string) (conds string, args []interface{}, rErr error) {
	amountValues := len(values)
	if amountValues == 0 {
		return "", args, rErr
	}

	args = make([]interface{}, 0, amountValues)
	var sb strings.Builder

	var num int64
	if amountValues == 1 {
		num, rErr = stoi(fieldName, values[0])
		if rErr != nil {
			return "", nil, rErr
		}
		args = append(args, num)

		sb.Grow(len(fieldName) + 4)
		sb.WriteByte('(')
		sb.WriteString(fieldName)
		sb.WriteString("=?")
		sb.WriteByte(')')
	} else {
		sb.Grow(len(fieldName) + 6 + 2*amountValues)
		sb.WriteString(fieldName)
		sb.WriteString(" IN (")
		sep := byte(' ')
		for _, v := range values {
			num, err := stoi(fieldName, v)
			if err != nil {
				return "", nil, err
			}
			args = append(args, num)
			sb.WriteByte(sep)
			sb.WriteByte('?')
			sep = ','
		}
		sb.WriteByte(')')
	}

	return sb.String(), args, nil
}

func stoi(fieldName, s string) (int64, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid value for filter %q. It must be a number.", fieldName)
		return 0, uerr.NewError(uerr.WrongInputParameterError, errMsg).WithCause(err)
	}

	return int64(num), nil
}
