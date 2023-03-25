package validate

import (
	"net/http"
	"testing"

	"github.com/ansel1/merry"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	type user struct {
		ID   string `json:"id" validate:"uuid"`
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=0,max=100"`
	}

	t.Run("ValidatingStructWithValidData_returnsNoError", func(t *testing.T) {
		u := &user{
			ID:   "7f64a006-773b-4fa8-a10f-2fc899387256",
			Name: "Juan Francisco",
			Age:  18,
		}

		err := Validate(u)
		require.NoError(t, err)
	})

	t.Run("ValidatingStructWithOneInvalidField_returnsStatusUnprocessableEntity", func(t *testing.T) {
		// ARRANGE
		u := &user{
			ID:   "INVALID_ID",
			Name: "Juan Francisco",
			Age:  18,
		}

		// ACT
		err := Validate(u)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, "Invalid field ID: the value must be 'uuid'. The value received is 'INVALID_ID'", err.Error())
		require.Equal(t, http.StatusUnprocessableEntity, merry.HTTPCode(err))
	})

	t.Run("ValidatingStructWithTwoInvalidFields_returnsStatusUnprocessableEntity", func(t *testing.T) {
		// ARRANGE
		u := &user{
			ID:   "INVALID_ID",
			Name: "Juan Francisco",
			Age:  -1,
		}

		// ACT
		err := Validate(u)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, "Invalid field ID: the value must be 'uuid'. The value received is 'INVALID_ID': "+
			"\nInvalid field Age: the value must be 'min=0'. The value received is '-1'", err.Error())
		require.Equal(t, http.StatusUnprocessableEntity, merry.HTTPCode(err))
	})
}
