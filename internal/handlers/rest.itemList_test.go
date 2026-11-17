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

func TestRestItemList(t *testing.T) {
	t.Parallel()

	target := "/v0/items?limit=10&offset=1"

	testCases := []struct {
		name string

		target string

		serviceMock *serviceMock[*core.ItemPage]

		expectStatus   int
		expectResponse any
	}{
		{
			name: "Success",

			target: target,

			serviceMock: &serviceMock[*core.ItemPage]{resp: &core.ItemPage{Items: []*core.Item{coreItem}, Total: 2}},

			expectStatus:   http.StatusOK,
			expectResponse: map[string]any{"items": []any{jsonItem}, "total": float64(2)},
		},
		{
			// A page past the end lists no item as an empty array, never null.
			name: "Success/Empty",

			target: target,

			serviceMock: &serviceMock[*core.ItemPage]{resp: &core.ItemPage{}},

			expectStatus:   http.StatusOK,
			expectResponse: map[string]any{"items": []any{}, "total": float64(0)},
		},
		{
			name: "Error/MalformedOffset",

			target: "/v0/items?offset=first",

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/InvalidRequest",

			target: target,

			serviceMock: &serviceMock[*core.ItemPage]{err: core.ErrInvalidRequest},

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/Internal",

			target: target,

			serviceMock: &serviceMock[*core.ItemPage]{err: errFoo},

			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockRestItemListService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, &core.ItemListRequest{Limit: 10, Offset: 1}).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.target, nil)

			w := httptest.NewRecorder()
			handlers.NewRestItemList(service, config.LoggerDev).ServeHTTP(w, request)

			requireResponse(t, w.Result(), testCase.expectStatus, testCase.expectResponse)
		})
	}
}
