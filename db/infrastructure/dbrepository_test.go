package infrastructure

import (
	"context"
	"ddd-hexa/shared/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	dbHolder := NewDBHolder(&DBConfig{SchemaName: "db_repository_test"})
	r := NewDBRepository(dbHolder)

	t.Run("SavingValidResource", func(t *testing.T) {
		dbHolder.Reset()

		ctx, err := r.BeginTx(context.Background())
		require.NoError(t, err)

		resource := domain.Resource{
			ID: "0ea57dec-5e79-40dc-b971-a52561fcc2c7",
		}

		func() {
			defer r.EndTx(ctx, &err)

			err = r.db.Exec("CREATE TABLE resources (id UUID);").Error
			require.NoError(t, err)

			err = r.Save(ctx, &resource)
			require.NoError(t, err)
		}()

		var actual domain.Resource
		err = r.FindByID(ctx, resource.ID, &actual)
		require.NoError(t, err)
	})

	t.Run("SavingInvalidResource", func(t *testing.T) {
		dbHolder.Reset()

		ctx, err := r.BeginTx(context.Background())
		require.NoError(t, err)

		resource := domain.Resource{
			ID: "INVALID_UUID",
		}

		func() {
			defer r.EndTx(ctx, &err)

			err = r.db.Exec("CREATE TABLE resources (id UUID);").Error
			require.NoError(t, err)

			err = r.Save(ctx, &resource)
			require.Error(t, err)
		}()
	})
}
