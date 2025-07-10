package usql

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/carlosarismendi/testhelper"
	"github.com/carlosarismendi/utils/udatabase"
	"github.com/carlosarismendi/utils/udatabase/filters"
	usqlFilters "github.com/carlosarismendi/utils/udatabase/usql/filters"
	"github.com/carlosarismendi/utils/uerr"
	"github.com/stretchr/testify/require"
)

type Resource struct {
	ID           string
	Name         string
	RandomNumber int
}

func createResourceTable(t testing.TB, r *DBrepository[*Resource]) {
	_, err := r.GetDBInstance().
		Exec("CREATE TABLE resources (id UUID PRIMARY KEY, name TEXT, random_number INTEGER);")
	require.NoError(t, err)
}

func TestTransactions(t *testing.T) {
	dbHolder := NewTestDBHolder("db_sql_repository_test_transactions")
	r := NewDBRepository[*Resource](dbHolder.DBHolder, nil, nil)

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
			return save(ctx, r, &resource)
		}()
		require.NoError(t, err)

		// ASSERT
		var actual Resource
		err = findByID(r, resource.ID, &actual)
		require.NoError(t, err)
		testhelper.RequireEqual(t, resource, actual)
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
			err = save(ctx, r, &validResource)
			if err != nil {
				return err
			}

			return fmt.Errorf("err")
		}()
		require.Error(t, err)

		// ASSERT
		var actual Resource
		err = findByID(r, validResource.ID, &actual)
		require.Error(t, err)
		require.Equal(t, uerr.ResourceNotFoundError, uerr.GetKey(err), err)
		testhelper.RequireEqual(t, Resource{}, actual)
	})

	t.Run("savingAResourceWithPanic_rollbacksTransactionAndRelaunchesPanic", func(t *testing.T) {
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
			err := findByID(r, resource.ID, &actual)
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
			err = save(ctx, r, &resource)
			if err != nil {
				return err
			}

			panic(fmt.Errorf("Fake panic error"))
		}()
		require.NoError(t, err)
	})
}

func TestInsertErrors(t *testing.T) {
	dbHolder := NewTestDBHolder("db_sql_repository_test_insert_errors")
	r := NewDBRepository[*Resource](dbHolder.DBHolder, nil, nil)

	t.Run("Inserting two elements with same value for primary key", func(t *testing.T) {
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

			return save(ctx, r, &resource)
		}()
		require.NoError(t, err)

		ctx, err := udatabase.BeginTx(context.Background(), r)
		require.NoError(t, err)
		defer r.Rollback(ctx)

		// ACT
		err = save(ctx, r, &resource)

		// ASSERT
		require.Error(t, err)
		require.Equal(t, uerr.ResourceAlreadyExistsError, uerr.GetKey(err))
	})
}

func TestFindList(t *testing.T) {
	dbHolder := NewTestDBHolder("db_sql_repository_test_find")
	dbHolder.Reset()

	r := NewDBRepository[*Resource](dbHolder.DBHolder, nil, nil)
	createResourceTable(t, r)

	ctx := context.Background()
	ctx, err := r.Begin(ctx)
	require.NoError(t, err)
	defer r.Rollback(ctx)

	r1 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4601",
		Name:         "Resource1",
		RandomNumber: 1,
	}
	require.NoError(t, save(ctx, r, r1))

	r2 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4602",
		Name:         "Resource2",
		RandomNumber: 2,
	}
	require.NoError(t, save(ctx, r, r2))

	r3 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4603",
		Name:         "Resource3",
		RandomNumber: 2,
	}
	require.NoError(t, save(ctx, r, r3))

	r4 := &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4604",
		Name:         "Resource3",
		RandomNumber: 1,
	}
	require.NoError(t, save(ctx, r, r4))
	require.NoError(t, r.Commit(ctx))

	t.Run("Select list of resources with LIMIT", func(t *testing.T) {
		// ACT
		var dst []*Resource
		err = findList(r, &dst, "1", "")

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, 1, len(dst))
	})

	t.Run("Select list of resources with OFFSET", func(t *testing.T) {
		// ACT
		var dst []*Resource
		err = findList(r, &dst, "", "1")

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, 3, len(dst))
	})
}

func TestSelectContextWithFilters(t *testing.T) {
	dbHolder := NewTestDBHolder("db_usql_repository_test_select_context")
	dbHolder.Reset()

	filtersMap := map[string]usqlFilters.Filter{
		"id":            usqlFilters.TextField("id"),
		"name":          usqlFilters.TextField("name"),
		"random_number": usqlFilters.NumField("random_number"),
	}

	sortersMap := map[string]usqlFilters.Sorter{
		"sort": usqlFilters.Sort("name", "random_number"),
	}

	r := NewDBRepository[*Resource](dbHolder.DBHolder, filtersMap, sortersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:    "SelectContextWithoutFilters",
			filters: url.Values{},
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r4, r3},
			},
			considerOrder: false,
		},
		{
			name:    "SelectContextFilteringByTextFieldName",
			filters: createFilter("name", "Resource1"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextFilteringByTextFieldID",
			filters: createFilter("id", "5ceff18d-9039-44b5-a5d3-3d99653f4603"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextFilteringByMultipleValuesInNumberField",
			filters: createFilter("name", "Resource1", "Resource2"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     2,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2},
			},
			considerOrder: false,
		},
		{
			name:    "SelectContextFilteringByMultipleValuesInNumberField",
			filters: createFilter("random_number", "2", "0"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     3,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r4},
			},
			considerOrder: false,
		},
		{
			name:    "SelectContextFilteringBothByNumberAndTextField",
			filters: createFilters(newFilter("random_number", "2"), newFilter("name", "Resource3")),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
			considerOrder: false,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testSelectContext(r))
	}
}

func TestSelectContextWithSorters(t *testing.T) {
	dbHolder := NewTestDBHolder("db_usql_repository_test_select_context")
	dbHolder.Reset()

	filtersMap := map[string]usqlFilters.Filter{
		"id":            usqlFilters.TextField("id"),
		"name":          usqlFilters.TextField("name"),
		"random_number": usqlFilters.NumField("random_number"),
	}

	sortersMap := map[string]usqlFilters.Sorter{
		"sort": usqlFilters.Sort("name", "random_number"),
	}

	r := NewDBRepository[*Resource](dbHolder.DBHolder, filtersMap, sortersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:    "SelectContextSortingByTextFieldNameAsc",
			filters: createFilter("sort", "name"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3, r4},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextSortingByTextFieldNameDesc",
			filters: createFilter("sort", "-name"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3, r4, r2, r1},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextSortingByFieldRandomNumberAsc",
			filters: createFilter("sort", "random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r4, r1, r2, r3},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextSortingByNumFieldRandomNumberDesc",
			filters: createFilter("sort", "-random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r2, r3, r1, r4},
			},
			considerOrder: true,
		},
		{
			name: "SelectContextSortingByNumFieldRandomNumberDescAnLimitTwo",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("limit", "2"),
			),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     2,
				Limit:     2,
				Offset:    0,
				Resources: []*Resource{r2, r3},
			},
			considerOrder: true,
		},
		{
			name: "SelectContextSortingByNumFieldRandomNumberDescAnOffset2",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("offset", "2"),
			),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     2,
				Limit:     10,
				Offset:    2,
				Resources: []*Resource{r1, r4},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextSortingByNameAscAndRandomNumberAsc",
			filters: createFilter("sort", "name", "random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r4, r3},
			},
			considerOrder: true,
		},
		{
			name:    "SelectContextSortingByNameAscAndRandomNumberDesc",
			filters: createFilter("sort", "name", "-random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     4,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r1, r2, r3, r4},
			},
			considerOrder: true,
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testSelectContext(r))
	}
}

func TestGetContext(t *testing.T) {
	dbHolder := NewTestDBHolder("db_usql_repository_test_get_context")
	dbHolder.Reset()

	filtersMap := map[string]usqlFilters.Filter{
		"id":            usqlFilters.TextField("id"),
		"name":          usqlFilters.TextField("name"),
		"random_number": usqlFilters.NumField("random_number"),
	}

	sortersMap := map[string]usqlFilters.Sorter{
		"sort": usqlFilters.Sort("name", "random_number"),
	}

	r := NewDBRepository[*Resource](dbHolder.DBHolder, filtersMap, sortersMap)
	createResourceTable(t, r)

	r1, r2, r3, r4 := populateDB(context.Background(), t, r)

	tests := []findTest{
		{
			name:    "GetContextWithoutFilters",
			filters: url.Values{},
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextFilteringByTextFieldName",
			filters: createFilter("name", "Resource1"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextFilteringByTextFieldID",
			filters: createFilter("id", "5ceff18d-9039-44b5-a5d3-3d99653f4603"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r3},
			},
		},
		{
			name:    "GetContextFilteringByMultipleValuesInNumberField",
			filters: createFilter("name", "Resource1", "Resource2"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextFilteringByMultipleValuesInNumberField",
			filters: createFilter("random_number", "2", "0"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r2},
			},
		},
		{
			name:    "GetContextFilteringBothByNumberAndTextField",
			filters: createFilters(newFilter("random_number", "2"), newFilter("name", "Resource3")),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Limit:     10,
				Offset:    0,
				Resources: []*Resource{r3},
			},
		},
		{
			name:    "GetContextSortingByTextFieldNameAsc",
			filters: createFilter("sort", "name"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextSortingByTextFieldNameDesc",
			filters: createFilter("sort", "-name"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r3},
			},
		},
		{
			name:    "GetContextSortingByFieldRandomNumberAsc",
			filters: createFilter("sort", "random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r4},
			},
		},
		{
			name:    "GetContextSortingByNumFieldRandomNumberDesc",
			filters: createFilter("sort", "-random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r2},
			},
		},
		{
			name: "GetContextSortingByNumFieldRandomNumberDescAnLimitTwo",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("limit", "2"),
			),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r2},
			},
		},
		{
			name: "GetContextSortingByNumFieldRandomNumberDescAnOffset2",
			filters: createFilters(
				newFilter("sort", "-random_number"),
				newFilter("offset", "2"),
			),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextSortingByNameAscAndRandomNumberAsc",
			filters: createFilter("sort", "name", "random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
		{
			name:    "GetContextSortingByNameAscAndRandomNumberDesc",
			filters: createFilter("sort", "name", "-random_number"),
			expected: &udatabase.ResourcePage[*Resource]{
				Total:     1,
				Resources: []*Resource{r1},
			},
		},
	}

	for _, ft := range tests {
		t.Run(ft.name, ft.testGetContext(r))
	}
}

func populateDB(ctx context.Context, t testing.TB, r *DBrepository[*Resource]) (r1, r2, r3, r4 *Resource) {
	ctx, err := r.Begin(ctx)
	require.NoError(t, err)
	defer r.Rollback(ctx)
	r1 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4601",
		Name:         "Resource1",
		RandomNumber: 1,
	}
	require.NoError(t, save(ctx, r, r1))

	r2 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4602",
		Name:         "Resource2",
		RandomNumber: 2,
	}
	require.NoError(t, save(ctx, r, r2))

	r3 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4603",
		Name:         "Resource3",
		RandomNumber: 2,
	}
	require.NoError(t, save(ctx, r, r3))

	r4 = &Resource{
		ID:           "5ceff18d-9039-44b5-a5d3-3d99653f4604",
		Name:         "Resource3",
		RandomNumber: 0,
	}
	require.NoError(t, save(ctx, r, r4))

	err = r.Commit(ctx)
	require.NoError(t, err)
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
	expected      *udatabase.ResourcePage[*Resource]
	considerOrder bool
}

func (ft *findTest) testSelectContext(r *DBrepository[*Resource]) func(*testing.T) {
	return func(t *testing.T) {
		// ARRANGE
		expectedResources := ft.expected.Resources
		require.EqualValues(t, len(expectedResources), ft.expected.Total, "Expected Total and Resources doesn't match.")

		// ACT
		query := "SELECT id, name, random_number as RandomNumber FROM resources"
		rp, err := r.SelectContext(context.Background(), r.GetDBInstance(), query, ft.filters)

		// ASSERT
		require.NoError(t, err)
		require.EqualValues(t, ft.expected.Total, rp.Total)
		require.EqualValues(t, ft.expected.Total, len(rp.Resources))
		require.EqualValues(t, ft.expected.Offset, rp.Offset)
		require.EqualValues(t, ft.expected.Limit, rp.Limit)

		if ft.considerOrder {
			testhelper.RequireEqual(t, expectedResources, rp.Resources)
		} else {
			for _, expRes := range expectedResources {
				require.Contains(t, rp.Resources, expRes)
			}
		}
	}
}

func (ft *findTest) testGetContext(r *DBrepository[*Resource]) func(*testing.T) {
	return func(t *testing.T) {
		// ARRANGE
		expectedResources := ft.expected.Resources
		require.EqualValues(t, 1, ft.expected.Total, "Total must be 1.")
		require.EqualValues(t, len(expectedResources), ft.expected.Total, "Expected Total and Resources doesn't match.")

		// ACT
		var resource Resource
		query := "SELECT id, name, random_number as RandomNumber FROM resources"
		actual, err := r.GetContext(context.Background(), r.GetDBInstance(), &resource, query, ft.filters)

		// ASSERT
		require.NoError(t, err)
		testhelper.RequireEqual(t, *expectedResources[0], resource)
		testhelper.RequireEqual(t, *expectedResources[0], actual)
	}
}

func save(ctx context.Context, r *DBrepository[*Resource], res *Resource) error {
	tx := r.GetTransaction(ctx)
	result, err := tx.NamedExec("INSERT INTO resources (id, name, random_number) VALUES (:id, :name, :randomnumber);", res)
	return r.HandleSaveOrUpdateError(result, err)
}

func findByID(r *DBrepository[*Resource], id string, dst *Resource) error {
	err := r.db.db.Get(dst, "SELECT id, name, random_number as RandomNumber FROM resources WHERE id=$1", id)
	return r.HandleSearchError(err)
}

func findList(r *DBrepository[*Resource], dst *[]*Resource, limit, offset string) error {
	query := "SELECT id, name, random_number as RandomNumber FROM resources"
	query, _, err := filters.ApplyLimit(query, limit)
	if err != nil {
		return err
	}
	query, _, err = filters.ApplyOffset(query, offset)
	if err != nil {
		return err
	}

	return r.db.db.Select(dst, query)
}
