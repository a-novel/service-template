package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/postgres"
	"github.com/a-novel-kit/golib/postgres/postgrestest"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/models/migrations"
)

func TestPgItemList(t *testing.T) {
	t.Parallel()

	one := dao.Item{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:      "item one",
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	two := dao.Item{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Name:      "item two",
		CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	three := dao.Item{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Name:      "item three",
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
	}

	row := func(item dao.Item) *dao.ItemListRow {
		return &dao.ItemListRow{Item: item, Total: 3}
	}

	testCases := []struct {
		name string

		fixtures []dao.Item

		request *dao.ItemListRequest

		expect    []*dao.ItemListRow
		expectErr error
	}{
		{
			name: "Success",

			fixtures: []dao.Item{one, two, three},

			request: &dao.ItemListRequest{Limit: 10},

			expect: []*dao.ItemListRow{row(three), row(two), row(one)},
		},
		{
			name: "Success/Empty",

			request: &dao.ItemListRequest{Limit: 10},
		},
		{
			name: "Success/Page",

			fixtures: []dao.Item{one, two, three},

			request: &dao.ItemListRequest{Limit: 1, Offset: 1},

			expect: []*dao.ItemListRow{row(two)},
		},
		{
			name: "Success/PastEnd",

			fixtures: []dao.Item{one, two, three},

			request: &dao.ItemListRequest{Limit: 10, Offset: 3},
		},
	}

	daoItemList := dao.NewPgItemList()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			postgrestest.RunDBTest(t, configtest.PostgresPreset, migrations.Migrations, func(ctx context.Context, t *testing.T) {
				t.Helper()

				db, err := postgres.GetContext(ctx)
				require.NoError(t, err)

				if len(testCase.fixtures) > 0 {
					_, err = db.NewInsert().Model(&testCase.fixtures).Exec(ctx)
					require.NoError(t, err)
				}

				result, err := daoItemList.Exec(ctx, testCase.request)
				require.ErrorIs(t, err, testCase.expectErr)
				require.Equal(t, testCase.expect, result)
			})
		})
	}
}
