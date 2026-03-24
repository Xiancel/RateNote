package item

import (
	"RateNote/internal/entity"
	"context"

	"github.com/google/uuid"
)

type ItemService interface {
	GetItem(ctx context.Context, id uuid.UUID) (*entity.Item, error)
	ListItem(ctx context.Context, filter ItemFilter) (*ItemListResponse, error)
	AddItem(ctx context.Context, req CreateItemRequest) (*entity.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req UpdateItemRequest) (*entity.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}
