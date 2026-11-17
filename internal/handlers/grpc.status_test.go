package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel/service-template/internal/handlers"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

func TestGrpcStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		// postgres hands the probe the test database; without it, the probe meets no database.
		postgres bool

		expect       *templatev0.StatusResponse
		expectStatus codes.Code
	}{
		{
			name:     "Success",
			postgres: true,
			expect: &templatev0.StatusResponse{
				Postgres: &templatev0.DependencyHealth{Status: templatev0.DependencyStatus_DEPENDENCY_STATUS_UP},
			},
			expectStatus: codes.OK,
		},
		{
			name:         "Error/PostgresDown",
			expectStatus: codes.Unavailable,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if testCase.postgres {
				ctx = postgresContext(t)
			}

			response, err := handlers.NewGrpcStatus().Status(ctx, &templatev0.StatusRequest{})
			result, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, testCase.expectStatus, result.Code(), "got %s (%v)", result.Code(), err)
			require.Equal(t, testCase.expect, response)

			if testCase.expectStatus != codes.OK {
				// The message carries no raw dependency error. The fault keeps that error, so the server
				// span names its type.
				require.Equal(t, "service dependencies unavailable", result.Message())
				require.Error(t, errors.Unwrap(err))
			}
		})
	}
}
