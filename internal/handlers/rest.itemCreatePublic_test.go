package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers"
	handlersmocks "github.com/a-novel/service-template/internal/handlers/mocks"
)

func TestRestItemCreatePublic(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type serviceMock struct {
		req  *core.ItemCreateRequest
		resp *core.Item
		err  error
	}

	testCases := []struct {
		name string

		request *http.Request

		serviceMock *serviceMock

		expectStatus   int
		expectResponse any
	}{
		{
			name: "Success",

			request: httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/items", strings.NewReader(`{
				"name": "test item",
				"description": "test description"
			}`)),

			serviceMock: &serviceMock{
				req: &core.ItemCreateRequest{
					Name:        "test item",
					Description: "test description",
				},
				resp: &core.Item{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Name:        "test item",
					Description: "test description",
					CreatedAt:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expectStatus: http.StatusCreated,
			expectResponse: map[string]any{
				"id":          "00000000-0000-0000-0000-000000000001",
				"name":        "test item",
				"description": "test description",
				"createdAt":   "2021-01-01T00:00:00Z",
				"updatedAt":   "2021-01-01T00:00:00Z",
			},
		},
		{
			name: "Error/InvalidBody",

			request: httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/items", strings.NewReader(`not json`)),

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/EmptyName",

			request: httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/items", strings.NewReader(`{
				"name": ""
			}`)),

			serviceMock: &serviceMock{
				req: &core.ItemCreateRequest{Name: ""},
				err: core.ErrInvalidRequest,
			},

			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/Internal",

			request: httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/items", strings.NewReader(`{
				"name": "test item"
			}`)),

			serviceMock: &serviceMock{
				req: &core.ItemCreateRequest{Name: "test item"},
				err: errFoo,
			},

			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := handlersmocks.NewMockItemCreatePublicService(t)

			if testCase.serviceMock != nil {
				service.EXPECT().
					Exec(mock.Anything, testCase.serviceMock.req).
					Return(testCase.serviceMock.resp, testCase.serviceMock.err)
			}

			handler := handlers.NewItemCreatePublic(service, config.LoggerDevHttp)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, testCase.request)

			res := w.Result()

			require.Equal(t, testCase.expectStatus, res.StatusCode)

			if testCase.expectResponse != nil {
				data, err := io.ReadAll(res.Body)
				require.NoError(t, errors.Join(err, res.Body.Close()))

				var jsonRes any
				require.NoError(t, json.Unmarshal(data, &jsonRes))
				require.Equal(t, testCase.expectResponse, jsonRes)
			}
		})
	}
}

// The 201 response must carry Content-Type: application/json, which openapi.yaml declares.
// It goes through a real server on purpose: httptest.ResponseRecorder's Header() stays live
// after WriteHeader, so a header the client never receives still reads back as set — which is
// exactly why the old WriteHeader-then-SendJSON ordering shipped a text/plain 201 unnoticed.
func TestRestItemCreatePublicSetsJSONContentType(t *testing.T) {
	t.Parallel()

	service := handlersmocks.NewMockItemCreatePublicService(t)
	service.EXPECT().
		Exec(mock.Anything, &core.ItemCreateRequest{Name: "n", Description: "d"}).
		Return(&core.Item{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:        "n",
			Description: "d",
			CreatedAt:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

	server := httptest.NewServer(handlers.NewItemCreatePublic(service, config.LoggerDevHttp))
	t.Cleanup(server.Close)

	req, err := http.NewRequestWithContext(
		t.Context(), http.MethodPost, server.URL, strings.NewReader(`{"name":"n","description":"d"}`),
	)
	require.NoError(t, err)

	res, err := server.Client().Do(req)
	require.NoError(t, err)

	t.Cleanup(func() { _ = res.Body.Close() })

	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"),
		"a 201 that declares application/json must actually send it")
}
