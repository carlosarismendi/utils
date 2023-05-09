package testhelper

import (
	"testing"
)

func TestRequireCompare(t *testing.T) {
	type NestedResource struct {
		InnerID string
		Number  int
	}

	type Resource struct {
		ID     string
		Name   string
		Nested *NestedResource
	}

	t.Run("CompareEqualPointerStructs_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := &Resource{
			ID:   "id",
			Name: "name",
			Nested: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		r2 := &Resource{
			ID:   "id",
			Name: "name",
			Nested: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2)
	})

	t.Run("CompareEqualValueStructs_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
			Nested: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		r2 := Resource{
			ID:   "id",
			Name: "name",
			Nested: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2)
	})

	t.Run("CompareDifferentValueStructsIgnoringTheDifferentField_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
		}

		r2 := Resource{
			ID:   "id",
			Name: "name",
			Nested: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2, "Nested")
	})
}
