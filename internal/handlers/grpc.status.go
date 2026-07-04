package handlers

import (
	"context"

	"github.com/samber/lo"
	"github.com/uptrace/bun"

	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/handlers/protogen"
)

// NewGrpcHealthStatus reports a dependency as up when err is nil, or down with
// the error message attached otherwise.
func NewGrpcHealthStatus(err error) *protogen.DependencyHealth {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return &protogen.DependencyHealth{
		Status: lo.Ternary(
			err == nil,
			protogen.DependencyStatus_DEPENDENCY_STATUS_UP,
			protogen.DependencyStatus_DEPENDENCY_STATUS_DOWN,
		),
		Err: errMsg,
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

// Status probes each dependency and returns its current health.
func (handler *GrpcStatus) Status(ctx context.Context, _ *protogen.StatusRequest) (*protogen.StatusResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.Status")
	defer span.End()

	return &protogen.StatusResponse{
		Postgres: NewGrpcHealthStatus(handler.reportPostgres(ctx)),
	}, nil
}

func (handler *GrpcStatus) reportPostgres(ctx context.Context) error {
	ctx, span := otel.Tracer().Start(ctx, "grpc.Status(reportPostgres)")
	defer span.End()

	pg, err := postgres.GetContext(ctx)
	if err != nil {
		return otel.ReportError(span, err)
	}

	pgdb, ok := pg.(*bun.DB)
	if !ok {
		// In transaction mode the context carries a transaction, not a *bun.DB,
		// so there is no pooled connection to ping. Report healthy rather than fail.
		return nil
	}

	err = pgdb.Ping()
	if err != nil {
		return otel.ReportError(span, err)
	}

	otel.ReportSuccessNoContent(span)

	return nil
}
