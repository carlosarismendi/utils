package dbrepository

import (
	"context"
	"ddd-hexa/shared/domain"
	"ddd-hexa/shared/infrastructure/dbholder"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	t.Run("SavingValidResource", func(t *testing.T) {
		schemaName := fmt.Sprintf("jamon%v", time.Now().Unix())
		dbHolder := dbholder.NewDBHolder(schemaName)
		r := NewDBRepository(dbHolder)

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
}
