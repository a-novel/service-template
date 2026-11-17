package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers"
	handlersmocks "github.com/a-novel/service-template/internal/handlers/mocks"
)

func TestRestItemCreate(t *testing.T) {
	t.Parallel()

	body := `{"name": "test item", "description": "test description"}`

	testCases := []struct {
		name string

		body string

		serviceMock *serviceMock[*core.Item]

		expectStatus   int
		expectResponse any
		expectLocation string
	}{
		{
			name: "Success",

			body: body,

			serviceMock: &serviceMock[*core.Item]{resp: coreItem},

			expectStatus:   http.StatusCreated,
			expectResponse: jsonItem,
			expectLocation: "item?id=" + itemID.String(),
		},
		{
			name: "Error/MalformedBody",

			body: "not json",

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/InvalidRequest",

			body: body,

			serviceMock: &serviceMock[*core.Item]{err: core.ErrInvalidRequest},

			expectStatus: http.StatusUnprocessableEntity,
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

			service := handlersmocks.NewMockRestItemCreateService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, &core.ItemCreateRequest{Name: "test item", Description: "test description"}).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			request := httptest.NewRequestWithContext(
				t.Context(), http.MethodPost, "/v0/items", strings.NewReader(testCase.body),
			)

			w := httptest.NewRecorder()
			handlers.NewRestItemCreate(service, config.LoggerDev).ServeHTTP(w, request)

			response := w.Result()
			require.Equal(t, testCase.expectLocation, response.Header.Get("Location"))
			requireResponse(t, response, testCase.expectStatus, testCase.expectResponse)
		})
	}
}
