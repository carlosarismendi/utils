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

// If the comparison fails, it kills the execution.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/require.
func RequireCompare(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		require.Fail(t, fmt.Sprintf("Expected (%T) and actual (%T) are not the same type.", expected, actual))
	}

	equal := compare(expected, actual, ignoreFields...)
	require.True(t, equal, getErrorMessage(expected, actual))
}

// If the comparison fails, it does not kill the execution.
// Uses methods from https://pkg.go.dev/github.com/stretchr/testify/assert.
func AssertCompare(t testing.TB, expected, actual interface{}, ignoreFields ...string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		assert.Fail(t, fmt.Sprintf("Expected (%T) and actual (%T) are not the same type.", expected, actual))
	}

	equal := compare(expected, actual, ignoreFields...)
	assert.True(t, equal, getErrorMessage(expected, actual))
}

func compare(expected, actual interface{}, ignoreFields ...string) bool {
	ignoreType := expected
	t := reflect.TypeOf(expected)
	if t.Kind() == reflect.Ptr {
		ignoreType = reflect.ValueOf(t.Elem())
	}

	equal := cmp.Equal(expected, actual, cmpopts.IgnoreFields(ignoreType, ignoreFields...))
	return equal
}

func getErrorMessage(expected, actual interface{}) string {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)

	msg := fmt.Sprintf("Expected: %s", string(expectedJSON))
	msg += fmt.Sprintf("\nActual: %s", string(actualJSON))
	return msg
}
