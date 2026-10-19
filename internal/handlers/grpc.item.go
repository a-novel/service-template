package handlers

import (
	"time"

	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/handlers/protogen"
)

// itemToProto converts a core item into its protobuf form, encoding timestamps
// as RFC 3339 strings.
func itemToProto(item *core.Item) *protogen.Item {
	return &protogen.Item{
		Id:          item.ID.String(),
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}
