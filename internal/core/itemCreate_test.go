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

func TestItemCreate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		itemName string

		daoMock *daoMock[*dao.Item]

		expect    *core.Item
		expectErr error
	}{
		{
			name: "Success",

			itemName: "test item",

			daoMock: &daoMock[*dao.Item]{resp: daoItem},

			expect: coreItem,
		},
		{
			name: "Error/EmptyName",

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/WhitespaceOnlyName",

			itemName: "   ",

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/NameTooLong",

			itemName: strings.Repeat("a", 257),

			expectErr: core.ErrInvalidRequest,
		},
		{
			name: "Error/Dao",

			itemName: "test item",

			daoMock: &daoMock[*dao.Item]{err: errFoo},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			itemDao := coremocks.NewMockItemCreateDao(t)

			if testCase.daoMock != nil {
				itemDao.EXPECT().
					Exec(mock.Anything, mock.MatchedBy(func(request *dao.ItemCreateRequest) bool {
						// A UUIDv7 sorts by creation time, which the list order relies on.
						return assert.Equal(t, uuid.Version(7), request.ID.Version()) &&
							assert.WithinDuration(t, time.Now(), request.Now, time.Minute) &&
							assert.Equal(t, "test item", request.Name) &&
							assert.Equal(t, "test description", request.Description)
					})).
					Return(testCase.daoMock.resp, testCase.daoMock.err)
			}

			item, err := core.NewItemCreate(itemDao).Exec(t.Context(), &core.ItemCreateRequest{
				Name: testCase.itemName, Description: "test description",
			})
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, item)
		})
	}
}
