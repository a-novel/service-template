package handlers

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/golib/grpcf"
	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

// GrpcItemListService is the core operation [GrpcItemList] delegates to.
type GrpcItemListService interface {
	Exec(ctx context.Context, request *core.ItemListRequest) (*core.ItemPage, error)
}

// GrpcItemList is the gRPC handler for the ItemList RPC.
type GrpcItemList struct {
	templatev0.UnimplementedItemListServiceServer

	service GrpcItemListService
}

func NewGrpcItemList(service GrpcItemListService) *GrpcItemList {
	return &GrpcItemList{service: service}
}

// ItemList returns a page of items, newest first, with their total.
func (handler *GrpcItemList) ItemList(
	ctx context.Context, request *templatev0.ItemListRequest,
) (*templatev0.ItemListResponse, error) {
	ctx, span := otel.Tracer().Start(ctx, "grpc.ItemList")
	defer span.End()

	serviceRequest := &core.ItemListRequest{
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetOffset()),
	}

	page, err := handler.service.Exec(ctx, serviceRequest)
	if errors.Is(err, core.ErrInvalidRequest) {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err != nil {
		return nil, grpcf.Fault(codes.Internal, "internal error", otel.ReportError(span, err))
	}

	return &templatev0.ItemListResponse{
		Items: lo.Map(page.Items, func(item *core.Item, _ int) *templatev0.Item { return itemToProto(item) }),
		Total: int64(page.Total),
	}, nil
}
