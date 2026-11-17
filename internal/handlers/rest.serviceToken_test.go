package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/servicetoken"
	"github.com/a-novel-kit/golib/servicetoken/servicetokentest"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/handlers"
)

func TestRestServiceToken(t *testing.T) {
	t.Parallel()

	issuer := servicetokentest.NewIssuer(t)
	tokenConfig := servicetoken.Config{
		Issuer:   issuer.URL,
		Audience: serviceTokenAudience,
		Callers:  []string{serviceTokenCaller},
	}

	testCases := []struct {
		name string

		tokenConfig    servicetoken.Config
		issuerOptional bool
		token          string

		expectStatus int
		expectErr    error
	}{
		{
			name: "Success",

			tokenConfig: tokenConfig,
			token:       issuer.Sign(t, issuer.Claims(serviceTokenCaller, serviceTokenAudience)),

			expectStatus: http.StatusNoContent,
		},
		{
			// A local stack opts out, and the server takes every call.
			name: "Success/IssuerOptional",

			issuerOptional: true,

			expectStatus: http.StatusNoContent,
		},
		{
			// An issuer keeps the gate on whatever the opt-out says.
			name: "Error/IssuerOptionalWithIssuer",

			tokenConfig:    tokenConfig,
			issuerOptional: true,

			expectStatus: http.StatusUnauthorized,
		},
		{
			// A deployment that forgets the issuer stops at startup rather than admit every caller.
			name: "Error/IssuerRequired",

			expectErr: handlers.ErrServiceTokenIssuerRequired,
		},
		{
			name: "Error/Missing",

			tokenConfig: tokenConfig,

			expectStatus: http.StatusUnauthorized,
		},
	}

	// The router answers 204, so a 401 proves the middleware refused the request before it.
	router := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/v0/items", nil)
			if testCase.token != "" {
				request.Header.Set(servicetoken.Header, testCase.token)
			}

			withServiceToken, err := handlers.RestServiceToken(
				testCase.tokenConfig, testCase.issuerOptional, &http.Client{}, config.LoggerDev, "/v0/ping",
			)
			require.ErrorIs(t, err, testCase.expectErr)

			if testCase.expectErr != nil {
				return
			}

			w := httptest.NewRecorder()
			withServiceToken(router).ServeHTTP(w, request)

			require.Equal(t, testCase.expectStatus, w.Code)
		})
	}
}
