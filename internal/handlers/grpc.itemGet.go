package handlers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

// GrpcItemGetService returns the item identified by the request. The core layer
// supplies the implementation.
type GrpcItemGetService interface {
	Exec(ctx context.Context, request *core.ItemGetRequest) (*core.Item, error)
}

// GrpcItemGet is the gRPC handler for the ItemGet RPC.
type GrpcItemGet struct {
	templatev0.UnimplementedItemGetServiceServer

	service GrpcItemGetService
}

func NewGrpcItemGet(service GrpcItemGetService) *GrpcItemGet {
	return &GrpcItemGet{service: service}
}

// ItemGet returns the item with the given ID. A malformed ID or rejected request
// maps to InvalidArgument, a missing item to NotFound, and any other failure to
// Internal.
func (handler *GrpcItemGet) ItemGet(
	ctx context.Context, request *templatev0.ItemGetRequest,
) (*templatev0.ItemGetResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.ItemGet")
	defer span.End()

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item id")
	}

	item, err := handler.service.Exec(ctx, &core.ItemGetRequest{ID: id})
	if errors.Is(err, core.ErrInvalidRequest) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if errors.Is(err, core.ErrItemGetNotFound) {
		return nil, status.Error(codes.NotFound, "item not found")
	}

	if err != nil {
		_ = otel.ReportError(span, err)

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &templatev0.ItemGetResponse{Item: itemToProto(item)}, nil
}
