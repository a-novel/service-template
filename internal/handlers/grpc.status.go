package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/handlers/protogen"
)

// NewGrpcHealthStatus returns the health report for a dependency whose probe succeeded. A failed
// probe answers Unavailable instead, so the report never carries a raw dependency error.
func NewGrpcHealthStatus() *protogen.DependencyHealth {
	return &protogen.DependencyHealth{
		Status: protogen.DependencyStatus_DEPENDENCY_STATUS_UP,
	}
}

// GrpcStatus is the gRPC handler for the Status RPC, reporting the health of the
// service's external dependencies.
type GrpcStatus struct {
	protogen.UnimplementedStatusServiceServer
}

func NewGrpcStatus() *GrpcStatus {
	return new(GrpcStatus)
}

// Status probes each dependency and returns their health, or Unavailable when a probe fails.
func (handler *GrpcStatus) Status(ctx context.Context, _ *protogen.StatusRequest) (*protogen.StatusResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.Status")
	defer span.End()

	// reportPostgres records a failure on its own span.
	if handler.reportPostgres(ctx) != nil {
		return nil, status.Error(codes.Unavailable, "service dependencies unavailable")
	}

	return &protogen.StatusResponse{
		Postgres: NewGrpcHealthStatus(),
	}, nil
}

func (handler *GrpcStatus) reportPostgres(ctx context.Context) error {
	ctx, span := otel.Tracer().Start(ctx, "grpc.Status(reportPostgres)")
	defer span.End()

	err := postgres.Health(ctx)
	if err != nil {
		return otel.ReportError(span, err)
	}

	return nil
}
