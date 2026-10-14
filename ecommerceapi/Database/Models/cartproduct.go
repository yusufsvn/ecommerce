package models

import (
	"time"

	"github.com/google/uuid"
)

type CartProduct struct {
	CartID    uuid.UUID  `db:"cart_id" json:"cart_id"`
	ProductID uuid.UUID  `db:"product_id" json:"product_id"`
	Quantity  int        `db:"quantity" json:"quantity"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
