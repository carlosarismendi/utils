package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	t.Run("CreatingResourceWithValidData_returnsNoError", func(t *testing.T) {
		// ARRANGE
		id := uuid.New().String()
		timestamp := time.Now().UTC()

		expectedResource := &Resource{
			ID:        id,
			CreatedAt: timestamp,
		}

		// ACT
		r, err := NewResource(id, timestamp)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, expectedResource, r)
	})

	t.Run("CreatingResourceWithInvalidID_returnsInvalidIDError", func(t *testing.T) {
		// ARRANGE
		id :=  "INVALID_ID"
		timestamp := time.Now().UTC()

		// ACT
		r, err := NewResource(id, timestamp)

		// ASSERT
		require.Nil(t, r)
		require.Error(t, err)
		require.Equal(t, "Invalid field ID: it must be a valid uuid. The value received is 'INVALID_ID'.", err.Error())
	})
}
