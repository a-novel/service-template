package core

import (
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-template/internal/dao"
)

// Item is the placeholder resource type.
type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ItemPage is a page of items.
type ItemPage struct {
	Items []*Item
	// Total counts every item, so a client paging through them knows when it holds them all. A page
	// past the end reports 0.
	Total int
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
