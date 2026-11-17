package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers"
	handlersmocks "github.com/a-novel/service-template/internal/handlers/mocks"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

func TestGrpcItemList(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		serviceMock *serviceMock[*core.ItemPage]

		expect       *templatev0.ItemListResponse
		expectStatus codes.Code
	}{
		{
			name: "Success",

			serviceMock: &serviceMock[*core.ItemPage]{resp: &core.ItemPage{Items: []*core.Item{coreItem}, Total: 2}},

			expect:       &templatev0.ItemListResponse{Items: []*templatev0.Item{protoItem}, Total: 2},
			expectStatus: codes.OK,
		},
		{
			name: "Error/InvalidRequest",

			serviceMock: &serviceMock[*core.ItemPage]{err: core.ErrInvalidRequest},

			expectStatus: codes.InvalidArgument,
		},
		{
			name: "Error/Internal",

			serviceMock: &serviceMock[*core.ItemPage]{err: errFoo},

			expectStatus: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockGrpcItemListService(t)
			service.EXPECT().
				Exec(mock.Anything, &core.ItemListRequest{Limit: 10, Offset: 1}).
				Return(testCase.serviceMock.resp, testCase.serviceMock.err)

			response, err := handlers.NewGrpcItemList(service).ItemList(
				t.Context(), &templatev0.ItemListRequest{Limit: 10, Offset: 1},
			)
			requireStatus(t, err, testCase.expectStatus)
			require.Equal(t, testCase.expect, response)
		})
	}
}
