package dao_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/postgres/postgrestest"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/models/migrations"
)

func TestPgItemCreate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		description string

		expect *dao.Item
	}{
		{
			name: "Success",

			description: "test description",

			expect: &dao.Item{
				ID: newItemID, Name: "test item", Description: "test description", CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "Success/NoDescription",

			expect: &dao.Item{ID: newItemID, Name: "test item", CreatedAt: now, UpdatedAt: now},
		},
	}

	daoItemCreate := dao.NewPgItemCreate()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				item, err := daoItemCreate.Exec(ctx, &dao.ItemCreateRequest{
					ID: newItemID, Name: "test item", Description: testCase.description, Now: now,
				})
				require.NoError(t, err)
				require.Equal(t, testCase.expect, item)
			})
		})
	}
}
