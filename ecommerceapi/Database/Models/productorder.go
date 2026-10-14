package models

type ProductOrder struct {
	ProductID       string `db:"product_id" json:"product_id"`
	OrderID         string `db:"order_id" json:"order_id"`
	Quantity        int    `db:"quantity" json:"quantity"`
	PriceAtPurchase int64  `db:"price_at_purchase" json:"price_at_purchase"`
}
