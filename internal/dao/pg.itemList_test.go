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

func TestPgItemList(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		// seed fills the table with the three [itemFixtures].
		seed    bool
		request *dao.ItemListRequest

		// expect lists the fixtures the page holds, by their index in [itemFixtures].
		expect []int
	}{
		{
			// Newest first.
			name: "Success",

			seed:    true,
			request: &dao.ItemListRequest{Limit: 10},

			expect: []int{2, 1, 0},
		},
		{
			name: "Success/Empty",

			request: &dao.ItemListRequest{Limit: 10},
		},
		{
			name: "Success/Page",

			seed:    true,
			request: &dao.ItemListRequest{Limit: 1, Offset: 1},

			expect: []int{1},
		},
		{
			name: "Success/PastEnd",

			seed:    true,
			request: &dao.ItemListRequest{Limit: 10, Offset: 3},
		},
	}

	daoItemList := dao.NewPgItemList()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				var fixtures []*dao.Item
				if testCase.seed {
					fixtures = insertItems(ctx, t)
				}

				var expect []*dao.ItemListRow
				// Every row carries the count of the whole table.
				for _, index := range testCase.expect {
					expect = append(expect, &dao.ItemListRow{Item: *fixtures[index], Total: len(fixtures)})
				}

				rows, err := daoItemList.Exec(ctx, testCase.request)
				require.NoError(t, err)
				require.Equal(t, expect, rows)
			})
		})
	}
}
