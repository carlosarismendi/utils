package uerr

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUError_MarshalJSON(t *testing.T) {
	t.Run("createFromBytesErrorWithUErrAsCause", func(t *testing.T) {
		// ARRANGE
		originalErr := NewError("testKey", "testMessage").
			WithCause(NewError("causeKey", "causeMessage"))

		// ACT
		actual, err := json.Marshal(originalErr)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t,
			`{"error":{"key":"testKey","message":"testMessage","cause":{"error":{"key":"causeKey","message":"causeMessage"}}}}`,
			string(actual),
		)
	})

	t.Run("createFromBytesErrorWithFmtErrAsCause", func(t *testing.T) {
		// ARRANGE
		originalErr := NewError("testKey", "testMessage").
			WithCause(fmt.Errorf("causeErr"))

		// ACT
		actual, err := json.Marshal(originalErr)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t,
			`{"error":{"key":"testKey","message":"testMessage","cause":"causeErr"}}`,
			string(actual),
		)
	})
}

func TestUError_FromBytes(t *testing.T) {
	t.Run("createFromBytesErrorWithCause", func(t *testing.T) {
		// ARRANGE
		originalErr := NewError("testKey", "testMessage").
			WithCause(NewError("causeKey", "causeMessage"))
		errBytes, err := originalErr.MarshalJSON()
		require.NoError(t, err)

		// ACT
		actualErr, err := FromBytes(errBytes)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, originalErr.key, actualErr.key)
		require.Equal(t, originalErr.message, actualErr.message)
	})
}

func TestUError_UnmarshalJSON(t *testing.T) {
	t.Run("createFromBytesErrorWithCause", func(t *testing.T) {
		// ARRANGE
		originalErr := NewError("testKey", "testMessage").
			WithCause(NewError("causeKey", "causeMessage"))
		errBytes, err := originalErr.MarshalJSON()
		require.NoError(t, err)

		// ACT
		var actualErr UError
		err = json.Unmarshal(errBytes, &actualErr)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, originalErr.key, actualErr.key)
		require.Equal(t, originalErr.message, actualErr.message)
	})
}
