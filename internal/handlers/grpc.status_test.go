package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/handlers"
	"github.com/a-novel/service-template/internal/handlers/protogen"
)

func TestGrpcStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		skipPostgres bool

		expect       *protogen.StatusResponse
		expectStatus codes.Code
	}{
		{
			name: "Success",

			expect: &protogen.StatusResponse{
				Postgres: &protogen.DependencyHealth{
					Status: protogen.DependencyStatus_DEPENDENCY_STATUS_UP,
				},
			},
		},
		{
			// Omitting postgres from the context makes the probe fail. The RPC then answers
			// Unavailable with no report, so a raw dependency error never reaches the caller.
			name: "Error/Unavailable",

			skipPostgres: true,

			expectStatus: codes.Unavailable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			handler := handlers.NewGrpcStatus()

			ctx := t.Context()

			if !testCase.skipPostgres {
				var err error

				ctx, err = postgres.NewContext(ctx, configtest.PostgresPreset)
				require.NoError(t, err)
			}

			res, err := handler.Status(ctx, new(protogen.StatusRequest))
			resSt, ok := status.FromError(err)
			require.True(t, ok, resSt.Code().String())
			require.Equal(
				t,
				testCase.expectStatus, resSt.Code(),
				"expected status code %s, got %s (%v)", testCase.expectStatus, resSt.Code(), err,
			)
			require.Equal(t, testCase.expect, res)
		})
	}
}
