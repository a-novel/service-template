package core_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/core"
	coremocks "github.com/a-novel/service-template/internal/core/mocks"
	"github.com/a-novel/service-template/internal/dao"
)

func TestItemGet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		id uuid.UUID

		daoMock *daoMock[*dao.Item]

		expect    *core.Item
		expectErr error
	}{
		{
			name: "Success",

			id: itemID,

			daoMock: &daoMock[*dao.Item]{resp: daoItem},

			expect: coreItem,
		},
		{
			name: "Error/InvalidID",

			id: uuid.Nil,

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/NotFound",

			id: itemID,

			daoMock: &daoMock[*dao.Item]{err: dao.ErrItemGetNotFound},

			expectErr: core.ErrItemGetNotFound,
		},
		{
			name: "Error/Dao",

			id: itemID,

			daoMock: &daoMock[*dao.Item]{err: errFoo},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			itemDao := coremocks.NewMockItemGetDao(t)

			if testCase.daoMock != nil {
				itemDao.EXPECT().
					Exec(mock.Anything, &dao.ItemGetRequest{ID: itemID}).
					Return(testCase.daoMock.resp, testCase.daoMock.err)
			}

			item, err := core.NewItemGet(itemDao).Exec(t.Context(), &core.ItemGetRequest{ID: testCase.id})
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, item)
		})
	}
}
