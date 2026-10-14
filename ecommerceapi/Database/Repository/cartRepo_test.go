package repository

import (
	"context"
	"testing"

	Customerrors "ecommerceapi/CustomErrors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCartByUserID(t *testing.T) {
	cartRepo := NewCartStore(testDbInstance)

	t.Run("test: get cart info by user id", func(t *testing.T) {

		// arrange
		// create user, create cart
		ctx := context.Background()
		user_id := uuid.NewString()
		createTestUser(t, ctx, cartRepo.db, user_id)
		cartId := createTestCart(t, ctx, cartRepo.db, user_id)
		// act

		cartIdFromResult, err := cartRepo.GetCartByUserID(ctx, user_id)
		if err != nil {
			t.Errorf("Hata: %v", err)
		}

		// assert

		if !assert.Equal(t, cartIdFromResult, cartId) {
			t.Error("test sonuçları beklenenden farklı geldi")
		}
	})

}

func TestCreateCart(t *testing.T) {
	cartRepo := NewCartStore(testDbInstance)

	t.Run("sepet oluşturma testi", func(t *testing.T) {
		// user gerekli

		// arrange
		user_id := uuid.NewString()
		ctx := context.Background()
		createTestUser(t, ctx, cartRepo.db, user_id)
		// act

		cartId, err := cartRepo.Create(ctx, user_id)
		if err != nil {
			t.Errorf("cart oluşturulamadı %v", err)
		}

		// assert
		var userIdFromDB string
		var status string
		cartRepo.db.QueryRow(ctx, "Select user_id,status from carts where cart_id=$1", cartId).Scan(&userIdFromDB, &status)

		if !assert.Equal(t, userIdFromDB, user_id) {
			t.Error("user id ler eşleşmiyor")
		}

		if !assert.Equal(t, status, "active") {
			t.Error("status is not valid must be active")
		}
	})
}

func TestDeleteCart(t *testing.T) {
	cartRepo := NewCartStore(testDbInstance)

	t.Run("sepet silme testi", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.NewString()
		createTestUser(t, ctx, cartRepo.db, userID)
		cartID := createTestCart(t, ctx, cartRepo.db, userID)

		err := cartRepo.Delete(ctx, cartID)
		if err != nil {
			t.Fatalf("cart silme başarısız oldu: %v", err)
		}

		var exists bool
		err = cartRepo.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM carts WHERE cart_id=$1)", cartID).Scan(&exists)
		if err != nil {
			t.Fatalf("db kontrol sorgusu başarısız oldu: %v", err)
		}

		if !assert.False(t, exists, "cart silinmemiş olmalı") {
			t.Error("cart silme testi başarısız")
		}
	})
}

func TestGetCartProductsByUserID(t *testing.T) {
	cartRepo := NewCartStore(testDbInstance)

	t.Run("kullanıcının sepet ürünlerini getirir", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.NewString()
		createTestUser(t, ctx, cartRepo.db, userID)
		cartID := createTestCart(t, ctx, cartRepo.db, userID)
		productID := createTestProduct(t, ctx, cartRepo.db)
		quantity := 3
		createTestCartProduct(t, ctx, cartRepo.db, cartID, productID, quantity)

		products, err := cartRepo.GetCartProductsByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("sepet ürünleri getirilemedi: %v", err)
		}

		assert.Len(t, products, 1)
		if len(products) > 0 {
			assert.Equal(t, cartID, products[0].Cart_id)
			assert.Equal(t, productID, products[0].Product_id)
			assert.Equal(t, quantity, products[0].Quantity)
		}
	})

}

func TestGetCartProductsByUserID2(t *testing.T) {
	cartRepo := NewCartStore(testDbInstance)

	t.Run("sepet boşsa ZeroListResult döner", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.NewString()
		createTestUser(t, ctx, cartRepo.db, userID)
		createTestCart(t, ctx, cartRepo.db, userID)

		products, err := cartRepo.GetCartProductsByUserID(ctx, userID)
		assert.ErrorIs(t, err, Customerrors.ZeroListResult)
		assert.Nil(t, products)
	})
}
