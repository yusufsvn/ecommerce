package models

import "time"

type Cart struct {
	CartID    string     `db:"cart_id" json:"cart_id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Status    string     `db:"status" json:"status"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
