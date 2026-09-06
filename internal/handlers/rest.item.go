package handlers

import (
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-template/internal/core"
)

// Item is the REST representation of an item returned by the public endpoints.
type Item struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func loadItem(s *core.Item) Item {
	return Item{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// loadItemMap adapts loadItem to lo.Map, ignoring the slice index.
func loadItemMap(item *core.Item, _ int) Item {
	return loadItem(item)
}
