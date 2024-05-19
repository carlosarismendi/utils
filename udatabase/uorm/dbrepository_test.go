package uorm

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/carlosarismendi/utils/udatabase"
	uormFilters "github.com/carlosarismendi/utils/udatabase/uorm/filters"
	"github.com/carlosarismendi/utils/uerr"
	"github.com/stretchr/testify/require"
)

type Resource struct {
	ID           string
	Name         string
	RandomNumber int
}

func createResourceTable(t testing.TB, r *DBrepository) {
	err := r.GetDBInstance(context.Background()).
		Exec("CREATE TABLE resources (id UUID PRIMARY KEY, name TEXT, random_number INTEGER);").Error
	require.NoError(t, err)
}

func TestTransactions(t *testing.T) {
	dbHolder := NewTestDBHolder("db_orm_repository_test_transactions")
	r := NewDBRepository(dbHolder.DBHolder, nil)

	t.Run("savingResourceWithoutError_commitsTransaction", func(t *testing.T) {
		// ARRANGE
		dbHolder.Reset()
		createResourceTable(t, r)

		resource := Resource{
			ID:           "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
			Name:         "Resource name",
			RandomNumber: 4,
		}

		err := func() (rErr error) {
			ctx, err := udatabase.BeginTx(context.Background(), r)
			if err != nil {
				return err
			}
			defer udatabase.EndTx(ctx, r, &rErr)

			// ACT
			return r.Save(ctx, &resource)
		}()
		require.NoError(t, err)

		// ASSERT
		var actual Resource
		err = r.FindByID(context.Background(), resource.ID, &actual)
		require.NoError(t, err)
	})

	t.Run("savingAResourceWithError_rollbacksTransaction", func(t *testing.T) {
		// ARRANGE
		dbHolder.Reset()
		createResourceTable(t, r)

		validResource := Resource{
			ID:           "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
			Name:         "Resource name",
			RandomNumber: 4,
		}

		err := func() (rErr error) {
			ctx, err := udatabase.BeginTx(context.Background(), r)
			if err != nil {
				return err
			}
			defer udatabase.EndTx(ctx, r, &rErr)

			// ACT
			err = r.Save(ctx, &validResource)
			if err != nil {
				return err
			}

			return fmt.Errorf("err")
		}()
		require.Error(t, err)

		// ASSERT
		var actual Resource
		err = r.FindByID(context.Background(), validResource.ID, &actual)
		require.Error(t, err)
		require.Equal(t, uerr.ResourceNotFoundError, uerr.GetKey(err), err)
	})

	t.Run("savingAResourceWithPanic_rollbacksTransactionAndReemitsPanic", func(t *testing.T) {
		// ARRANGE
		dbHolder.Reset()
		createResourceTable(t, r)

		resource := Resource{
			ID:           "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
			Name:         "Resource name",
			RandomNumber: 4,
		}

		defer func() {
			pErr := recover()
			require.NotNil(t, pErr, "Panic error expected not nil.")
			require.Error(t, pErr.(error))

			// ASSERT
			var actual Resource
			err := r.FindByID(context.Background(), resource.ID, &actual)
			require.Error(t, err)
			require.Equal(t, uerr.ResourceNotFoundError, uerr.GetKey(err), err)
		}()

		err := func() (rErr error) {
			ctx, err := udatabase.BeginTx(context.Background(), r)
			if err != nil {
				return err
			}
			defer udatabase.EndTx(ctx, r, &rErr)

			// ACT
			err = r.Save(ctx, &resource)
			if err != nil {
				return err
			}

			panic(fmt.Errorf("Fake panic error"))
		}()
		require.Error(t, err)
	})
}

func TestSave(t *testing.T) {
	dbHolder := NewTestDBHolder("db_orm_repository_test_save")
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

func TestFindWithFilters(t *testing.T) {
	dbHolder := NewTestDBHolder("db_orm_repository_test_find")
	dbHolder.Reset()

	filtersMap := map[string]uormFilters.Filter{
		"id":            uormFilters.TextField("id"),
		"name":          uormFilters.TextField("name"),
		"random_number": uormFilters.NumField("random_number"),
		"sort":          uormFilters.Sorter("name", "random_number"),
	}

	r := NewDBRepository(dbHolder.DBHolder, filtersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:    "FindingWithoutFilters",
			filters: url.Values{},
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r4, r3},
			},
			considerOrder: false,
		},
		{
			name:    "FindingFilteringByTextFieldName",
			filters: createFilter("name", "Resource1"),
			expected: &udatabase.ResourcePage{
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
			expected: &udatabase.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingFilteringByMultipleValuesInNumberField",
			filters: createFilter("name", "Resource1", "Resource2"),
			expected: &udatabase.ResourcePage{
				Total:     2,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2},
			},
			considerOrder: false,
		},
		{
			name:    "FindingFilteringByMultipleValuesInNumberField",
			filters: createFilter("random_number", "2", "0"),
			expected: &udatabase.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r4},
			},
			considerOrder: false,
		},
		{
			name:    "FindingFilteringBothByNumberAndTextField",
			filters: createFilters(newFilter("random_number", "2"), newFilter("name", "Resource3")),
			expected: &udatabase.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: false,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testRun(r))
	}
}

func TestFindWithValuedFilters(t *testing.T) {
	dbHolder := NewTestDBHolder("db_orm_repository_test_find")
	dbHolder.Reset()

	filtersMap := map[string]uormFilters.Filter{
		"id":            uormFilters.TextField("id"),
		"name":          uormFilters.TextField("name"),
		"random_number": uormFilters.NumField("random_number"),
		"sort":          uormFilters.Sorter("name", "random_number"),
	}

	r := NewDBRepository(dbHolder.DBHolder, filtersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:          "FindingWithoutFilters",
			valuedFilters: nil,
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r4, r3},
			},
			considerOrder: false,
		},
		{
			name: "FindingFilteringByTextFieldName",
			valuedFilters: []uormFilters.ValuedFilter{
				uormFilters.TextFieldWithValue("name", "Resource1"),
			},
			expected: &udatabase.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1},
			},
			considerOrder: true,
		},
		{
			name: "FindingFilteringByTextFieldID",
			valuedFilters: []uormFilters.ValuedFilter{
				uormFilters.TextFieldWithValue("id", "5ceff18d-9039-44b5-a5d3-3d99653f4603"),
			},
			expected: &udatabase.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: true,
		},
		{
			name: "FindingFilteringByMultipleValuesInNumberField",
			valuedFilters: []uormFilters.ValuedFilter{
				uormFilters.TextFieldWithValue("name", "Resource1", "Resource2"),
			},
			expected: &udatabase.ResourcePage{
				Total:     2,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2},
			},
			considerOrder: false,
		},
		{
			name: "FindingFilteringByMultipleValuesInNumberField",
			valuedFilters: []uormFilters.ValuedFilter{
				uormFilters.NumFieldWithValue("random_number", "2", "0"),
			},
			expected: &udatabase.ResourcePage{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r4},
			},
			considerOrder: false,
		},
		{
			name: "FindingFilteringBothByNumberAndTextField",
			valuedFilters: []uormFilters.ValuedFilter{
				uormFilters.NumFieldWithValue("random_number", "2"),
				uormFilters.TextFieldWithValue("name", "Resource3"),
			},
			expected: &udatabase.ResourcePage{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: false,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testRunWithValuedFilters(r))
	}
}

func TestFindWithSorters(t *testing.T) {
	dbHolder := NewTestDBHolder("db_orm_repository_test_find")
	dbHolder.Reset()

	filtersMap := map[string]uormFilters.Filter{
		"id":            uormFilters.TextField("id"),
		"name":          uormFilters.TextField("name"),
		"random_number": uormFilters.NumField("random_number"),
		"sort":          uormFilters.Sorter("name", "random_number"),
	}

	r := NewDBRepository(dbHolder.DBHolder, filtersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:    "FindingSortingByTextFieldNameAsc",
			filters: createFilter("sort", "name"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3, r4},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByTextFieldNameDesc",
			filters: createFilter("sort", "-name"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3, r4, r2, r1},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByFieldRandomNumberAsc",
			filters: createFilter("sort", "random_number"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r4, r1, r2, r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByNumFieldRandomNumberDesc",
			filters: createFilter("sort", "-random_number"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r1, r4},
			},
			considerOrder: true,
		},
		{
			name: "FindingSortingByNumFieldRandomNumberDescAnLimitTwo",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("limit", "2"),
			),
			expected: &udatabase.ResourcePage{
				Total:     2,
				Limit:     2,
				Offset:    0,
				Resources: []*Resource{r2, r3},
			},
			considerOrder: true,
		},
		{
			name: "FindingSortingByNumFieldRandomNumberDescAnOffset2",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("offset", "2"),
			),
			expected: &udatabase.ResourcePage{
				Total:     2,
				Limit:     10,
				Offset:    2,
				Resources: []*Resource{r1, r4},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByNameAscAndRandomNumberAsc",
			filters: createFilter("sort", "name", "random_number"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r4, r3},
			},
			considerOrder: true,
		},
		{
			name:    "FindingSortingByNameAscAndRandomNumberDesc",
			filters: createFilter("sort", "name", "-random_number"),
			expected: &udatabase.ResourcePage{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3, r4},
			},
			considerOrder: true,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testRun(r))
	}
}

func populateDB(ctx context.Context, t testing.TB, r *DBrepository) (r1, r2, r3, r4 *Resource) {
	r1 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4601",
		Name:         "Resource1",
		RandomNumber: 1,
	}
	require.NoError(t, r.Save(ctx, r1))

	r2 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4602",
		Name:         "Resource2",
		RandomNumber: 2,
	}
	require.NoError(t, r.Save(ctx, r2))

	r3 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4603",
		Name:         "Resource3",
		RandomNumber: 2,
	}
	require.NoError(t, r.Save(ctx, r3))

	r4 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4604",
		Name:         "Resource3",
		RandomNumber: 0,
	}
	require.NoError(t, r.Save(ctx, r4))
	return r1, r2, r3, r4
}

type filter struct {
	key    string
	values []string
}

func newFilter(key string, values ...string) *filter {
	return &filter{
		key:    key,
		values: values,
	}
}

func createFilters(fs ...*filter) url.Values {
	v := url.Values{}
	for _, f := range fs {
		for _, value := range f.values {
			v.Add(f.key, value)
		}
	}
	return v
}

func createFilter(key string, values ...string) url.Values {
	v := url.Values{}
	for _, value := range values {
		v.Add(key, value)
	}
	return v
}

type findTest struct {
	name          string
	filters       url.Values
	valuedFilters []uormFilters.ValuedFilter
	expected      *udatabase.ResourcePage
	considerOrder bool
}

func (ft *findTest) testRun(r *DBrepository) func(*testing.T) {
	return func(t *testing.T) {
		// ARRANGE
		expectedResources := ft.expected.Resources.([]*Resource)
		require.EqualValues(t, len(expectedResources), ft.expected.Total, "Expected Total and Resources doesn't match.")

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
	}
}

func (ft *findTest) testRunWithValuedFilters(r *DBrepository) func(*testing.T) {
	return func(t *testing.T) {
		// ARRANGE
		expectedResources := ft.expected.Resources.([]*Resource)
		require.EqualValues(t, len(expectedResources), ft.expected.Total, "Expected Total and Resources doesn't match.")

		// ACT
		var resources []*Resource
		rp, err := r.FindWithFilters(context.Background(), &resources, ft.valuedFilters...)

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
	}
}
