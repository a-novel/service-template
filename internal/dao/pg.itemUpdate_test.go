package dao_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/postgres/postgrestest"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/models/migrations"
)

func TestPgItemUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		id uuid.UUID

		expect    *dao.Item
		expectErr error
	}{
		{
			name: "Success",

			id: itemID,

			// The creation time stays.
			expect: &dao.Item{
				ID: itemID, Name: "updated item", Description: "updated description", CreatedAt: createdAt, UpdatedAt: now,
			},
		},
		{
			name: "Error/NotFound",

			id: absentID,

			expectErr: dao.ErrItemNotFound,
		},
	}

	daoItemUpdate := dao.NewPgItemUpdate()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				insertItems(ctx, t)

				item, err := daoItemUpdate.Exec(ctx, &dao.ItemUpdateRequest{
					ID: testCase.id, Name: "updated item", Description: "updated description", Now: now,
				})
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expect, item)
			})
		})
	}
}
