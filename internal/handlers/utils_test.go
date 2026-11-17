package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/internal/core"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

var (
	// errFoo is the unexpected error a service mock returns, which every handler answers as a server
	// fault.
	errFoo = errors.New("foo")

	itemID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// coreItem is the item the service mocks return. Its creation time carries microseconds, which
	// both transports must keep.
	coreItem = &core.Item{
		ID:          itemID,
		Name:        "test item",
		Description: "test description",
		CreatedAt:   time.Date(2021, 1, 1, 0, 0, 0, 123456000, time.UTC),
		UpdatedAt:   time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	// jsonItem is coreItem as the REST handlers send it.
	jsonItem = map[string]any{
		"id":          itemID.String(),
		"name":        "test item",
		"description": "test description",
		"createdAt":   "2021-01-01T00:00:00.123456Z",
		"updatedAt":   "2021-01-02T00:00:00Z",
	}
	// protoItem is coreItem as the gRPC handlers send it.
	protoItem = &templatev0.Item{
		Id:          itemID.String(),
		Name:        "test item",
		Description: "test description",
		CreatedAt:   &timestamppb.Timestamp{Seconds: 1609459200, Nanos: 123456000},
		UpdatedAt:   &timestamppb.Timestamp{Seconds: 1609545600},
	}
)

const (
	// serviceTokenAudience is the audience the service token tests answer to.
	serviceTokenAudience = "service-template"
	// serviceTokenCaller is the one caller the service token tests allow.
	serviceTokenCaller = "system:serviceaccount:agora:caller"
	// serviceTokenStranger is a caller outside the allowed ones.
	serviceTokenStranger = "system:serviceaccount:agora:stranger"
)

// serviceMock is what the service mock of a handler test case returns. A nil one means the handler
// must not call the service.
type serviceMock[Resp any] struct {
	resp Resp
	err  error
}

// requireResponse checks the status of a REST handler's response, and its JSON body when one is
// expected.
func requireResponse(t *testing.T, response *http.Response, expectStatus int, expectBody any) {
	t.Helper()

	require.Equal(t, expectStatus, response.StatusCode)

	if expectBody == nil {
		return
	}

	data, err := io.ReadAll(response.Body)
	require.NoError(t, errors.Join(err, response.Body.Close()))

	var body any
	require.NoError(t, json.Unmarshal(data, &body))
	require.Equal(t, expectBody, body)
}

// requireStatus checks the gRPC status a handler answers with. A server fault must keep errFoo, which
// every Internal case's service mock returns, as its cause, so the server span names its type.
func requireStatus(t *testing.T, err error, expectStatus codes.Code) {
	t.Helper()

	result, ok := status.FromError(err)
	require.True(t, ok, "not a gRPC status: %v", err)
	require.Equal(t, expectStatus, result.Code(), "got %s (%v)", result.Code(), err)

	if expectStatus == codes.Internal {
		require.ErrorIs(t, err, errFoo)
	}
}

// postgresContext returns the test's context carrying the test database, as the servers hand it to
// a health probe.
func postgresContext(t *testing.T) context.Context {
	t.Helper()

	ctx, err := postgres.NewContext(t.Context(), configtest.PostgresPreset)
	require.NoError(t, err)

	return ctx
}
