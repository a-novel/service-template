package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel-kit/golib/otel"

	"github.com/a-novel/service-template/internal/dao"
)

// ItemListDao is the persistence dependency ItemList uses to read a page of items.
type ItemListDao interface {
	Exec(ctx context.Context, request *dao.ItemListRequest) ([]*dao.ItemListRow, error)
}

// ItemListLimitDefault is the page size of an [ItemListRequest] whose Limit is zero.
const ItemListLimitDefault = 20

// ItemListRequest selects a page of items.
type ItemListRequest struct {
	// Limit caps the page at 100 items. It defaults to [ItemListLimitDefault] when zero.
	Limit int `validate:"min=0,max=100"`
	// Offset is how many items, newest first, come before the page.
	Offset int `validate:"min=0"`
}

// ItemList reads the items page by page, newest first. Each page carries the total.
type ItemList struct {
	dao ItemListDao
}

func NewItemList(dao ItemListDao) *ItemList {
	return &ItemList{dao: dao}
}

func (service *ItemList) Exec(ctx context.Context, request *ItemListRequest) (*ItemPage, error) {
	ctx, span := otel.Tracer().Start(ctx, "core.ItemList")
	defer span.End()

	span.SetAttributes(
		attribute.Int("item.limit", request.Limit),
		attribute.Int("item.offset", request.Offset),
	)

	err := validate.Struct(request)
	if err != nil {
		return nil, otel.ReportError(span, errors.Join(err, ErrInvalidRequest))
	}

	if request.Limit == 0 {
		request.Limit = ItemListLimitDefault
	}

	daoRequest := &dao.ItemListRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}

	rows, err := service.dao.Exec(ctx, daoRequest)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list items: %w", err))
	}

	page := &ItemPage{
		Items: lo.Map(rows, func(row *dao.ItemListRow, _ int) *Item { return newItem(&row.Item) }),
	}

	// Every row carries the total, and a page past the end carries none.
	if len(rows) > 0 {
		page.Total = rows[0].Total
	}

	return page, nil
}
