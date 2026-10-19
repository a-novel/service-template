package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

//go:embed pg.itemGet.sql
var itemGetQuery string

// ErrItemGetNotFound is returned by [PgItemGet.Exec] when no item matches the requested
// ID. It is joined onto the underlying sql.ErrNoRows so callers can branch on it with
// errors.Is.
var ErrItemGetNotFound = errors.New("item not found")

// ItemGetRequest is the input to [PgItemGet.Exec].
type ItemGetRequest struct {
	ID uuid.UUID
}

// PgItemGet retrieves an item by its ID.
type PgItemGet struct{}

func NewPgItemGet() *PgItemGet {
	return &PgItemGet{}
}

func (dao *PgItemGet) Exec(ctx context.Context, request *ItemGetRequest) (*Item, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.PgItemGet")
	defer span.End()

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get transaction: %w", err))
	}

	entity := &Item{}

	err = tx.NewRaw(itemGetQuery, request.ID).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, ErrItemGetNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("execute query: %w", err))
	}

	return entity, nil
}
