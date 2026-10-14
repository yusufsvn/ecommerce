package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID         uuid.UUID  `db:"order_id" json:"order_id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	PaymentMethod   string     `db:"payment_method" json:"payment_method"`
	Status          string     `db:"status" json:"status"`
	TotalAmount     int64      `db:"total_amount" json:"total_amount"` // Dikkat: Hassas finansal işlemler için 'decimal' kütüphanesi önerilir
	ShippingAddress string     `db:"shipping_address" json:"shipping_address"`
	BillingAddress  string     `db:"billing_address" json:"billing_address"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
}
