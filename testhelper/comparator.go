package testhelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If the comparison fails, it kills the execution. When comparing maps, ignoreFields
// only applies for structs that are the values of the map.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/require.
func RequireEqual(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	equal := compare(expected, actual, ignoreFields...)
	if !equal {
		require.Fail(t, getErrorMessage(expected, actual))
	}
}

// If the comparison fails, it does not kill the execution.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/assert.
func AssertEqual(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	equal := compare(expected, actual, ignoreFields...)
	if !equal {
		assert.Fail(t, getErrorMessage(expected, actual))
	}
}

func compare(expected, actual interface{}, ignoreFields ...string) bool {
	exp, act, sameTypes := getValues(expected, actual)
	if !sameTypes {
		return false
	}

	if !exp.IsValid() && !act.IsValid() {
		return true
	}

	if exp.IsValid() && !act.IsValid() {
		return false
	}

	if !exp.IsValid() && act.IsValid() {
		return false
	}

	equal := compareValue(exp, act, ignoreFields...)

	return equal
}

func compareValue(exp, act reflect.Value, ignoreFields ...string) bool {
	switch exp.Kind() {
	case reflect.Struct:
		return compareStructs(exp, act, ignoreFields...)

	case reflect.Map:
		return compareMaps(exp, act, ignoreFields...)

	case reflect.Slice, reflect.Array:
		return compareSlices(exp, act, ignoreFields...)

	default:
		// Primitive types (int, uint, float, string, complex, bool) and
		// others that don't fit in previous cases.
		return exp.Interface() == act.Interface()
	}
}

func compareStructs(exp, act reflect.Value, ignoreFields ...string) bool {
	typ := exp.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if isIgnoreField(field.Name, ignoreFields) {
			continue
		}

		equal := compare(exp.Field(i).Interface(), act.Field(i).Interface(), ignoreFields...)
		if !equal {
			return false
		}
	}

	return true
}

func compareMaps(exp, act reflect.Value, ignoreFields ...string) bool {
	if exp.Len() != act.Len() {
		return false
	}

	for _, key := range exp.MapKeys() {
		actValue := act.MapIndex(key)
		// Key from expected Map doest not exist in actual Map
		if !actValue.IsValid() {
			return false
		}

		equal := compare(exp.MapIndex(key).Interface(), actValue.Interface(), ignoreFields...)
		if !equal {
			return false
		}
	}

	return true
}

func compareSlices(exp, act reflect.Value, ignoreFields ...string) bool {
	if exp.Len() != act.Len() {
		return false
	}

	for i := 0; i < exp.Len(); i++ {
		equal := compare(exp.Index(i).Interface(), act.Index(i).Interface(), ignoreFields...)
		if !equal {
			return false
		}
	}

	return true
}

func getErrorMessage(expected, actual interface{}) string {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	msg := fmt.Sprintf("Expected: %s", string(expectedJSON))
	msg += fmt.Sprintf("\nActual: %s", string(actualJSON))
	return msg
}

func getValues(expected, actual interface{}) (exp, act reflect.Value, sameType bool) {
	exp = reflect.ValueOf(expected)
	if exp.Kind() == reflect.Ptr || exp.Kind() == reflect.Interface && !exp.IsNil() {
		exp = exp.Elem()
	}

	act = reflect.ValueOf(actual)
	if act.Kind() == reflect.Ptr || act.Kind() == reflect.Interface && !act.IsNil() {
		act = act.Elem()
	}

	if !exp.IsValid() && !act.IsValid() {
		return exp, act, true
	}

	if exp.IsValid() && !act.IsValid() {
		return exp, act, false
	}

	if !exp.IsValid() && act.IsValid() {
		return exp, act, false
	}

	if exp.Kind() != act.Kind() || exp.Type() != act.Type() {
		if act.CanConvert(exp.Type()) {
			act = act.Convert(exp.Type())
		} else {
			return exp, act, false
		}
	}

	return exp, act, true
}

func isIgnoreField(field string, ignoreFields []string) bool {
	for i := range ignoreFields {
		if field == ignoreFields[i] {
			return true
		}
	}

	return false
}
