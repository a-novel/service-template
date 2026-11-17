package dao_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/postgres"
	"github.com/a-novel-kit/golib/postgres/postgrestest"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/models/migrations"
)

func TestPgItemDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		id uuid.UUID

		expectFixture bool
		expectErr     error
	}{
		{
			name: "Success",

			id: itemID,

			expectFixture: true,
		},
		{
			name: "Error/NotFound",

			id: absentID,

			expectErr: dao.ErrItemDeleteNotFound,
		},
	}

	daoItemDelete := dao.NewPgItemDelete()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				fixtures := insertItems(ctx, t)

				item, err := daoItemDelete.Exec(ctx, &dao.ItemDeleteRequest{ID: testCase.id})
				require.ErrorIs(t, err, testCase.expectErr)

				db, err := postgres.GetContext(ctx)
				require.NoError(t, err)

				exists, err := db.NewSelect().Model((*dao.Item)(nil)).Where("id = ?", itemID).Exists(ctx)
				require.NoError(t, err)

				if !testCase.expectFixture {
					require.Nil(t, item)
					require.True(t, exists)

					return
				}

				require.Equal(t, fixtures[0], item)
				require.False(t, exists)
			})
		})
	}
}
