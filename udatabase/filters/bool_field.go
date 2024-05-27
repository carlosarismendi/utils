package filters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/carlosarismendi/utils/uerr"
)

// ApplyBoolField only applies first value from values.
func ApplyBoolField(fieldName string, values ...string) (conds string, args []interface{}, rErr error) {
	amountValues := len(values)
	if amountValues == 0 {
		return "", args, rErr
	}

	args = make([]interface{}, 0, 1)
	var sb strings.Builder

	var value bool
	value, rErr = stob(fieldName, values[0])
	if rErr != nil {
		return "", nil, rErr
	}
	args = append(args, value)

	sb.Grow(len(fieldName) + 4)
	sb.WriteByte('(')
	sb.WriteString(fieldName)
	sb.WriteString("=?")
	sb.WriteByte(')')

	return sb.String(), args, nil
}

func stob(fieldName, s string) (bool, error) {
	value, err := strconv.ParseBool(s)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid value for filter %q. It must be 'true' or 'false'.", fieldName)
		return false, uerr.NewError(uerr.WrongInputParameterError, errMsg).WithCause(err)
	}

	return value, nil
}
