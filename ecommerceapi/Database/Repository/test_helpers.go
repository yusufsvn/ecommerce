// helpers_test.go
package repository

import (
	"context"
	"ecommerceapi/Utils"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Geçici test kullanıcı üretir
func createTestUser(t *testing.T, ctx context.Context, db *pgxpool.Pool, userId string) {
	// Mark for Helper. tek başına test oarak kullanılamaz
	t.Helper()

	email := fmt.Sprintf("test_%s@test.com", uuid.NewString())
	query := `INSERT INTO users (user_id, name, surname, email) VALUES ($1, 'testuser', 'testsurname', $2)`

	_, err := db.Exec(ctx, query, userId, email)
	if err != nil {
		t.Fatalf("Test kullanıcısı oluşturulamadı: %v", err)
	}
}

func createTestCart(t *testing.T, ctx context.Context, db *pgxpool.Pool, userId string) string {
	t.Helper()

	query := `insert into carts (user_id,status) values ($1,$2) returning cart_id`

	var cartId string
	err := db.QueryRow(ctx, query, userId, "active").Scan(&cartId)

	if err != nil {
		t.Fatalf("Test sepeti oluşturulamadı: %v", err)
	}

	return cartId
}

func createTestProduct(t *testing.T, ctx context.Context, db *pgxpool.Pool) string {
	t.Helper()

	query := `insert into products (product_name, stock, thumbnail_image, description, price, product_publish_state,product_code) values ($1, $2, $3, $4, $5, $6,$7) returning product_id`

	// Birden fazla ürün gereken testler için product Code
	productCode, err := Utils.CreateNanoId(8)
	if err != nil {
		t.Fatalf("Test ürün kodu oluşturulamadı: %v", err)
	}

	var productId string
	err = db.QueryRow(ctx, query, "Test Product", 10, "thumbnail.png", "test description", 1000, "unpublished", productCode).Scan(&productId)
	if err != nil {
		t.Fatalf("Test ürün oluşturulamadı: %v", err)
	}
	return productId
}

func createTestCartProduct(t *testing.T, ctx context.Context, db *pgxpool.Pool, cartId string, productId string, quantity int) {
	t.Helper()

	query := `insert into cartproducts (cart_id, product_id, quantity) values ($1, $2, $3)`
	_, err := db.Exec(ctx, query, cartId, productId, quantity)
	if err != nil {
		t.Fatalf("Test cartproduct oluşturulamadı: %v", err)
	}
}
