package postgres

import (
	"RateNote/internal/entity"
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *entity.Item) error
	List(ctx context.Context, limit, offset int) ([]*entity.Item, error)
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	UpdateItem(ctx context.Context, item *entity.Item) error
}
