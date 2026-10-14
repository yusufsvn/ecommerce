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

type CartProductFunctions interface {
	GetByCartID(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error)
	Add(ctx context.Context, cartProduct *st.CartProductDTO) error
	Remove(ctx context.Context, cart_id string, product_id string) error
	UpdateQuantity(ctx context.Context, data *st.CartProductDTO) error
	GetAllProductsFromCart(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error)
	RemoveAll(ctx context.Context, cart_id string) error
}

type CartProductTable struct {
	db *pgxpool.Pool
}

func NewCartProductStore(dbpool *pgxpool.Pool) *CartProductTable {
	return &CartProductTable{
		db: dbpool,
	}
}

func (s *CartProductTable) GetByCartID(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error) {
	// Assuming cart_id is comparable in SQL
	rows, err := s.db.Query(ctx, "Select * from CartProducts where cart_id=$1", cartID)
	if err != nil {
		log.Printf("cartProductRepo--GetByCartID--error:%v", err)
		return nil, fmt.Errorf("sepet ürün sorgu hatası: %w", err)
	}
	cartProducts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entities.CartProduct])
	if err != nil {
		log.Printf("cartProductRepo--GetByCartID--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return cartProducts, nil
}

func (s *CartProductTable) Add(ctx context.Context, cartProduct *st.CartProductDTO) error {
	var Query string = `Insert into CartProducts(cart_id, product_id, quantity) Values ($1,$2,$3)`
	_, err := s.db.Exec(ctx, Query,
		cartProduct.Cart_id,
		cartProduct.Product_id,
		cartProduct.Quantity,
	)
	if err != nil {
		log.Printf("cartProductRepo--Add--error:%v", err)
		return fmt.Errorf("sepet ürün oluşturma hatası: %w", err)
	}
	return nil
}

func (s *CartProductTable) Remove(ctx context.Context, cart_id string, product_id string) error {
	query := `Delete From CartProducts where cart_id=$1 AND product_id=$2`
	_, err := s.db.Exec(ctx, query, cart_id, product_id)
	if err != nil {
		log.Printf("cartProductRepo--Remove--error:%v", err)
		return fmt.Errorf("sepet ürün silme hatası: %w", err)
	}
	return nil
}

func (s *CartProductTable) UpdateQuantity(ctx context.Context, data *st.CartProductDTO) error {

	query := "update CartProducts set quantity=$1 where cart_id=$2 and product_id=$3"

	_, err := s.db.Exec(ctx, query, data.Quantity, data.Cart_id, data.Product_id)

	if err != nil {
		log.Printf("cartProductRepo--UpdateQuantity--error:%v", err)
		return fmt.Errorf("sepet güncellenemedi")
	}

	return nil
}

func (s *CartProductTable) GetAllProductsFromCart(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error) {

	query :=
		`select p.product_id,p.product_name,p.price,p.item_type,cp.quantity from cartProducts cp 
		left join products p on p.product_id = cp.product_id
		where cp.cart_id = $1;`

	rows, err := s.db.Query(ctx, query, cart_id)

	if err != nil {
		log.Printf("db query error--cartProductRepo--GetAllProductsFromCart--%v", err)
		return nil, fmt.Errorf("db error")
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.ProductsFromCartDb])

	if err1 != nil {
		log.Printf("collect rows error--cartProductRepo--GetAllproductsfromcart--%v", err1)
		return nil, fmt.Errorf("row error")
	}

	return data, nil
}

func (s *CartProductTable) RemoveAll(ctx context.Context, cart_id string) error {
	query := `Delete From CartProducts where cart_id=$1`
	_, err := s.db.Exec(ctx, query, cart_id)
	if err != nil {
		log.Printf("cartProductRepo--RemoveAll--error:%v", err)
		return fmt.Errorf("sepet ürün silme hatası: %w", err)
	}
	return nil
}
