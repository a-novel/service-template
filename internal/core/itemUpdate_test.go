package core_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/core"
	coremocks "github.com/a-novel/service-template/internal/core/mocks"
	"github.com/a-novel/service-template/internal/dao"
)

func TestItemUpdate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		id       uuid.UUID
		itemName string

		daoMock *daoMock[*dao.Item]

		expect    *core.Item
		expectErr error
	}{
		{
			name: "Success",

			id:       itemID,
			itemName: "test item",

			daoMock: &daoMock[*dao.Item]{resp: daoItem},

			expect: coreItem,
		},
		{
			name: "Error/InvalidID",

			id:       uuid.Nil,
			itemName: "test item",

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/EmptyName",

			id: itemID,

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/WhitespaceOnlyName",

			id:       itemID,
			itemName: "   ",

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/NameTooLong",

			id:       itemID,
			itemName: strings.Repeat("a", 257),

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/NotFound",

			id:       itemID,
			itemName: "test item",

			daoMock: &daoMock[*dao.Item]{err: dao.ErrItemUpdateNotFound},

			expectErr: core.ErrItemUpdateNotFound,
		},
		{
			name: "Error/Dao",

			id:       itemID,
			itemName: "test item",

			daoMock: &daoMock[*dao.Item]{err: errFoo},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			itemDao := coremocks.NewMockItemUpdateDao(t)

			if testCase.daoMock != nil {
				itemDao.EXPECT().
					Exec(mock.Anything, mock.MatchedBy(func(request *dao.ItemUpdateRequest) bool {
						return assert.WithinDuration(t, time.Now(), request.Now, time.Minute) &&
							assert.Equal(t, itemID, request.ID) &&
							assert.Equal(t, "test item", request.Name) &&
							assert.Equal(t, "test description", request.Description)
					})).
					Return(testCase.daoMock.resp, testCase.daoMock.err)
			}

			item, err := core.NewItemUpdate(itemDao).Exec(t.Context(), &core.ItemUpdateRequest{
				ID: testCase.id, Name: testCase.itemName, Description: "test description",
			})
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, item)
		})
	}
}
