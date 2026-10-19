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

// GrpcItemListService returns a page of items bounded by the request's limit and
// offset. The core layer supplies the implementation.
type GrpcItemListService interface {
	Exec(ctx context.Context, request *core.ItemListRequest) ([]*core.Item, error)
}

// GrpcItemList is the gRPC handler for the ItemList RPC.
type GrpcItemList struct {
	protogen.UnimplementedItemListServiceServer

	service GrpcItemListService
}

func NewGrpcItemList(service GrpcItemListService) *GrpcItemList {
	return &GrpcItemList{service: service}
}

// ItemList returns a page of items. A rejected request maps to InvalidArgument;
// any other failure maps to Internal.
func (handler *GrpcItemList) ItemList(
	ctx context.Context, request *protogen.ItemListRequest,
) (*protogen.ItemListResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.ItemList")
	defer span.End()

	items, err := handler.service.Exec(ctx, &core.ItemListRequest{
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetOffset()),
	})
	if errors.Is(err, core.ErrInvalidRequest) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err != nil {
		_ = otel.ReportError(span, err)

		return nil, status.Error(codes.Internal, "internal error")
	}

	protoItems := make([]*protogen.Item, len(items))
	for i, item := range items {
		protoItems[i] = itemToProto(item)
	}

	return &protogen.ItemListResponse{Items: protoItems}, nil
}
