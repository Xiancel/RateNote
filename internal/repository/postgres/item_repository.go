package postgres

import (
	database "RateNote/internal/db"
	"RateNote/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *entity.Item) error
	List(ctx context.Context, limit, offset int) ([]*entity.Item, error)
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	UpdateItem(ctx context.Context, item *entity.Item) error
}

type ItemRepo struct {
	db *database.DB
}

func NewItemRepository(db *database.DB) ItemRepository {
	return &ItemRepo{db: db}
}

// AddItem implements ItemRepository.
func (i *ItemRepo) AddItem(ctx context.Context, item *entity.Item) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	query := `
	INSERT INTO items(id,name,comment,rating,image_path,created_at,updated_at)
	VALUES($1,$2,$3,$4,$5,$6,NOW(),NOW())
	`

	_, err := i.db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Comment,
		item.Rating,
		item.ImagePath,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}
	return nil
}

// DeleteItem implements ItemRepository.
func (i *ItemRepo) DeleteItem(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetItemByID implements ItemRepository.
func (i *ItemRepo) GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.Item, error) {
	panic("unimplemented")
}

// List implements ItemRepository.
func (i *ItemRepo) List(ctx context.Context, limit int, offset int) ([]*entity.Item, error) {
	panic("unimplemented")
}

// UpdateItem implements ItemRepository.
func (i *ItemRepo) UpdateItem(ctx context.Context, item *entity.Item) error {
	panic("unimplemented")
}
