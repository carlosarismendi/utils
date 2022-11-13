package infrastructure

import (
	"context"
	"testing"

	"github.com/carlosarismendi/dddhelper/db/domain"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	type Resource struct {
		ID string
	}

	dbHolder := NewDBHolder(&DBConfig{SchemaName: "db_repository_test"})
	r := NewDBRepository(dbHolder)

	t.Run("SavingValidResource", func(t *testing.T) {
		dbHolder.Reset()

		ctx, err := domain.BeginTx(context.Background(), r)
		require.NoError(t, err)

		resource := Resource{
			ID: "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
		}

		func() {
			defer domain.EndTx(ctx, r, &err)

			err = r.db.Exec("CREATE TABLE resources (id UUID);").Error
			require.NoError(t, err)

			err = r.Save(ctx, &resource)
			require.NoError(t, err)
		}()

		var actual Resource
		err = r.FindByID(ctx, resource.ID, &actual)
		require.NoError(t, err)
	})

	t.Run("SavingInvalidResource", func(t *testing.T) {
		dbHolder.Reset()

		ctx, err := domain.BeginTx(context.Background(), r)
		require.NoError(t, err)

		resource := Resource{
			ID: "INVALID_UUID",
		}

		func() {
			defer domain.EndTx(ctx, r, &err)

			err = r.db.Exec("CREATE TABLE resources (id UUID);").Error
			require.NoError(t, err)

			err = r.Save(ctx, &resource)
			require.Error(t, err)
		}()
	})
}
