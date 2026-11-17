package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/servicetoken"
	"github.com/a-novel-kit/golib/servicetoken/servicetokentest"

	"github.com/a-novel/service-template/internal/config"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	const (
		audience = "service-template"
		caller   = "system:serviceaccount:agora:caller"
	)

	issuer := servicetokentest.NewIssuer(t)
	gate := servicetoken.Config{Issuer: issuer.URL, Audience: audience, Callers: []string{caller}}

	// The test hands the router no database, so a request the gate admits to the health route answers
	// that it is down.
	const postgresDown = `{"client:postgres":{"status":"down"}}`

	testCases := []struct {
		name string

		serviceToken servicetoken.Config
		target       string
		token        string

		expectStatus int
		expectBody   string
		expectErr    error
	}{
		{
			// A caller's token passes the gate.
			name: "Success",

			serviceToken: gate,
			target:       "/v0/healthcheck",
			token:        issuer.Sign(t, issuer.Claims(caller, audience)),

			expectStatus: http.StatusServiceUnavailable,
			expectBody:   postgresDown,
		},
		{
			// The liveness probe passes the gate without a token.
			name: "Success/Ping",

			serviceToken: gate,
			target:       "/v0/ping",

			expectStatus: http.StatusOK,
			expectBody:   "pong",
		},
		{
			// A local stack opts out, and the request reaches the health route.
			name: "Success/IssuerOptional",

			serviceToken: servicetoken.Config{IssuerOptional: true},
			target:       "/v0/healthcheck",

			expectStatus: http.StatusServiceUnavailable,
			expectBody:   postgresDown,
		},
		{
			// The health route carries no auth of its own, so the refusal can only come from the gate.
			name: "Error/Unauthenticated",

			serviceToken: gate,
			target:       "/v0/healthcheck",

			expectStatus: http.StatusUnauthorized,
		},
		{
			// The gate wraps the whole router: an unknown path is refused before the router answers 404,
			// so a gate mounted inside /v0 fails this case.
			name: "Error/UnknownPath",

			serviceToken: gate,
			target:       "/nowhere",

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

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, testCase.target, nil)
			if testCase.token != "" {
				request.Header.Set(servicetoken.Header, testCase.token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			require.Equal(t, testCase.expectStatus, w.Code)

			if testCase.expectBody != "" {
				require.Equal(t, testCase.expectBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
