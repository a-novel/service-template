package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-template/internal/handlers"
)

func TestRestHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		// postgres hands the probe the test database; without it, the probe meets no database.
		postgres bool

		expectStatus int
		expectHealth string
	}{
		{
			name:         "Success",
			postgres:     true,
			expectStatus: http.StatusOK,
			expectHealth: handlers.RestHealthStatusUp,
		},
		{
			name:         "Error/PostgresDown",
			expectStatus: http.StatusServiceUnavailable,
			expectHealth: handlers.RestHealthStatusDown,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.postgres {
				ctx = postgresContext(t)
			}

			w := httptest.NewRecorder()
			handlers.NewRestHealth().ServeHTTP(w, httptest.NewRequestWithContext(ctx, http.MethodGet, "/v0/healthcheck", nil))

			response := w.Result()
			require.Contains(t, response.Header.Get("Content-Type"), "application/json")
			// Exact matching keeps raw dependency errors out of the public report.
			requireResponse(t, response, testCase.expectStatus, map[string]any{
				"client:postgres": map[string]any{"status": testCase.expectHealth},
			})
		})
	}
}
