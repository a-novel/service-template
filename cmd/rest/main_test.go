package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/servicetoken"

	"github.com/a-novel/service-template/internal/config"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	// The gate refuses a request without a token before it fetches any key, so the issuer never
	// answers.
	gate := servicetoken.Config{
		Issuer:   "https://issuer.test",
		Audience: "service-template",
		Callers:  []string{"system:serviceaccount:agora:caller"},
	}

	testCases := []struct {
		name string

		serviceToken servicetoken.Config
		target       string

		expectStatus int
		expectBody   string
		expectErr    error
	}{
		{
			// The liveness probe passes the gate without a token.
			name: "Success/Ping",

			serviceToken: gate,
			target:       "/v0/ping",

			expectStatus: http.StatusOK,
			expectBody:   "pong",
		},
		{
			// A local stack opts out, and the request reaches the health route. The test hands it no
			// database, so its probe reports the database down.
			name: "Success/IssuerOptional",

			serviceToken: servicetoken.Config{IssuerOptional: true},
			target:       "/v0/healthcheck",

			expectStatus: http.StatusServiceUnavailable,
			expectBody:   `{"client:postgres":{"status":"down"}}`,
		},
		{
			name: "Error/Unauthenticated",

			serviceToken: gate,
			target:       "/v0/healthcheck",

			expectStatus: http.StatusUnauthorized,
		},
		{
			// A deployment that forgets the issuer stops at startup rather than admit every caller.
			name: "Error/IssuerRequired",

			expectErr: servicetoken.ErrIssuerRequired,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.AppPresetDefault
			cfg.ServiceToken = testCase.serviceToken

			router, err := newRouter(cfg)
			require.ErrorIs(t, err, testCase.expectErr)

			if testCase.expectErr != nil {
				return
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.target, nil))

			require.Equal(t, testCase.expectStatus, w.Code)

			if testCase.expectBody != "" {
				require.Equal(t, testCase.expectBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
