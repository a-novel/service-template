package core

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/dao"
)

// ItemListDao is the persistence dependency ItemList uses to read a page of items.
type ItemListDao interface {
	Exec(ctx context.Context, request *dao.ItemListRequest) ([]*dao.Item, error)
}

// ItemListDefaultSize is the page size of an [ItemListRequest] whose Limit is zero or negative.
const ItemListDefaultSize = 20

// ItemListRequest selects a page of items.
type ItemListRequest struct {
	// Limit caps the page at 100 items. It defaults to [ItemListDefaultSize] when zero or negative.
	Limit int `validate:"max=100"`
	// Offset is the number of newest items to skip.
	Offset int `validate:"min=0"`
}

// ItemList retrieves a paginated list of items.
type ItemList struct {
	dao ItemListDao
}

func NewItemList(dao ItemListDao) *ItemList {
	return &ItemList{dao: dao}
}

func (service *ItemList) Exec(ctx context.Context, request *ItemListRequest) ([]*Item, error) {
	ctx, span := otel.Tracer().Start(ctx, "core.ItemList")
	defer span.End()

	// Resolved locally so the caller's request stays untouched.
	limit := request.Limit
	if limit <= 0 {
		limit = ItemListDefaultSize
	}

	span.SetAttributes(
		attribute.Int("item.limit", limit),
		attribute.Int("item.offset", request.Offset),
	)

	err := validate.Struct(request)
	if err != nil {
		return nil, otel.ReportError(span, errors.Join(err, ErrInvalidRequest))
	}

	daoRequest := &dao.ItemListRequest{
		Limit:  limit,
		Offset: request.Offset,
	}

	entities, err := service.dao.Exec(ctx, daoRequest)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list items: %w", err))
	}

	items := make([]*Item, len(entities))
	for i, entity := range entities {
		items[i] = newItem(entity)
	}

	return items, nil
}
