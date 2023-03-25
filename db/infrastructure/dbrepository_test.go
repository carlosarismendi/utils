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
		dbHolder.Reset()
		createResourceTable(t, r)

		ctx, err := domain.BeginTx(context.Background(), r)
		require.NoError(t, err)

		resource := Resource{
			ID:           "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
			Name:         "Resource name",
			RandomNumber: 4,
		}

		err = r.Save(ctx, &resource)
		require.NoError(t, err)

		domain.EndTx(ctx, r, &err)

		var actual Resource
		err = r.FindByID(ctx, resource.ID, &actual)
		require.NoError(t, err)
	})

	t.Run("SavingInvalidResource", func(t *testing.T) {
		dbHolder.Reset()
		createResourceTable(t, r)

		ctx, err := domain.BeginTx(context.Background(), r)
		require.NoError(t, err)

		resource := Resource{
			ID:           "INVALID_UUID",
			Name:         "Resource",
			RandomNumber: 2,
		}

		func() {
			defer domain.EndTx(ctx, r, &err)

			err = r.Save(ctx, &resource)
			require.Error(t, err)
		}()
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

	ctx, err := domain.BeginTx(context.Background(), r)
	require.NoError(t, err)

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

	domain.EndTx(ctx, r, nil)

	t.Run("FindingWithoutFilters", func(t *testing.T) {
		var resources []*Resource
		rp, err := r.Find(context.Background(), url.Values{}, &resources)

		require.NoError(t, err)
		require.Equal(t, 3, len(resources))
		require.Equal(t, int64(3), rp.Total)
		require.Equal(t, 3, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Contains(t, resources, r1)
		require.Contains(t, resources, r2)
		require.Contains(t, resources, r3)
	})

	t.Run("FindingFilteringByTextFieldName", func(t *testing.T) {
		v := url.Values{}
		v.Add("name", "Resource1")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 1, len(resources))
		require.Equal(t, int64(1), rp.Total)
		require.Equal(t, 1, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, r1, resources[0])
	})

	t.Run("FindingFilteringByTextFieldID", func(t *testing.T) {
		v := url.Values{}
		v.Add("id", "5ceff18d-9039-44b5-a5d3-3d99653f4603")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 1, len(resources))
		require.Equal(t, int64(1), rp.Total)
		require.Equal(t, 1, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, r3, resources[0])
	})

	t.Run("FindingFilteringByNumberFieldRandomNumber", func(t *testing.T) {
		v := url.Values{}
		v.Add("random_number", "2")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 2, len(resources))
		require.Equal(t, int64(2), rp.Total)
		require.Equal(t, 2, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Contains(t, resources, r2)
		require.Contains(t, resources, r3)
	})

	t.Run("FindingSortingByTextFieldNameAsc", func(t *testing.T) {
		v := url.Values{}
		v.Add("sort", "name")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 3, len(resources))
		require.Equal(t, int64(3), rp.Total)
		require.Equal(t, 3, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, []*Resource{r1, r2, r3}, resources)
	})

	t.Run("FindingSortingByTextFieldNameDesc", func(t *testing.T) {
		v := url.Values{}
		v.Add("sort", "-name")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 3, len(resources))
		require.Equal(t, int64(3), rp.Total)
		require.Equal(t, 3, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, []*Resource{r3, r2, r1}, resources)
	})

	t.Run("FindingSortingByFieldRandomNumberAsc", func(t *testing.T) {
		v := url.Values{}
		v.Add("sort", "random_number")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 3, len(resources))
		require.Equal(t, int64(3), rp.Total)
		require.Equal(t, 3, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, []*Resource{r1, r2, r3}, resources)
	})

	t.Run("FindingSortingByNumFieldRandomNumberDesc", func(t *testing.T) {
		v := url.Values{}
		v.Add("sort", "-random_number")

		var resources []*Resource
		rp, err := r.Find(context.Background(), v, &resources)

		require.NoError(t, err)
		require.Equal(t, 3, len(resources))
		require.Equal(t, int64(3), rp.Total)
		require.Equal(t, 3, len(*rp.Resources.(*[]*Resource)))
		require.Equal(t, int64(0), rp.Offset)
		require.Equal(t, int64(10), rp.Limit)

		require.Equal(t, []*Resource{r2, r3, r1}, resources)
	})
}
