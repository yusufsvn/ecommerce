package repository

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	structs "ecommerceapi/Structs"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartFunctions interface {
	GetCartByUserID(ctx context.Context, user_id string) (string, error)
	Create(ctx context.Context, user_id string) (string, error)
	Delete(ctx context.Context, cartId string) error
	GetCartProductsByUserID(ctx context.Context, user_id string) ([]structs.CartProductDTO, error)
}

type CartTable struct {
	db *pgxpool.Pool
}

func NewCartStore(dbpool *pgxpool.Pool) *CartTable {
	return &CartTable{
		db: dbpool,
	}
}

func (s *CartTable) GetCartByUserID(ctx context.Context, user_id string) (string, error) {
	rows, err := s.db.Query(ctx, "Select cart_id from Carts where user_id=$1", user_id)
	if err != nil {
		log.Printf("cartRepo--GetCartByUserID--error:%v", err)
		return "", fmt.Errorf("sepet sorgu hatası: %w", err)
	}
	cart, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		log.Printf("cartRepo--GetCartByUserID--error:%v", err)
		return "", fmt.Errorf("from cartRepo GetCartByUserID :satır dönüşüm hatası: %w", err)
	}
	return cart, nil
}

// Returning the CartID
func (s *CartTable) Create(ctx context.Context, user_id string) (string, error) {

	cart := Entities.Cart{
		UserID: user_id,
		Status: "active",
	}

	var Query string = `Insert into Carts(user_id, status) Values ($1,$2) returning cart_id`
	err := s.db.QueryRow(ctx, Query,
		cart.UserID,
		cart.Status,
	).Scan(&cart.CartID)

	if err != nil {
		log.Printf("cartRepo--Create--error:%v", err)
		return "", fmt.Errorf("sepet oluşturma hatası: %w", err)
	}
	return cart.CartID, nil
}

func (s *CartTable) Delete(ctx context.Context, cartId string) error {
	query := `Delete From Carts where cart_id=$1`
	_, err := s.db.Exec(ctx, query, cartId)
	if err != nil {
		log.Printf("cartRepo--Delete--error:%v", err)
		return fmt.Errorf("sepet silme hatası: %w", err)
	}
	return nil
}

func (s *CartTable) GetCartProductsByUserID(ctx context.Context, user_id string) ([]structs.CartProductDTO, error) {

	query := `SELECT p.product_id,c.cart_id,cp.quantity FROM products p 
		JOIN cartproducts cp ON p.product_id = cp.product_id 
		JOIN carts c ON cp.cart_id = c.cart_id WHERE c.user_id = $1`
	rows, err := s.db.Query(ctx, query, user_id)

	if err != nil {
		log.Printf("cartRepo--GetCartProductsByUserID--error:%v", err)
		return nil, fmt.Errorf("sepet sorgu hatası: %w", err)
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[structs.CartProductDTO])
	if err != nil {
		log.Printf("cartRepo--GetCartProductsByUserID--error:%v", err)
		return nil, fmt.Errorf("satır eşleme hatası: %w", err)
	}
	if len(products) == 0 || products == nil {
		return nil, Customerrors.ZeroListResult
	}
	return products, nil
}
