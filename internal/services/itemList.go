package services

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/dao"
)

type ItemListRepository interface {
	Exec(ctx context.Context, request *dao.ItemListRequest) ([]*dao.Item, error)
}

// ItemListMaxSize caps the number of items returned per page.
const ItemListMaxSize = 100

type ItemListRequest struct {
	Limit  int `validate:"min=1,max=100"`
	Offset int
}

// ItemList retrieves a paginated list of items.
type ItemList struct {
	repository ItemListRepository
}

func NewItemList(repository ItemListRepository) *ItemList {
	return &ItemList{repository: repository}
}

func (service *ItemList) Exec(ctx context.Context, request *ItemListRequest) ([]*Item, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.ItemList")
	defer span.End()

	span.SetAttributes(
		attribute.Int("item.limit", request.Limit),
		attribute.Int("item.offset", request.Offset),
	)

	err := validate.Struct(request)
	if err != nil {
		return nil, otel.ReportError(span, errors.Join(err, ErrInvalidRequest))
	}

	entities, err := service.repository.Exec(ctx, &dao.ItemListRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list items: %w", err))
	}

	span.SetAttributes(attribute.Int("items.count", len(entities)))

	items := make([]*Item, len(entities))
	for i, entity := range entities {
		items[i] = &Item{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		}
	}

	return otel.ReportSuccess(span, items), nil
}
