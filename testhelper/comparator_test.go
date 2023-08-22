package testhelper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	type NestedResource struct {
		ID     string
		Number int
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

	type ResourceWithSlice struct {
		ID  string
		Arr []*Resource
	}

	type compareTest struct {
		name         string
		resource1    interface{}
		resource2    interface{}
		ignoreFields []string
		equal        bool
	}

	tests := []compareTest{
		{
			name: "CompareEqualPointerStructs_returnsEqual",
			resource1: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			resource2: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareEqualValueStructs_returnsEqual",
			resource1: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			resource2: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareEqualStructsOnePointerOneValue_returnsEqual",
			resource1: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			resource2: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareDifferentValueStructsIgnoringTheDifferentField_returnsEqual",
			resource1: Resource{
				ID:   "id",
				Name: "name",
			},
			resource2: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{"Nested"},
			equal:        true,
		},

		{
			name: "CompareDifferentPointerStructsIgnoringTheDifferentFieldThatIsInheritance_returnsEqual",
			resource1: &InheritanceResource{
				ID:   "id",
				Name: "name",
			},
			resource2: &InheritanceResource{
				ID:   "id",
				Name: "name",
				NestedResource: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{"NestedResource"},
			equal:        true,
		},

		{
			name: "CompareDifferentValueStructsIgnoringTheDifferentFieldThatIsInheritance_returnsNoError",
			resource1: InheritanceResource{
				ID:   "id",
				Name: "name",
			},
			resource2: InheritanceResource{
				ID:   "id",
				Name: "name",
				NestedResource: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{"NestedResource"},
			equal:        true,
		},

		{
			name: "CompareDifferentStructsIgnoringDifferentFieldThatIsInheritanceOneValueOnePointer_returnsEqual",
			resource1: &InheritanceResource{
				ID:   "id",
				Name: "name",
			},
			resource2: InheritanceResource{
				ID:   "id",
				Name: "name",
				NestedResource: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{"NestedResource"},
			equal:        true,
		},

		{
			name: "CompareEqualMaps_returnsEqual",
			resource1: map[string]int{
				"field1": 1,
				"field2": 3,
			},
			resource2: map[string]int{
				"field1": 1,
				"field2": 3,
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareDifferentMapsIgnoringFieldOfInnerStruct_returnsEqual",
			resource1: map[string]*Resource{
				"field1": {
					ID:   "id1",
					Name: "name1",
				},
			},
			resource2: map[string]*Resource{
				"field1": {
					ID:   "id1",
					Name: "different",
				},
			},
			ignoreFields: []string{"Name"},
			equal:        true,
		},

		{
			name:         "CompareEqualSlicesOfInts_returnsEqual",
			resource1:    []int{1, 2, 3},
			resource2:    []int{1, 2, 3},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareEqualSlicesOfValueStructs_returnsEqual",
			resource1: []Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			resource2: []Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareEqualSlicesOfPointerStructs_returnsEqual",
			resource1: []*Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			resource2: []*Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareDifferentSlicesOfValueStructsIgnoringDifferentField_returnsEqual",
			resource1: []Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			resource2: []Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "different",
				},
			},
			ignoreFields: []string{"Name"},
			equal:        true,
		},

		{
			name: "CompareDifferentSlicesOfPointerStructsIgnoringDifferentField_returnsEqual",
			resource1: []*Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "name2",
				},
			},
			resource2: []*Resource{
				{
					ID:   "id",
					Name: "name",
				},
				{
					ID:   "id2",
					Name: "different",
				},
			},
			ignoreFields: []string{"Name"},
			equal:        true,
		},

		{
			name: "CompareDifferentResourcesWithDifferentFieldInInnerResourceOfSlice_returnEqual",
			resource1: ResourceWithSlice{
				ID: "id",
				Arr: []*Resource{
					{
						ID:   "id",
						Name: "name1",
					},
				},
			},
			resource2: ResourceWithSlice{
				ID: "id",
				Arr: []*Resource{
					{
						ID:   "id",
						Name: "different",
					},
				},
			},
			ignoreFields: []string{"Arr.Name"},
			equal:        true,
		},

		{
			name:         "CompareEqualTimestamps_returnEqual",
			resource1:    time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC),
			resource2:    time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC),
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name:         "CompareStructWithEqualTimestamps_returnEqual",
			resource1:    struct{ Timestamp time.Time }{Timestamp: time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC)},
			resource2:    struct{ Timestamp time.Time }{Timestamp: time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC)},
			ignoreFields: []string{},
			equal:        true,
		},

		{
			name: "CompareStructIgnoringInnerFieldOfSubStruct_returnEqual",
			resource1: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID_diff",
					Number: 25,
				},
			},
			resource2: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					ID:     "innerID",
					Number: 25,
				},
			},
			ignoreFields: []string{"Nested.ID"},
			equal:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			equal := compare(test.resource1, test.resource2, "", test.ignoreFields...)
			require.Equal(t, test.equal, equal, getErrorMessage(test.resource1, test.resource2))
		})
	}
}
