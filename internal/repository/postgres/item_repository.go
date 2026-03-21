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
	query := `
	DELETE FROM items WHERE id = $1
	`

	res, err := i.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("items not found")
	}

	return nil
}

// GetItemByID implements ItemRepository.
func (i *ItemRepo) GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.Item, error) {
	var item entity.Item

	query := `
	SELECT id,name,comment,rating,image_path,created_at,updated_at
	FROM items
	WHERE id = $1
	`

	err := i.db.GetContext(ctx, &item, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("Items not found: %w", err)
	}

	return &item, nil
}

// List implements ItemRepository.
func (i *ItemRepo) List(ctx context.Context, limit int, offset int) ([]*entity.Item, error) {
	query := `
	SELECT id,name,comment,rating,image_path,created_at,updated_at
	FROM items
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	var items []*entity.Item
	err := i.db.SelectContext(ctx, &items, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return items, nil
}

// UpdateItem implements ItemRepository.
func (i *ItemRepo) UpdateItem(ctx context.Context, item *entity.Item) error {
	query := `
	UPDATE items
	SET name = $1,
		comment = $2,
		rating = $3,
		image_path = $4,
		updated_at = NOW()
	WHERE id = $5
	`

	res, err := i.db.ExecContext(ctx, query,
		item.Name,
		item.Comment,
		item.Rating,
		item.ImagePath,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}
