package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// createUniqueUserAndCart testler için izole bir kullanıcı ve sepet üretir
func createUniqueUserAndCart(t *testing.T, ctx context.Context, db *pgxpool.Pool) (string, string) {
	t.Helper()

	userID := uuid.NewString()
	email := fmt.Sprintf("user_%s@test.com", uuid.NewString())

	queryUser := `INSERT INTO users (user_id, password, name, surname, phone, email) VALUES ($1, 'pass123', 'TestName', 'TestSurname', '05551234567', $2)`
	_, err := db.Exec(ctx, queryUser, userID, email)
	if err != nil {
		t.Fatalf("Test kullanıcısı oluşturulamadı: %v", err)
	}

	queryCart := `INSERT INTO carts (user_id, status) VALUES ($1, 'active') RETURNING cart_id`
	var cartID string
	err = db.QueryRow(ctx, queryCart, userID).Scan(&cartID)
	if err != nil {
		t.Fatalf("Test sepeti oluşturulamadı: %v", err)
	}

	return userID, cartID
}

func TestGetByCartID(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("sepetteki ürünleri başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)

		createTestCartProduct(t, ctx, testDbInstance, cartID, productID, 3)

		cartUUID, err := uuid.Parse(cartID)
		assert.NoError(t, err)

		// act
		cartProducts, err := cartProductRepo.GetByCartID(ctx, cartUUID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, cartProducts, 1)
		assert.Equal(t, cartID, cartProducts[0].CartID.String())
		assert.Equal(t, productID, cartProducts[0].ProductID.String())
		assert.Equal(t, 3, cartProducts[0].Quantity)
	})

	t.Run("sepet boşsa boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, emptyCartID := createUniqueUserAndCart(t, ctx, testDbInstance)

		cartUUID, err := uuid.Parse(emptyCartID)
		assert.NoError(t, err)

		// act
		cartProducts, err := cartProductRepo.GetByCartID(ctx, cartUUID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, cartProducts)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		cartProducts, err := cartProductRepo.GetByCartID(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Nil(t, cartProducts)
		assert.Contains(t, err.Error(), "sepet ürün sorgu hatası")
	})
}

func TestAddCartProduct(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("sepete başarıyla ürün eklemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)

		dto := st.CartProductDTO{
			Cart_id:    cartID,
			Product_id: productID,
			Quantity:   2,
		}

		// act
		err := cartProductRepo.Add(ctx, &dto)

		// assert
		assert.NoError(t, err)

		// Veritabanından doğrudan doğrula
		var quantity int
		err = testDbInstance.QueryRow(ctx, "SELECT quantity FROM cartproducts WHERE cart_id = $1 AND product_id = $2", cartID, productID).Scan(&quantity)
		assert.NoError(t, err)
		assert.Equal(t, 2, quantity)
	})

	t.Run("aynı ürünü ikinci kez eklemeye çalıştığında PK kısıtlaması hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)

		dto := st.CartProductDTO{
			Cart_id:    cartID,
			Product_id: productID,
			Quantity:   1,
		}

		err := cartProductRepo.Add(ctx, &dto)
		assert.NoError(t, err)

		// act: Tekrar ekleme denemesi
		err = cartProductRepo.Add(ctx, &dto)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürün oluşturma hatası")
	})

	t.Run("olmayan sepet veya ürün ID'si ile foreign key hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		dto := st.CartProductDTO{
			Cart_id:    uuid.NewString(),
			Product_id: uuid.NewString(),
			Quantity:   1,
		}

		// act
		err := cartProductRepo.Add(ctx, &dto)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürün oluşturma hatası")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dto := st.CartProductDTO{
			Cart_id:    uuid.NewString(),
			Product_id: uuid.NewString(),
			Quantity:   1,
		}

		// act
		err := cartProductRepo.Add(ctx, &dto)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürün oluşturma hatası")
	})
}

func TestRemoveCartProduct(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("sepetteki ürünü başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)

		createTestCartProduct(t, ctx, testDbInstance, cartID, productID, 2)

		// act
		err := cartProductRepo.Remove(ctx, cartID, productID)

		// assert
		assert.NoError(t, err)

		// Veritabanında kalmadığını doğrula
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT COUNT(*) FROM cartproducts WHERE cart_id = $1 AND product_id = $2", cartID, productID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("olmayan bir kaydı silme denemesi hata fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := cartProductRepo.Remove(ctx, uuid.NewString(), uuid.NewString())

		// assert
		assert.NoError(t, err)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := cartProductRepo.Remove(ctx, uuid.NewString(), uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürün silme hatası")
	})
}

func TestUpdateQuantity(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("ürün adedini başarıyla güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)

		createTestCartProduct(t, ctx, testDbInstance, cartID, productID, 1)

		dto := st.CartProductDTO{
			Cart_id:    cartID,
			Product_id: productID,
			Quantity:   5,
		}

		// act
		err := cartProductRepo.UpdateQuantity(ctx, &dto)

		// assert
		assert.NoError(t, err)

		// Veritabanını doğrula
		var actualQuantity int
		err = testDbInstance.QueryRow(ctx, "SELECT quantity FROM cartproducts WHERE cart_id = $1 AND product_id = $2", cartID, productID).Scan(&actualQuantity)
		assert.NoError(t, err)
		assert.Equal(t, 5, actualQuantity)
	})

	t.Run("olmayan bir sepet ürününün adedini güncelleme denemesi hata fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		dto := st.CartProductDTO{
			Cart_id:    uuid.NewString(),
			Product_id: uuid.NewString(),
			Quantity:   10,
		}

		// act
		err := cartProductRepo.UpdateQuantity(ctx, &dto)

		// assert
		assert.NoError(t, err)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dto := st.CartProductDTO{
			Cart_id:    uuid.NewString(),
			Product_id: uuid.NewString(),
			Quantity:   3,
		}

		// act
		err := cartProductRepo.UpdateQuantity(ctx, &dto)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet güncellenemedi")
	})
}

func TestGetAllProductsFromCart(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("sepetteki ürünleri ürün detayları ile birlikte getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		productID := createTestProduct(t, ctx, testDbInstance)
		_, err := testDbInstance.Exec(ctx, "UPDATE products SET item_type = 'PHYSICAL' WHERE product_id = $1", productID)
		assert.NoError(t, err)

		createTestCartProduct(t, ctx, testDbInstance, cartID, productID, 4)

		// act
		items, err := cartProductRepo.GetAllProductsFromCart(ctx, cartID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, productID, items[0].ProductID)
		assert.Equal(t, 4, items[0].Quantity)
		assert.NotEmpty(t, items[0].ProductName)
		assert.Greater(t, items[0].Price, int64(0))
	})

	t.Run("boş sepette boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, emptyCartID := createUniqueUserAndCart(t, ctx, testDbInstance)

		// act
		items, err := cartProductRepo.GetAllProductsFromCart(ctx, emptyCartID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		items, err := cartProductRepo.GetAllProductsFromCart(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestRemoveAllCartProducts(t *testing.T) {
	cartProductRepo := NewCartProductStore(testDbInstance)

	t.Run("sepetteki tüm ürünleri başarıyla temizlemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		_, cartID := createUniqueUserAndCart(t, ctx, testDbInstance)
		prod1 := createTestProduct(t, ctx, testDbInstance)
		prod2 := createTestProduct(t, ctx, testDbInstance)

		createTestCartProduct(t, ctx, testDbInstance, cartID, prod1, 1)
		createTestCartProduct(t, ctx, testDbInstance, cartID, prod2, 3)

		// act
		err := cartProductRepo.RemoveAll(ctx, cartID)

		// assert
		assert.NoError(t, err)

		// Veritabanında o sepete ait ürün kalmadığını doğrula
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT COUNT(*) FROM cartproducts WHERE cart_id = $1", cartID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("olmayan bir sepeti temizleme denemesi hata fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := cartProductRepo.RemoveAll(ctx, uuid.NewString())

		// assert
		assert.NoError(t, err)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := cartProductRepo.RemoveAll(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürün silme hatası")
	})
}
