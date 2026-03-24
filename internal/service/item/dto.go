package item

import "RateNote/internal/entity"

type CreateItemRequest struct {
	Name      string  `json:"name" validate:"required,min=1,max=255"`
	Comment   string  `json:"comment" validate:"max=1000"`
	Rating    float64 `json:"rating" validate:"required,gte=1,lte=10"`
	ImagePath string  `json:"image_path"`
}

type UpdateItemRequest struct {
	Name      *string  `json:"name" validate:"omitempty,min=1,max=255"`
	Comment   *string  `json:"comment" validate:"max=1000"`
	Rating    *float64 `json:"rating" validate:"required,gte=1,lte=10"`
	ImagePath *string  `json:"image_path"`
}
type ItemFilter struct {
	Name      string   `json:"name"`
	MinRating *float64 `json:"min_rating" validate:"omitempty,gte=1,lte=10"`
	MaxRating *float64 `json:"max_rating" validate:"omitempty,gte=1,lte=10"`
	Limit     int      `json:"limit" validate:"gte=0,lte=100"`
	Offset    int      `json:"offset" validate:"gte=0"`
}

type ItemListResponse struct {
	Items  []*entity.Item `json:"items"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
