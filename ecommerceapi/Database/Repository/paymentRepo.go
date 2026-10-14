package repository

import (
	"context"
	Entities "ecommerceapi/Database/Models"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentFunctions interface {
	GetAll(ctx context.Context) ([]Entities.Payment, error)
	Create(ctx context.Context, payment Entities.Payment) error
	Delete(ctx context.Context, paymentId uuid.UUID) error
}

type PaymentTable struct {
	db *pgxpool.Pool
}

func NewPaymentStore(dbpool *pgxpool.Pool) *PaymentTable {
	return &PaymentTable{
		db: dbpool,
	}
}

func (s *PaymentTable) GetAll(ctx context.Context) ([]Entities.Payment, error) {
	rows, err := s.db.Query(ctx, "Select * from Payments")
	if err != nil {
		log.Printf("paymentRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("ödeme sorgu hatası: %w", err)
	}
	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entities.Payment])
	if err != nil {
		log.Printf("paymentRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return payments, nil
}

func (s *PaymentTable) Create(ctx context.Context, payment *Entities.Payment) error {
	var Query string = `Insert into Payments(payment_id, order_id, payment_method, amount, status) Values ($1,$2,$3,$4,$5)`
	_, err := s.db.Exec(ctx, Query,
		payment.PaymentID,
		payment.OrderID,
		payment.PaymentMethod,
		payment.Amount,
		payment.Status,
	)
	if err != nil {
		log.Printf("paymentRepo--Create--error:%v", err)
		return fmt.Errorf("ödeme oluşturma hatası: %w", err)
	}
	return nil
}

func (s *PaymentTable) Delete(ctx context.Context, paymentId *uuid.UUID) error {
	query := `Delete From Payments where payment_id=$1`
	_, err := s.db.Exec(ctx, query, paymentId)
	if err != nil {
		log.Printf("paymentRepo--Delete--error:%v", err)
		return fmt.Errorf("ödeme silme hatası: %w", err)
	}
	return nil
}
