package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment struct {
	PaymentID     uuid.UUID       `db:"payment_id" json:"payment_id"`
	OrderID       uuid.UUID       `db:"order_id" json:"order_id"`
	PaymentMethod string          `db:"payment_method" json:"payment_method"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
	Status        string          `db:"status" json:"status"`
}
