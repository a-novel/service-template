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

		// seed fills the table with itemOne, itemTwo and itemThree.
		seed    bool
		request *dao.ItemListRequest

		// expect is the page; every row carries the count of the whole table.
		expect []*dao.ItemListRow
	}{
		{
			// Newest first.
			name: "Success",

			seed:    true,
			request: &dao.ItemListRequest{Limit: 10},

			expect: []*dao.ItemListRow{{Item: itemThree, Total: 3}, {Item: itemTwo, Total: 3}, {Item: itemOne, Total: 3}},
		},
		{
			name: "Success/Empty",

			request: &dao.ItemListRequest{Limit: 10},
		},
		{
			name: "Success/Page",

			seed:    true,
			request: &dao.ItemListRequest{Limit: 1, Offset: 1},

			expect: []*dao.ItemListRow{{Item: itemTwo, Total: 3}},
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

				if testCase.seed {
					insertItems(ctx, t)
				}

				rows, err := daoItemList.Exec(ctx, testCase.request)
				require.NoError(t, err)
				require.Equal(t, testCase.expect, rows)
			})
		})
	}
}
