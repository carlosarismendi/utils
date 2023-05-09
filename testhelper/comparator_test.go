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

	type InheritanceResource struct {
		ID   string
		Name string
		*NestedResource
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

	t.Run("CompareEqualStructsOnePointerOneValue_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := &Resource{
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

	t.Run("CompareDifferentPointerStructsIgnoringTheDifferentFieldThatIsInheritance_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := &InheritanceResource{
			ID:   "id",
			Name: "name",
		}

		r2 := &InheritanceResource{
			ID:   "id",
			Name: "name",
			NestedResource: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2, "NestedResource")
	})

	t.Run("CompareDifferentValueStructsIgnoringTheDifferentFieldThatIsInheritance_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := InheritanceResource{
			ID:   "id",
			Name: "name",
		}

		r2 := InheritanceResource{
			ID:   "id",
			Name: "name",
			NestedResource: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2, "NestedResource")
	})

	t.Run("CompareDifferentStructsIgnoringDifferentFieldThatIsInheritanceOneValueOnePointer"+
		"_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := &InheritanceResource{
			ID:   "id",
			Name: "name",
		}

		r2 := InheritanceResource{
			ID:   "id",
			Name: "name",
			NestedResource: &NestedResource{
				InnerID: "innerID",
				Number:  25,
			},
		}

		// ACT
		RequireCompare(t, r1, r2, "NestedResource")
	})

	t.Run("CompareEqualMaps_returnsNoError", func(t *testing.T) {
		// ARRANGE
		m1 := map[string]int{
			"field1": 1,
			"field2": 3,
		}

		m2 := map[string]int{
			"field1": 1,
			"field2": 3,
		}

		// ACT
		RequireCompare(t, m1, m2)
	})

	t.Run("CompareDifferentMapsIgnoringDifferentField_returnsNoError", func(t *testing.T) {
		// ARRANGE
		m1 := map[string]int{
			"field1": 1,
			"field2": 3,
		}

		m2 := map[string]int{
			"field1": 1,
			"field2": 3000,
		}

		// ACT
		RequireCompare(t, m1, m2, "field2")
	})

	t.Run("CompareEqualSlicesOfInts_returnsNoError", func(t *testing.T) {
		// ARRANGE
		s1 := []int{1, 2, 3}

		s2 := []int{1, 2, 3}

		// ACT
		RequireCompare(t, s1, s2)
	})

	t.Run("CompareEqualSlicesOfValueStructs_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
		}

		r2 := Resource{
			ID:   "id2",
			Name: "name2",
		}

		s1 := []Resource{r1, r2}

		s2 := []Resource{r1, r2}

		// ACT
		RequireCompare(t, s1, s2)
	})

	t.Run("CompareEqualSlicesOfPointerStructs_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
		}

		r2 := Resource{
			ID:   "id2",
			Name: "name2",
		}

		r11 := Resource{
			ID:   "id",
			Name: "name",
		}

		r22 := Resource{
			ID:   "id2",
			Name: "name2",
		}

		s1 := []*Resource{&r1, &r2}

		s2 := []*Resource{&r11, &r22}

		// ACT
		RequireCompare(t, s1, s2)
	})

	t.Run("CompareDifferentSlicesOfValueStructsIgnoringDifferentField_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
		}

		r2 := Resource{
			ID:   "id2",
			Name: "name2",
		}

		r11 := Resource{
			ID:   "id",
			Name: "name",
		}

		r22 := Resource{
			ID:   "id2",
			Name: "name",
		}

		s1 := []Resource{r1, r2}

		s2 := []Resource{r11, r22}

		// ACT
		RequireCompare(t, s1, s2, "Name")
	})

	t.Run("CompareDifferentSlicesOfPointerStructsIgnoringDifferentField_returnsNoError", func(t *testing.T) {
		// ARRANGE
		r1 := Resource{
			ID:   "id",
			Name: "name",
		}

		r2 := Resource{
			ID:   "id2",
			Name: "name2",
		}

		r11 := Resource{
			ID:   "id",
			Name: "name",
		}

		r22 := Resource{
			ID:   "id2",
			Name: "name",
		}

		s1 := []*Resource{&r1, &r2}

		s2 := []*Resource{&r11, &r22}

		// ACT
		RequireCompare(t, s1, s2, "Name")
	})
}
