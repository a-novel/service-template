package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/core"
)

// RestItemUpdateService is the core operation [RestItemUpdate] delegates to.
type RestItemUpdateService interface {
	Exec(ctx context.Context, request *core.ItemUpdateRequest) (*core.Item, error)
}

// RestItemUpdateRequest is the JSON body accepted by the update endpoint.
type RestItemUpdateRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// RestItemUpdate is the HTTP handler for the update-item endpoint. It responds
// with the updated item.
type RestItemUpdate struct {
	service RestItemUpdateService
	logger  logging.Log
}

func NewRestItemUpdate(service RestItemUpdateService, logger logging.Log) *RestItemUpdate {
	return &RestItemUpdate{service: service, logger: logger}
}

func (handler *RestItemUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer().Start(r.Context(), "rest.ItemUpdate")
	defer span.End()

	var request RestItemUpdateRequest

	err := httpf.DecodeJSON(r.Body, &request)
	if err != nil {
		httpf.HandleError(ctx, handler.logger, w, span, httpf.ErrMap{nil: http.StatusBadRequest}, err)

		return
	}

	item, err := handler.service.Exec(ctx, &core.ItemUpdateRequest{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		statuses := httpf.ErrMap{
			core.ErrItemUpdateNotFound: http.StatusNotFound,
			core.ErrInvalidRequest:     http.StatusBadRequest,
		}
		httpf.HandleError(ctx, handler.logger, w, span, statuses, err)

		return
	}

	httpf.SendJSONStatus(ctx, w, span, http.StatusOK, loadItem(item))
}
