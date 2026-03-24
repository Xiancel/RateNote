package item

import "errors"

var (
	ErrItemNameRequired = errors.New("item name is required")
	ErrInvalidRating    = errors.New("invalid rating")
	ErrItemNotFound     = errors.New("item not found")
	ErrNoFields         = errors.New("no fields")
)
