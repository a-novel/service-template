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

func TestGrpcItemCreate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		serviceMock *serviceMock[*core.Item]

		expect       *templatev0.ItemCreateResponse
		expectStatus codes.Code
	}{
		{
			name: "Success",

			serviceMock: &serviceMock[*core.Item]{resp: coreItem},

			expect:       &templatev0.ItemCreateResponse{Item: protoItem},
			expectStatus: codes.OK,
		},
		{
			name: "Error/InvalidRequest",

			serviceMock: &serviceMock[*core.Item]{err: core.ErrInvalidRequest},

			expectStatus: codes.InvalidArgument,
		},
		{
			name: "Error/Internal",

			serviceMock: &serviceMock[*core.Item]{err: errFoo},

			expectStatus: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockGrpcItemCreateService(t)
			service.EXPECT().
				Exec(mock.Anything, &core.ItemCreateRequest{Name: "test item", Description: "test description"}).
				Return(testCase.serviceMock.resp, testCase.serviceMock.err)

			response, err := handlers.NewGrpcItemCreate(service).ItemCreate(
				t.Context(), &templatev0.ItemCreateRequest{Name: "test item", Description: "test description"},
			)
			requireStatus(t, err, testCase.expectStatus)
			require.Equal(t, testCase.expect, response)
		})
	}
}
