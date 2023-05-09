package testhelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If the comparison fails, it kills the execution. When comparing maps,
// if no ignore field is indicated, this works fine. In case of comparing maps with
// ignore fields, these map keys must be of type string. Otherwise this function won't work.
// In this case, it is recommended to use RequireAdvanceCompare method.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/require.
func RequireCompare(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	equal := compare(expected, actual, ignoreFields...)
	require.True(t, equal, getErrorMessage(expected, actual))
}

// If the comparison fails, it does not kill the execution.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/assert.
func AssertCompare(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	equal := compare(expected, actual, ignoreFields...)
	assert.True(t, equal, getErrorMessage(expected, actual))
}

// If the comparison fails, it kills the execution.
// This is a wrapper function for cmp.Equal provided by the package github.com/google/go-cmp.
// This just offers a message showing the JSON string for both expected and actual parameter.
// See https://pkg.go.dev/github.com/google/go-cmp/cmp
func RequireAdvanceCompare(t testing.TB, expected, actual interface{}, options ...cmp.Option) {
	equal := cmp.Equal(expected, actual, options...)
	require.True(t, equal, getErrorMessage(expected, actual))
}

// If the comparison fails, it does not kill the execution.
// This is a wrapper function for cmp.Equal provided by the package github.com/google/go-cmp.
// This just offers a message showing the JSON string for both expected and actual parameter.
// See https://pkg.go.dev/github.com/google/go-cmp/cmp
func AssertAdvanceCompare(t testing.TB, expected, actual interface{}, options ...cmp.Option) {
	equal := cmp.Equal(expected, actual, options...)
	assert.True(t, equal, getErrorMessage(expected, actual))
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

	var equal bool
	if len(ignoreFields) > 0 {
		if exp.Kind() == reflect.Map {
			equal = compareMap(exp, act, ignoreFields...)
		} else if exp.Kind() == reflect.Slice || exp.Kind() == reflect.Array {
			equal = compareSlice(exp, act, ignoreFields...)
		} else {
			equal = cmp.Equal(exp.Interface(), act.Interface(), cmpopts.IgnoreFields(exp.Interface(), ignoreFields...))
		}
	} else {
		equal = cmp.Equal(exp.Interface(), act.Interface())
	}

	return equal
}

func compareMap(exp, act reflect.Value, ignoreFields ...string) bool {
	ignoreFieldsOption := cmpopts.IgnoreMapEntries(func(k, v interface{}) bool {
		for i := range ignoreFields {
			if k == ignoreFields[i] {
				return true
			}
		}

		return false
	})
	return cmp.Equal(exp.Interface(), act.Interface(), ignoreFieldsOption)
}

func compareSlice(exp, act reflect.Value, ignoreFields ...string) bool {
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

	if exp.Kind() != act.Kind() || exp.Type() != act.Type() {
		if act.CanConvert(exp.Type()) {
			act = act.Convert(exp.Type())
		} else {
			return exp, act, false
		}
	}

	return exp, act, true
}
