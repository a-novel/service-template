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

func TestPgItemGet(t *testing.T) {
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

			expectErr: dao.ErrItemGetNotFound,
		},
	}

	daoItemGet := dao.NewPgItemGet()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				fixtures := insertItems(ctx, t)

				item, err := daoItemGet.Exec(ctx, &dao.ItemGetRequest{ID: testCase.id})
				require.ErrorIs(t, err, testCase.expectErr)

				if !testCase.expectFixture {
					require.Nil(t, item)

					return
				}

				require.Equal(t, fixtures[0], item)
			})
		})
	}
}
