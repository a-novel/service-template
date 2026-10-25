package handlers

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/a-novel/service-template/internal/core"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

// itemToProto converts a core item into its protobuf form.
func itemToProto(item *core.Item) *templatev0.Item {
	return &templatev0.Item{
		Id:          item.ID.String(),
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
	}
}

// parseOptionalItemID parses an item ID that a request may leave empty. Empty parses to uuid.Nil.
func parseOptionalItemID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, nil
	}

	return uuid.Parse(id)
}
