package models

import (
	"time"
)

type Image struct {
	ImageID   string    `json:"image_id" db:"image_id"`
	ImagePath string    `json:"image_path" db:"image_path"`
	ProductID string    `json:"product_id" db:"product_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
