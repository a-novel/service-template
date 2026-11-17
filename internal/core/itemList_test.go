package core_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/core"
	coremocks "github.com/a-novel/service-template/internal/core/mocks"
	"github.com/a-novel/service-template/internal/dao"
)

func TestItemList(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		request *core.ItemListRequest

		daoRequest *dao.ItemListRequest
		daoMock    *daoMock[[]*dao.ItemListRow]

		expect    *core.ItemPage
		expectErr error
	}{
		{
			name: "Success",

			request: &core.ItemListRequest{Limit: 10, Offset: 2},

			daoRequest: &dao.ItemListRequest{Limit: 10, Offset: 2},
			daoMock:    &daoMock[[]*dao.ItemListRow]{resp: []*dao.ItemListRow{{Item: *daoItem, Total: 4}}},

			expect: &core.ItemPage{Items: []*core.Item{coreItem}, Total: 4},
		},
		{
			// A page past the end has no row to carry the total.
			name: "Success/Empty",

			request: &core.ItemListRequest{Limit: 10, Offset: 50},

			daoRequest: &dao.ItemListRequest{Limit: 10, Offset: 50},
			daoMock:    &daoMock[[]*dao.ItemListRow]{resp: []*dao.ItemListRow{}},

			expect: &core.ItemPage{Items: []*core.Item{}},
		},
		{
			name: "Success/LimitDefaulted",

			request: &core.ItemListRequest{},

			daoRequest: &dao.ItemListRequest{Limit: core.ItemListLimitDefault},
			daoMock:    &daoMock[[]*dao.ItemListRow]{resp: []*dao.ItemListRow{}},

			expect: &core.ItemPage{Items: []*core.Item{}},
		},
		{
			name: "Error/NegativeLimit",

			request: &core.ItemListRequest{Limit: -1},

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/LimitTooHigh",

			request: &core.ItemListRequest{Limit: 101},

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/NegativeOffset",

			request: &core.ItemListRequest{Offset: -1},

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/Dao",

			request: &core.ItemListRequest{Limit: 10},

			daoRequest: &dao.ItemListRequest{Limit: 10},
			daoMock:    &daoMock[[]*dao.ItemListRow]{err: errFoo},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			itemDao := coremocks.NewMockItemListDao(t)

			if testCase.daoMock != nil {
				itemDao.EXPECT().
					Exec(mock.Anything, testCase.daoRequest).
					Return(testCase.daoMock.resp, testCase.daoMock.err)
			}

			page, err := core.NewItemList(itemDao).Exec(t.Context(), testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, page)
		})
	}
}
