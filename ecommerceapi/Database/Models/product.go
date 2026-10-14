package models

import (
	"time"
)

type Product struct {
	ProductID           string     `db:"product_id" json:"product_id"`
	ProductName         string     `db:"product_name" json:"product_name"`
	Stock               int        `db:"stock" json:"stock"`
	ThumbnailImage      string     `db:"thumbnail_image" json:"thumbnail_image"`
	Description         string     `db:"description" json:"description"`
	Price               int64      `db:"price" json:"price"`
	CreatedAt           *time.Time `db:"created_at" json:"created_at"` // Ürün ilk oluştuğunda null olabilir mi? SQL'de default now() var ama nullable tanımlanmış.
	UpdatedAt           *time.Time `db:"updated_at" json:"updated_at"`
	ItemType            string     `db:"item_type" json:"item_type"`
	TechnicalDetails    []byte     `db:"technical_details" json:"technical_details"`
	ProductCode         string     `db:"product_code" json:"product_code"`
	ProductPublishState string     `db:"product_publish_state" json:"product_publish_state"`
}
