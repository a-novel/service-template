package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers"
	handlersmocks "github.com/a-novel/service-template/internal/handlers/mocks"
)

func TestRestItemGet(t *testing.T) {
	t.Parallel()

	target := "/v0/item?id=" + itemID.String()

	testCases := []struct {
		name string

		target string

		serviceMock *serviceMock[*core.Item]

		expectStatus   int
		expectResponse any
	}{
		{
			name: "Success",

			target: target,

			serviceMock: &serviceMock[*core.Item]{resp: coreItem},

			expectStatus:   http.StatusOK,
			expectResponse: jsonItem,
		},
		{
			name: "Error/MalformedID",

			target: "/v0/item?id=not-a-uuid",

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/InvalidRequest",

			target: target,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrInvalidRequest},

			expectStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error/NotFound",

			target: target,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrItemNotFound},

			expectStatus: http.StatusNotFound,
		},
		{
			name: "Error/Internal",

			target: target,

			serviceMock: &serviceMock[*core.Item]{err: errFoo},

			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockRestItemGetService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, &core.ItemGetRequest{ID: itemID}).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.target, nil)

			w := httptest.NewRecorder()
			handlers.NewRestItemGet(service, config.LoggerDev).ServeHTTP(w, request)

			requireResponse(t, w.Result(), testCase.expectStatus, testCase.expectResponse)
		})
	}
}
