package testhelper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
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

	type ResourceWithSlice struct {
		ID  string
		Arr []*Resource
	}

	type compareTest struct {
		name         string
		resource1    interface{}
		resource2    interface{}
		ignoreFields []string
	}

	tests := []compareTest{
		{
			name: "CompareEqualPointerStructs_returnsEqual",
			resource1: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			resource2: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{},
		},

		{
			name: "CompareEqualValueStructs_returnsEqual",
			resource1: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			resource2: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{},
		},

		{
			name: "CompareEqualStructsOnePointerOneValue_returnsEqual",
			resource1: &Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			resource2: Resource{
				ID:   "id",
				Name: "name",
				Nested: &NestedResource{
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{},
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
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{"Nested"},
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
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{"NestedResource"},
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
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{"NestedResource"},
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
					InnerID: "innerID",
					Number:  25,
				},
			},
			ignoreFields: []string{"NestedResource"},
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
		},

		{
			name:         "CompareEqualSlicesOfInts_returnsEqual",
			resource1:    []int{1, 2, 3},
			resource2:    []int{1, 2, 3},
			ignoreFields: []string{},
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
			ignoreFields: []string{"Name"},
		},

		{
			name:         "CompareEqualTimestamps_returnEqual",
			resource1:    time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC),
			resource2:    time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC),
			ignoreFields: []string{},
		},

		{
			name:         "CompareStructWithEqualTimestamps_returnEqual",
			resource1:    struct{ Timestamp time.Time }{Timestamp: time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC)},
			resource2:    struct{ Timestamp time.Time }{Timestamp: time.Date(2022, 10, 5, 10, 49, 8, 1, time.UTC)},
			ignoreFields: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			equal := compare(test.resource1, test.resource2, test.ignoreFields...)
			require.True(t, equal, getErrorMessage(test.resource1, test.resource2))
		})
	}
}
