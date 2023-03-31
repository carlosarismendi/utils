package infrastructure

import (
	"context"
	"net/url"
	"testing"

	"github.com/carlosarismendi/utils/db/domain"
	"github.com/carlosarismendi/utils/db/infrastructure/filters"
	"github.com/stretchr/testify/require"
)

type Resource struct {
	ID           string
	Name         string
	RandomNumber int
}

func createResourceTable(t testing.TB, r *DBrepository) {
	ctx, err := domain.BeginTx(context.Background(), r)
	require.NoError(t, err)

	err = r.db.Exec("CREATE TABLE resources (id UUID, name TEXT, random_number INTEGER);").Error
	domain.EndTx(ctx, r, &err)
}

func TestSave(t *testing.T) {
	dbHolder := NewTestDBHolder("db_repository_test_save")
	r := NewDBRepository(dbHolder.DBHolder, nil)

	t.Run("SavingValidResource", func(t *testing.T) {
		// ARRANGE
		dbHolder.Reset()
		createResourceTable(t, r)

		resource := Resource{
			ID:           "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
			Name:         "Resource name",
			RandomNumber: 4,
		}

		// ACT
		ctx := context.Background()
		err := r.Save(ctx, &resource)
		require.NoError(t, err)

		// ASSERT
		var actual Resource
		err = r.FindByID(ctx, resource.ID, &actual)
		require.NoError(t, err)
	})

	t.Run("SavingInvalidResource", func(t *testing.T) {
		// ARRANGE
		dbHolder.Reset()
		createResourceTable(t, r)

		resource := Resource{
			ID:           "INVALID_UUID",
			Name:         "Resource",
			RandomNumber: 2,
		}

		// ACT
		err := r.Save(context.Background(), &resource)

		// ASSERT
		require.Error(t, err)
	})
}

func TestFind(t *testing.T) {
	dbHolder := NewTestDBHolder("db_repository_test_find")
	dbHolder.Reset()

	filtersMap := map[string]filters.Filter{
		"id":            filters.TextField("id"),
		"name":          filters.TextField("name"),
		"random_number": filters.NumField("random_number"),
		"sort":          filters.Sorter(),
	}

	r := NewDBRepository(dbHolder.DBHolder, filtersMap)
	createResourceTable(t, r)

	ctx := context.Background()

	r1 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4601",
		Name:         "Resource1",
		RandomNumber: 1,
	}
	require.NoError(t, r.Save(ctx, r1))

	r2 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4602",
		Name:         "Resource2",
		RandomNumber: 2,
	}
	require.NoError(t, r.Save(ctx, r2))

	r3 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4603",
		Name:         "Resource3",
		RandomNumber: 2,
	}
	require.NoError(t, r.Save(ctx, r3))

	type findTest struct {
		name          string
		filters       url.Values
		expected      *domain.ResourcePage
		considerOrder bool
	}

	tests := []findTest{
		{
			name:    "FindingWithoutFilters",
			filters: url.Values{},
			expected: &domain.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3},
			},
			considerOrder: false,
		},
		{
			name:    "FindingFilteringByTextFieldName",
			filters: createFilter("name", "Resource1"),
			expected: &domain.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1},
			},
			considerOrder: true,
		},
		{
			name:    "FindingFilteringByTextFieldID",
			filters: createFilter("id", "5ceff18d-9039-44b5-a5d3-3d99653f4603"),
			expected: &domain.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingFilteringByNumberFieldRandomNumber",
			filters: createFilter("random_number", "2"),
			expected: &domain.ResourcePage{
				Total:     2,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3},
			},
			considerOrder: false,
		},
		{
			name:    "FindingSortingByTextFieldNameAsc",
			filters: createFilter("sort", "name"),
			expected: &domain.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByTextFieldNameDesc",
			filters: createFilter("sort", "-name"),
			expected: &domain.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3, r2, r1},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByFieldRandomNumberAsc",
			filters: createFilter("sort", "random_number"),
			expected: &domain.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByNumFieldRandomNumberDesc",
			filters: createFilter("sort", "-random_number"),
			expected: &domain.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r1},
			},
			considerOrder: true,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, func(t *testing.T) {
			// ARRANGE
			expectedResources := ft.expected.Resources.([]*Resource)
			require.EqualValues(t, len(expectedResources), ft.expected.Total, "Expected Total and Resources are not properly set.")

			// ACT
			var resources []*Resource
			rp, err := r.Find(context.Background(), ft.filters, &resources)

			// ASSERT
			require.NoError(t, err)
			require.EqualValues(t, ft.expected.Total, len(resources))
			require.EqualValues(t, ft.expected.Total, rp.Total)
			require.EqualValues(t, ft.expected.Total, len(*rp.Resources.(*[]*Resource)))
			require.EqualValues(t, ft.expected.Offset, rp.Offset)
			require.EqualValues(t, ft.expected.Limit, rp.Limit)

			if ft.considerOrder {
				require.Equal(t, expectedResources, resources)
			} else {
				for _, expRes := range expectedResources {
					require.Contains(t, resources, expRes)
				}
			}
		})
	}
}

func createFilter(key string, values ...string) url.Values {
	v := url.Values{}
	for i := range values {
		v.Add(key, values[i])
	}
	return v
}
