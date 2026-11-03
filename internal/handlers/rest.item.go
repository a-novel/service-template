package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

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

// loadItem converts a core item into its REST form.
func loadItem(item *core.Item) Item {
	return Item{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

// ItemPage is the REST representation of a page of items.
type ItemPage struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

// loadItemPage converts a core page of items into its REST form.
func loadItemPage(page *core.ItemPage) ItemPage {
	return ItemPage{
		Items: lo.Map(page.Items, func(item *core.Item, _ int) Item { return loadItem(item) }),
		Total: page.Total,
	}
}
