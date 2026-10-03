package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

//go:embed pg.itemUpdate.sql
var itemUpdateQuery string

// ErrItemUpdateNotFound is returned by [ItemUpdate.Exec] when no item matches the
// requested ID. It is joined onto the underlying sql.ErrNoRows so callers can branch
// on it with errors.Is.
var ErrItemUpdateNotFound = errors.New("item not found")

// ItemUpdateRequest is the input to [ItemUpdate.Exec].
type ItemUpdateRequest struct {
	ID          uuid.UUID
	Name        string
	Description string
	Now         time.Time
}

// ItemUpdate modifies the name and description of an existing item.
type ItemUpdate struct{}

func NewItemUpdate() *ItemUpdate {
	return new(ItemUpdate)
}

func (dao *ItemUpdate) Exec(ctx context.Context, request *ItemUpdateRequest) (*Item, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.ItemUpdate")
	defer span.End()

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get transaction: %w", err))
	}

	entity := new(Item)

	err = tx.NewRaw(itemUpdateQuery, request.Name, request.Description, request.Now, request.ID).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, ErrItemUpdateNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("execute query: %w", err))
	}

	return entity, nil
}
