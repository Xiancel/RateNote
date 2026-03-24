package item

import (
	"RateNote/internal/entity"
	repository "RateNote/internal/repository/postgres"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type service struct {
	itemRepo repository.ItemRepository
}

func NewService(itemRepo repository.ItemRepository) ItemService {
	return &service{itemRepo: itemRepo}
}

// AddItem implements ItemService.
func (s *service) AddItem(ctx context.Context, req CreateItemRequest) (*entity.Item, error) {
	if req.Name == "" {
		return nil, ErrItemNameRequired
	}
	if req.Rating < 0 || req.Rating > 10 {
		return nil, ErrInvalidRating
	}

	var comment *string
	if req.Comment != "" {
		comment = &req.Comment
	}

	var imagePath *string
	if req.ImagePath != "" {
		imagePath = &req.ImagePath
	}

	item := &entity.Item{
		Name:      req.Name,
		Comment:   *comment,
		Rating:    req.Rating,
		ImagePath: *imagePath,
	}

	if err := s.itemRepo.AddItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return item, nil

}

// DeleteItem implements ItemService.
func (s *service) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := s.itemRepo.GetItemByID(ctx, id)
	if err != nil {
		return ErrItemNotFound
	}

	if err := s.itemRepo.DeleteItem(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

// GetItem implements ItemService.
func (s *service) GetItem(ctx context.Context, id uuid.UUID) (*entity.Item, error) {
	item, err := s.itemRepo.GetItemByID(ctx, id)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

// ListItem implements ItemService.
func (s *service) ListItem(ctx context.Context, filter ItemFilter) (*ItemListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	if filter.MinRating != nil && *filter.MinRating < 0 {
		return nil, ErrInvalidRating
	}

	if filter.MaxRating != nil && *filter.MaxRating < 0 {
		return nil, ErrInvalidRating
	}

	items, err := s.itemRepo.List(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list item: %w", err)
	}

	filtered := []*entity.Item{}
	
	for _, i := range items {
		if filter.Name != "" &&
			!strings.Contains(strings.ToLower(i.Name),
				strings.ToLower(filter.Name)) {
			continue
		}

		if filter.MinRating != nil &&
			i.Rating < *filter.MinRating {
			continue
		}

		if filter.MaxRating != nil &&
			i.Rating > *filter.MaxRating {
			continue
		}

		filtered = append(filtered, i)
	}

	resp := &ItemListResponse{
		Items: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateItem implements ItemService.
func (s *service) UpdateItem(ctx context.Context, id uuid.UUID, req UpdateItemRequest) (*entity.Item, error) {
	item, err := s.itemRepo.GetItemByID(ctx, id)
	if err != nil {
		return nil, ErrItemNotFound
	}

	if req.Name == nil && req.Comment == nil && req.Rating == nil && req.ImagePath == nil {
		return nil, ErrNoFields
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, ErrItemNameRequired
		}
		item.Name = *req.Name
	}

	if req.Comment != nil {
		item.Comment = *req.Comment
	}

	if req.Rating != nil {
		if *req.Rating < 0 || *req.Rating > 10 {
			return nil, ErrInvalidRating
		}
		item.Rating = *req.Rating
	}

	if req.ImagePath != nil {
		item.ImagePath = *req.ImagePath
	}

	if err := s.itemRepo.UpdateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return item, nil
}
