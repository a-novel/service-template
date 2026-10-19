package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers/protogen"
)

// GrpcItemCreateService creates an item from the request and returns the stored
// record. The core layer supplies the implementation.
type GrpcItemCreateService interface {
	Exec(ctx context.Context, request *core.ItemCreateRequest) (*core.Item, error)
}

// GrpcItemCreate is the gRPC handler for the ItemCreate RPC.
type GrpcItemCreate struct {
	protogen.UnimplementedItemCreateServiceServer

	service GrpcItemCreateService
}

func NewGrpcItemCreate(service GrpcItemCreateService) *GrpcItemCreate {
	return &GrpcItemCreate{service: service}
}

// ItemCreate creates an item and returns it. A rejected request maps to
// InvalidArgument; any other failure maps to Internal.
func (handler *GrpcItemCreate) ItemCreate(
	ctx context.Context, request *protogen.ItemCreateRequest,
) (*protogen.ItemCreateResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.ItemCreate")
	defer span.End()

	item, err := handler.service.Exec(ctx, &core.ItemCreateRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	})
	if errors.Is(err, core.ErrInvalidRequest) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err != nil {
		_ = otel.ReportError(span, err)

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &protogen.ItemCreateResponse{Item: itemToProto(item)}, nil
}
