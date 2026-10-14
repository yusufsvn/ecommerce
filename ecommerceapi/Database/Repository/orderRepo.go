package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderFunctions interface {
	GetAll(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error)
	Create(ctx context.Context, order *st.OrderDTO) error
	Delete(ctx context.Context, orderId uuid.UUID) error
	GetOrdersByUserId(ctx context.Context, user_id string) ([]st.GetOrder, error)
	UpdateOrderState(ctx context.Context, newState string, order_id string) error
	GetOrderCount(ctx context.Context) (int, error)
	UpdateMultiState(ctx context.Context, idlist []string, new_state string) error
	GetOrderCountByStatus(ctx context.Context, status string) (int, error)
	GetOrdersByStatus(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error)
	GetOrderByOrderCode(ctx context.Context, order_code string) (*st.GetOrderAddress, error)
	GetOrderCountForWaiting(ctx context.Context) (int, error)
	SearchOrder(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error)
	SearchOrderCount(ctx context.Context, search_query string) (int, error)
}

type OrderTable struct {
	db *pgxpool.Pool
}

func NewOrderStore(dbpool *pgxpool.Pool) *OrderTable {
	return &OrderTable{
		db: dbpool,
	}
}

// for admin
func (s *OrderTable) GetAll(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {

	query := `Select o.order_code,o.status,o.total_amount,o.created_at,u.name,u.surname from Orders o 
				left join users u on u.user_id = o.user_id 
				limit $1 offset $2`

	rows, err := s.db.Query(ctx, query, limit, offset)

	if err != nil {
		log.Printf("orderRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("sipariş sorgu hatası: %w", err)
	}
	orders, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.OrderForAdmin])
	if err1 != nil {
		log.Printf("orderRepo--GetAll--error:%v", err1)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err1)
	}
	return orders, nil
}

func (s *OrderTable) Create(ctx context.Context, order *st.OrderDTO) error {
	now := time.Now()

	var Query string = `Insert into Orders(
	order_id, 
	user_id, 
	status, 
	total_amount, 
	shipping_address,
	billing_address, 
	created_at, updated_at,order_code) Values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := s.db.Exec(ctx, Query, order.Order_id,
		order.User_id, order.Status, order.TotalAmount,
		order.ShippingAddress, order.BillingAddress, now, now, order.OrderCode)

	if err != nil {
		log.Printf("orderRepo--Create--error:%v", err)
		return fmt.Errorf("sipariş oluşturma hatası: %w", err)
	}
	return nil
}

func (s *OrderTable) Delete(ctx context.Context, orderId uuid.UUID) error {
	query := `Delete From Orders where order_id=$1`
	_, err := s.db.Exec(ctx, query, orderId)
	if err != nil {
		log.Printf("orderRepo--Delete--error:%v", err)
		return fmt.Errorf("sipariş silme hatası: %w", err)
	}
	return nil
}

func (s *OrderTable) GetOrdersByUserId(ctx context.Context, user_id string) ([]st.GetOrder, error) {
	query := `select order_code,user_id,total_amount,status,created_at,
	COALESCE((SELECT SUM(quantity) FROM productorders po WHERE po.order_id = o.order_id), 0) AS total_quantity
	from orders o where user_id=$1 order by created_at desc`
	rows, err := s.db.Query(ctx, query, user_id)

	if err != nil {
		log.Printf("db error:GetOrdersByUserId--%v", err)
		return nil, fmt.Errorf("db error")
	}

	datalist, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.GetOrder])

	if err1 != nil {
		log.Printf("OrderRepo--GetOrdersByUserId--%v", err1)
		return nil, fmt.Errorf("Beklenmeyen Bir sorunla ")
	}

	return datalist, nil

}

func (s *OrderTable) UpdateOrderState(ctx context.Context, newState string, order_id string) error {

	query := `Update Orders set status = $1 where order_code=$2`

	_, err := s.db.Exec(ctx, query, newState, order_id)

	if err != nil {
		log.Printf("db error : /OrderRepo/UpdateOrderState %v", err)
		return fmt.Errorf("Beklenmeyen bir hatayla karşılaşıldı")
	}

	return nil
}

func (s *OrderTable) GetOrderCount(ctx context.Context) (int, error) {

	var count int
	query := `select count(*) from orders`

	err := s.db.QueryRow(ctx, query).Scan(&count)

	if err != nil {
		return -1, fmt.Errorf("Beklenmeyen Bir hatayla karşılaşıldı")
	}

	return count, nil

}

func (s *OrderTable) GetOrderCountByStatus(ctx context.Context, status string) (int, error) {

	var count int
	query := `select count(*) from orders where status=$1`

	err := s.db.QueryRow(ctx, query, status).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil

}

func (s *OrderTable) UpdateMultiState(ctx context.Context, idlist []string, new_state string) error {

	query := `update orders set status = $1 where order_code=any($2)`

	_, err := s.db.Exec(ctx, query, new_state, idlist)

	if err != nil {
		log.Printf("toplu update de hata %v", err)
		return err
	}

	return nil
}

func (s *OrderTable) GetOrdersByStatus(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {

	query := `Select o.order_code,o.status,o.total_amount,o.created_at,u.name,u.surname from Orders o 
				left join users u on u.user_id = o.user_id 
				where status=$1 limit $2 offset $3`

	rows, err := s.db.Query(ctx, query, status, limit, offset)

	if err != nil {
		log.Printf("orderRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("sipariş sorgu hatası: %w", err)
	}
	orders, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.OrderForAdmin])
	if err1 != nil {
		log.Printf("orderRepo--GetAll--error:%v", err1)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err1)
	}
	return orders, nil

}

// Get : order_id, shipping and billing address
func (s *OrderTable) GetOrderByOrderCode(ctx context.Context, order_code string) (*st.GetOrderAddress, error) {

	query := "Select order_id,billing_address,shipping_address from orders where order_code=$1"

	rows, err := s.db.Query(ctx, query, order_code)

	if err != nil {
		log.Printf("db error : %v", err)
		return nil, err
	}

	result, err1 := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[st.GetOrderAddress])

	if err1 != nil {
		log.Printf("collect rows error : %v", err1)
		return nil, err1
	}
	return &result, nil

}

func (s *OrderTable) GetOrderCountForWaiting(ctx context.Context) (int, error) {

	var count int

	query := `select count(*) from orders where status='pending'`

	err := s.db.QueryRow(ctx, query).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil

}

func (s *OrderTable) SearchOrder(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error) {

	query := `select order_code,status,total_amount,o.created_at,u.name,u.surname from orders o 
	left join users u on u.user_id = o.user_id 
	where search_vector @@ to_tsquery('simple', $1) limit $2`

	rows, err := s.db.Query(ctx, query, search_query, limit)

	if err != nil {
		log.Printf("orderRepo--SearchOrder--error:%v", err)
		return nil, fmt.Errorf("sipariş arama hatası: %w", err)
	}

	orders, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.OrderForAdmin])
	if err1 != nil {
		log.Printf("orderRepo--SearchOrder--error:%v", err1)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err1)
	}
	return orders, nil

}

func (s *OrderTable) SearchOrderCount(ctx context.Context, search_query string) (int, error) {

	query := `select count(*) from orders where search_vector @@ to_tsquery('simple', $1)`

	var count int

	err := s.db.QueryRow(ctx, query, search_query).Scan(&count)

	if err != nil {
		log.Printf("orderRepo--SearchOrderCount--error:%v", err)
		return -1, fmt.Errorf("sipariş arama sayısı hatası: %w", err)
	}

	return count, nil

}
