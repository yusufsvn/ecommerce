package repository

import (
	"context"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductOrderFunctions interface {
	GetByOrderID(ctx context.Context, orderID string) ([]st.OrderFullDetail, error)
	Add(ctx context.Context, item Entities.ProductOrder) error
	AddMany(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error
	Delete(ctx context.Context, product_id uuid.UUID) error
}

type ProductOrderTable struct {
	db *pgxpool.Pool
}

func NewProductOrderStore(dbpool *pgxpool.Pool) *ProductOrderTable {
	return &ProductOrderTable{
		db: dbpool,
	}
}

// add : thumbnail image, product name
func (s *ProductOrderTable) GetByOrderID(ctx context.Context, orderID string) ([]st.OrderFullDetail, error) {

	query := `Select p.product_id,po.quantity,po.price_at_purchase,p.thumbnail_image,p.product_name from ProductOrders po
	inner join products p on p.product_id = po.product_id
	where order_id =$1`

	rows, err := s.db.Query(ctx, query, orderID)
	if err != nil {
		log.Printf("productOrderRepo--GetByOrderID--error:%v", err)
		return nil, fmt.Errorf("sipariş ürün sorgu hatası: %w", err)
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[st.OrderFullDetail])
	if err != nil {
		log.Printf("productOrderRepo--GetByOrderID--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return items, nil
}

func (s *ProductOrderTable) Add(ctx context.Context, item Entities.ProductOrder) error {
	var Query string = `Insert into ProductOrders(product_id, order_id, quantity, price_at_purchase) Values ($1,$2,$3,$4)`
	_, err := s.db.Exec(ctx, Query,
		item.ProductID,
		item.OrderID,
		item.Quantity,
		item.PriceAtPurchase,
	)
	if err != nil {
		log.Printf("productOrderRepo--Add--error:%v", err)
		return fmt.Errorf("sipariş ürün oluşturma hatası: %w", err)
	}
	return nil
}

func (s *ProductOrderTable) AddMany(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {

	// bulk insert
	_, err := s.db.CopyFrom(
		ctx,
		pgx.Identifier{"productorders"},
		[]string{"product_id", "order_id", "quantity", "price_at_purchase"},

		// CopyFromSlice, veri dizisinin uzunluğunu ve her adımda
		// satır verisini döndüren bir callback (geri çağırma) fonksiyonu alır.
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			return []any{
				items[i].ProductID,
				order_id,
				items[i].Quantity,
				items[i].Price,
			}, nil
		}),
	)

	if err != nil {
		log.Printf("db error : productOrders/AddMany: %v", err)
		return fmt.Errorf("db error")
	}

	return nil
}

func (s *ProductOrderTable) Delete(ctx context.Context, product_id uuid.UUID) error {
	query := `Delete from productorders where product_id = $1`
	_, err := s.db.Exec(ctx, query, product_id)
	if err != nil {
		log.Printf("productOrderRepo--Delete--error:%v", err)
		return fmt.Errorf("sipariş ürün silme hatası: %w", err)
	}
	return nil
}
