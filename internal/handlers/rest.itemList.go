package handlers

import (
	"context"
	"net/http"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
)

// RestItemListService is the core operation [RestItemList] delegates to.
type RestItemListService interface {
	Exec(ctx context.Context, request *core.ItemListRequest) ([]*core.Item, error)
}

// RestItemListRequest carries the pagination query parameters accepted by the
// list endpoint.
type RestItemListRequest struct {
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
}

// RestItemList is the HTTP handler for the list-items endpoint.
type RestItemList struct {
	service RestItemListService
	logger  logging.Log
}

func NewRestItemList(service RestItemListService, logger logging.Log) *RestItemList {
	return &RestItemList{service: service, logger: logger}
}

func (handler *RestItemList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer().Start(r.Context(), "rest.ItemList")
	defer span.End()

	var request RestItemListRequest

	err := muxDecoder.Decode(&request, r.URL.Query())
	if err != nil {
		httpf.HandleError(ctx, handler.logger, w, span, httpf.ErrMap{nil: http.StatusBadRequest}, err)

		return
	}

	items, err := handler.service.Exec(ctx, &core.ItemListRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		httpf.HandleError(ctx, handler.logger, w, span, httpf.ErrMap{
			core.ErrInvalidRequest: http.StatusBadRequest,
		}, err)

		return
	}

	httpf.SendJSONStatus(ctx, w, span, http.StatusOK, lo.Map(items, loadItemMap))
}
