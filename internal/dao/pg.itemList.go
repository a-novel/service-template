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
	Limit  int
	Offset int
}

// PgItemList retrieves a paginated list of items ordered by creation date (newest first), with the id
// as a tiebreaker so the ordering is total. Offset pagination is only coherent over a total order.
type PgItemList struct{}

func NewPgItemList() *PgItemList {
	return &PgItemList{}
}

func (dao *PgItemList) Exec(ctx context.Context, request *ItemListRequest) ([]*Item, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.PgItemList")
	defer span.End()

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get transaction: %w", err))
	}

	var entities []*Item

	err = tx.NewRaw(itemListQuery, request.Limit, request.Offset).Scan(ctx, &entities)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("execute query: %w", err))
	}

	span.SetAttributes(attribute.Int("items.count", len(entities)))

	return entities, nil
}
