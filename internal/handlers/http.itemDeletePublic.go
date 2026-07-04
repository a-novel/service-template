package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/dao"
)

// ItemDeletePublicService is the core operation the delete endpoint delegates to.
type ItemDeletePublicService interface {
	Exec(ctx context.Context, request *core.ItemDeleteRequest) (*core.Item, error)
}

// ItemDeletePublicRequest carries the query parameters accepted by the delete endpoint.
type ItemDeletePublicRequest struct {
	ID uuid.UUID `schema:"id"`
}

// ItemDeletePublic is the HTTP handler for the delete-item endpoint. It responds
// with the deleted item.
type ItemDeletePublic struct {
	service ItemDeletePublicService
	logger  logging.Log
}

func NewItemDeletePublic(service ItemDeletePublicService, logger logging.Log) *ItemDeletePublic {
	return &ItemDeletePublic{service: service, logger: logger}
}

func (handler *ItemDeletePublic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer().Start(r.Context(), "rest.ItemDeletePublic")
	defer span.End()

	var request ItemDeletePublicRequest

	err := muxDecoder.Decode(&request, r.URL.Query())
	if err != nil {
		httpf.HandleError(ctx, handler.logger, w, span, httpf.ErrMap{nil: http.StatusBadRequest}, err)

		return
	}

	item, err := handler.service.Exec(ctx, &core.ItemDeleteRequest{ID: request.ID})
	if err != nil {
		httpf.HandleError(ctx, handler.logger, w, span, httpf.ErrMap{
			dao.ErrItemDeleteNotFound: http.StatusNotFound,
			core.ErrInvalidRequest:    http.StatusBadRequest,
		}, err)

		return
	}

	httpf.SendJSON(ctx, w, span, loadItem(item))
}
