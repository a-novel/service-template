package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers"
	handlersmocks "github.com/a-novel/service-template/internal/handlers/mocks"
)

func TestRestItemUpdate(t *testing.T) {
	t.Parallel()

	body := `{"id": "` + itemID.String() + `", "name": "test item", "description": "test description"}`

	testCases := []struct {
		name string

		body string

		serviceMock *serviceMock[*core.Item]

		expectStatus   int
		expectResponse any
	}{
		{
			name: "Success",

			body: body,

			serviceMock: &serviceMock[*core.Item]{resp: coreItem},

			expectStatus:   http.StatusOK,
			expectResponse: jsonItem,
		},
		{
			name: "Error/MalformedBody",

			body: "not json",

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/MalformedID",

			body: `{"id": "not-a-uuid", "name": "test item"}`,

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/InvalidRequest",

			body: body,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrInvalidRequest},

			expectStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error/NotFound",

			body: body,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrItemNotFound},

			expectStatus: http.StatusNotFound,
		},
		{
			name: "Error/Internal",

			body: body,

			serviceMock: &serviceMock[*core.Item]{err: errFoo},

			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockRestItemUpdateService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, &core.ItemUpdateRequest{
						ID: itemID, Name: "test item", Description: "test description",
					}).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			request := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/v0/item", strings.NewReader(testCase.body))

			w := httptest.NewRecorder()
			handlers.NewRestItemUpdate(service, config.LoggerDev).ServeHTTP(w, request)

			requireResponse(t, w.Result(), testCase.expectStatus, testCase.expectResponse)
		})
	}
}
