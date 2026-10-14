package models

import "time"

type ProductProductCategory struct {
	ProductID  string     `db:"product_id" json:"product_id"`
	CategoryID string     `db:"category_id" json:"category_id"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
}

