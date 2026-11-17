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

		expect    *dao.Item
		expectErr error
		// expectCount is how many items the table holds after the call.
		expectCount int64
	}{
		{
			name: "Success",

			id: itemID,

			expect:      &itemOne,
			expectCount: 2,
		},
		{
			name: "Error/NotFound",

			id: absentID,

			expectErr:   dao.ErrItemNotFound,
			expectCount: 3,
		},
	}

	daoItemDelete := dao.NewPgItemDelete()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				insertItems(ctx, t)

				item, err := daoItemDelete.Exec(ctx, &dao.ItemDeleteRequest{ID: testCase.id})
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expect, item)

				db, err := postgres.GetContext(ctx)
				require.NoError(t, err)

				count, err := db.NewSelect().Model((*dao.Item)(nil)).Count(ctx)
				require.NoError(t, err)
				require.Equal(t, testCase.expectCount, count)
			})
		})
	}
}
