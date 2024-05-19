package filters

import "strings"

func ApplyTextField(fieldName string, values ...string) (conds string, args []interface{}, rErr error) {
	amountValues := len(values)
	if amountValues == 0 {
		return "", args, rErr
	}

	args = make([]interface{}, 0, amountValues)
	var sb strings.Builder
	if amountValues == 1 {
		args = append(args, values[0])

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
			args = append(args, v)
			sb.WriteByte(sep)
			sb.WriteByte('?')
			sep = ','
		}
		sb.WriteByte(')')
	}

	return sb.String(), args, nil
}
