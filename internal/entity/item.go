package entity

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Comment   string    `db:"comment" json:"comment"`
	Rating    float64   `db:"rating" json:"rating"`
	ImagePath string    `db:"image_path" json:"image_path"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
