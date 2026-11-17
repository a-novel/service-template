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

func TestGrpcItemDelete(t *testing.T) {
	t.Parallel()

	request := &templatev0.ItemDeleteRequest{Id: itemID.String()}

	testCases := []struct {
		name string

		request *templatev0.ItemDeleteRequest

		serviceMock *serviceMock[*core.Item]

		expect       *templatev0.ItemDeleteResponse
		expectStatus codes.Code
	}{
		{
			name: "Success",

			request: request,

			serviceMock: &serviceMock[*core.Item]{resp: coreItem},

			expect:       &templatev0.ItemDeleteResponse{Item: protoItem},
			expectStatus: codes.OK,
		},
		{
			name: "Error/MalformedID",

			request: &templatev0.ItemDeleteRequest{Id: "not-a-uuid"},

			expectStatus: codes.InvalidArgument,
		},
		{
			name: "Error/InvalidRequest",

			request: request,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrInvalidRequest},

			expectStatus: codes.InvalidArgument,
		},
		{
			name: "Error/NotFound",

			request: request,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrItemDeleteNotFound},

			expectStatus: codes.NotFound,
		},
		{
			name: "Error/Internal",

			request: request,

			serviceMock: &serviceMock[*core.Item]{err: errFoo},

			expectStatus: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockGrpcItemDeleteService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, &core.ItemDeleteRequest{ID: itemID}).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			response, err := handlers.NewGrpcItemDelete(service).ItemDelete(t.Context(), testCase.request)
			requireStatus(t, err, testCase.expectStatus)
			require.Equal(t, testCase.expect, response)
		})
	}
}
