package dao

import (
	"context"
	_ "embed"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

//go:embed pg.itemList.sql
var itemListQuery string

// ItemListRequest is the input to [PgItemList.Exec].
type ItemListRequest struct {
	Limit int
	// Offset is how many items, newest first, come before the page.
	Offset int
}

// ItemListRow is one item of a page.
type ItemListRow struct {
	Item

	// Total counts every item, beyond the page.
	Total int `bun:"total"`
}

// PgItemList retrieves a page of items, newest first. Each row carries the total, and a page past
// the end comes back empty.
type PgItemList struct{}

func NewPgItemList() *PgItemList {
	return &PgItemList{}
}

func (dao *PgItemList) Exec(ctx context.Context, request *ItemListRequest) ([]*ItemListRow, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.PgItemList")
	defer span.End()

	db, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres context: %w", err))
	}

	var entities []*ItemListRow

	err = db.NewRaw(itemListQuery, request.Limit, request.Offset).Scan(ctx, &entities)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("execute query: %w", err))
	}

	span.SetAttributes(attribute.Int("items.count", len(entities)))

	return entities, nil
}
