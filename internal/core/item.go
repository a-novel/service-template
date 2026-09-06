package core

import (
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-template/internal/dao"
)

var (
	// ErrItemGetNotFound is returned when no item matches a get request.
	ErrItemGetNotFound = dao.ErrItemGetNotFound
	// ErrItemDeleteNotFound is returned when no item matches a delete request.
	ErrItemDeleteNotFound = dao.ErrItemDeleteNotFound
	// ErrItemUpdateNotFound is returned when no item matches an update request.
	ErrItemUpdateNotFound = dao.ErrItemUpdateNotFound
)

// Item is the placeholder resource type.
type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newItem(entity *dao.Item) *Item {
	return &Item{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
