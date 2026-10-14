package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createOrderTestUser foreign key kısıtını sağlamak için izole bir kullanıcı oluşturur
func createOrderTestUser(t *testing.T, ctx context.Context, db *pgxpool.Pool) (string, string, string) {
	t.Helper()

	userID := uuid.NewString()
	name := "OrderUserName"
	surname := "OrderUserSurname"
	email := fmt.Sprintf("order_user_%s@test.com", uuid.NewString())
	query := `INSERT INTO users (user_id, password, name, surname, phone, email) VALUES ($1, 'pass123', $2, $3, '05559998877', $4)`

	_, err := db.Exec(ctx, query, userID, name, surname, email)
	if err != nil {
		t.Fatalf("Test kullanıcısı oluşturulamadı: %v", err)
	}

	return userID, name, surname
}

// createOrderHelper testler için izole bir sipariş oluşturur
func createOrderHelper(t *testing.T, ctx context.Context, repo *OrderTable, userID string, status string) *st.OrderDTO {
	t.Helper()

	orderID := uuid.NewString()
	orderCode, err := Utils.CreateNanoId(8)
	require.NoError(t, err)

	dto := &st.OrderDTO{
		Order_id:        orderID,
		User_id:         userID,
		Status:          status,
		TotalAmount:     "2500",
		ShippingAddress: "Test Shipping Address",
		BillingAddress:  "Test Billing Address",
		OrderCode:       orderCode,
	}

	err = repo.Create(ctx, dto)
	require.NoError(t, err)

	return dto
}

func TestCreateOrder(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("başarıyla sipariş oluşturmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		orderID := uuid.NewString()
		orderCode, err := Utils.CreateNanoId(8)
		require.NoError(t, err)

		order := &st.OrderDTO{
			Order_id:        orderID,
			User_id:         userID,
			Status:          "pending",
			TotalAmount:     "5000",
			ShippingAddress: "Kadıköy, İstanbul",
			BillingAddress:  "Kadıköy, İstanbul",
			OrderCode:       orderCode,
		}

		// act
		err = orderRepo.Create(ctx, order)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var dbStatus, dbShipping, dbBilling, dbOrderCode string
		var dbTotalAmount int64
		query := `SELECT status, total_amount, shipping_address, billing_address, order_code FROM orders WHERE order_id = $1`
		err = testDbInstance.QueryRow(ctx, query, orderID).Scan(&dbStatus, &dbTotalAmount, &dbShipping, &dbBilling, &dbOrderCode)

		assert.NoError(t, err)
		assert.Equal(t, "pending", dbStatus)
		assert.Equal(t, int64(5000), dbTotalAmount)
		assert.Equal(t, "Kadıköy, İstanbul", dbShipping)
		assert.Equal(t, "Kadıköy, İstanbul", dbBilling)
		assert.Equal(t, orderCode, dbOrderCode)
	})

	t.Run("olmayan user_id ile oluşturulmaya çalışıldığında foreign key hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		orderCode, _ := Utils.CreateNanoId(8)
		order := &st.OrderDTO{
			Order_id:        uuid.NewString(),
			User_id:         uuid.NewString(),
			Status:          "pending",
			TotalAmount:     "1000",
			ShippingAddress: "Address",
			BillingAddress:  "Address",
			OrderCode:       orderCode,
		}

		// act
		err := orderRepo.Create(ctx, order)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sipariş oluşturma hatası")
	})

	t.Run("aynı order_code ile tekrar sipariş eklenmeye çalışıldığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order1 := createOrderHelper(t, ctx, orderRepo, userID, "pending")

		duplicateOrder := &st.OrderDTO{
			Order_id:        uuid.NewString(),
			User_id:         userID,
			Status:          "pending",
			TotalAmount:     "2000",
			ShippingAddress: "Address",
			BillingAddress:  "Address",
			OrderCode:       order1.OrderCode, // aynı order_code
		}

		// act
		err := orderRepo.Create(ctx, duplicateOrder)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sipariş oluşturma hatası")
	})

	t.Run("bağlam iptal edildiğinde (context canceled) hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		order := &st.OrderDTO{
			Order_id:  uuid.NewString(),
			User_id:   uuid.NewString(),
			OrderCode: "CAN12345",
		}

		// act
		err := orderRepo.Create(ctx, order)

		// assert
		assert.Error(t, err)
	})
}

func TestGetAllOrders(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("admin için sayfalanmış tüm siparişleri başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, userName, userSurname := createOrderTestUser(t, ctx, testDbInstance)
		ord1 := createOrderHelper(t, ctx, orderRepo, userID, "pending")
		ord2 := createOrderHelper(t, ctx, orderRepo, userID, "completed")

		// act
		orders, err := orderRepo.GetAll(ctx, "50", "0")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, orders)

		found1, found2 := false, false
		for _, o := range orders {
			if o.OrderId == ord1.OrderCode {
				found1 = true
				assert.Equal(t, "pending", o.Status)
				assert.Equal(t, userName, o.Name)
				assert.Equal(t, userSurname, o.Surname)
			}
			if o.OrderId == ord2.OrderCode {
				found2 = true
				assert.Equal(t, "completed", o.Status)
			}
		}
		assert.True(t, found1, "ord1 siparişi admin listesinde bulunmalı")
		assert.True(t, found2, "ord2 siparişi admin listesinde bulunmalı")
	})

	t.Run("limit ve offset parametreleri doğru çalışmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		createOrderHelper(t, ctx, orderRepo, userID, "status1")
		createOrderHelper(t, ctx, orderRepo, userID, "status2")

		// act
		page1, err1 := orderRepo.GetAll(ctx, "1", "0")
		page2, err2 := orderRepo.GetAll(ctx, "1", "1")

		// assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Len(t, page1, 1)
		assert.Len(t, page2, 1)
		assert.NotEqual(t, page1[0].OrderId, page2[0].OrderId, "Farklı sayfalar farklı siparişleri getirmeli")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		orders, err := orderRepo.GetAll(ctx, "10", "0")

		// assert
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Contains(t, err.Error(), "sipariş sorgu hatası")
	})
}

func TestDeleteOrder(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("siparişi başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order := createOrderHelper(t, ctx, orderRepo, userID, "pending")
		orderUUID, err := uuid.Parse(order.Order_id)
		require.NoError(t, err)

		// act
		err = orderRepo.Delete(ctx, orderUUID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: satır silinmiş olmalı
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT count(*) FROM orders WHERE order_id = $1", orderUUID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("olmayan sipariş silinmek istendiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := orderRepo.Delete(ctx, uuid.New())

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := orderRepo.Delete(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sipariş silme hatası")
	})
}

func TestGetOrdersByUserId(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("kullanıcıya ait siparişleri ve toplam ürün adedini getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order1 := createOrderHelper(t, ctx, orderRepo, userID, "processing")
		order2 := createOrderHelper(t, ctx, orderRepo, userID, "delivered")

		// order1 için 2 farklı ürün ekle (productorders)
		productID1 := createTestProduct(t, ctx, testDbInstance)
		productID2 := createTestProduct(t, ctx, testDbInstance)
		queryPO := `INSERT INTO productorders (order_id, product_id, price_at_purchase, quantity) VALUES ($1, $2, 1000, 3), ($1, $3, 1000, 2)`
		_, err := testDbInstance.Exec(ctx, queryPO, order1.Order_id, productID1, productID2)
		require.NoError(t, err)

		// act
		orders, err := orderRepo.GetOrdersByUserId(ctx, userID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, orders, 2)

		for _, o := range orders {
			if o.OrderId == order1.OrderCode {
				assert.Equal(t, 5, o.TotalQuantity, "Order1 toplam ürün miktarı 3 + 2 = 5 olmalı")
				assert.Equal(t, "processing", o.Status)
			}
			if o.OrderId == order2.OrderCode {
				assert.Equal(t, 0, o.TotalQuantity, "Order2 için ürün eklenmediğinden 0 olmalı")
				assert.Equal(t, "delivered", o.Status)
			}
		}
	})

	t.Run("siparişi olmayan kullanıcı için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomUserID := uuid.NewString()

		// act
		orders, err := orderRepo.GetOrdersByUserId(ctx, randomUserID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		orders, err := orderRepo.GetOrdersByUserId(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Nil(t, orders)
	})
}

func TestUpdateOrderState(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("order_code ile sipariş durumunu başarıyla güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order := createOrderHelper(t, ctx, orderRepo, userID, "pending")

		// act
		err := orderRepo.UpdateOrderState(ctx, "completed", order.OrderCode)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var updatedStatus string
		query := `SELECT status FROM orders WHERE order_code = $1`
		err = testDbInstance.QueryRow(ctx, query, order.OrderCode).Scan(&updatedStatus)
		assert.NoError(t, err)
		assert.Equal(t, "completed", updatedStatus)
	})

	t.Run("olmayan order_code ile çağrıldığında hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := orderRepo.UpdateOrderState(ctx, "completed", "NOTEXIST")

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := orderRepo.UpdateOrderState(ctx, "completed", "ORD12345")

		// assert
		assert.Error(t, err)
	})
}

func TestGetOrderCount(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("toplam sipariş sayısını doğru dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		var countBefore int
		err := testDbInstance.QueryRow(ctx, "SELECT count(*) FROM orders").Scan(&countBefore)
		require.NoError(t, err)

		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		createOrderHelper(t, ctx, orderRepo, userID, "pending")

		// act
		count, err := orderRepo.GetOrderCount(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, countBefore+1, count)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := orderRepo.GetOrderCount(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestGetOrderCountByStatus(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("belirli durumdaki siparişlerin sayısını doğru dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		statusKey := fmt.Sprintf("st_%s", uuid.NewString()[:5])
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)

		createOrderHelper(t, ctx, orderRepo, userID, statusKey)
		createOrderHelper(t, ctx, orderRepo, userID, statusKey)

		// act
		count, err := orderRepo.GetOrderCountByStatus(ctx, statusKey)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("eşleşmeyen durum için 0 dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		count, err := orderRepo.GetOrderCountByStatus(ctx, "non_existent_status_123")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := orderRepo.GetOrderCountByStatus(ctx, "pending")

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestUpdateMultiState(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("birden fazla siparişin durumunu toplu olarak güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		ord1 := createOrderHelper(t, ctx, orderRepo, userID, "initial_status")
		ord2 := createOrderHelper(t, ctx, orderRepo, userID, "initial_status")
		ord3 := createOrderHelper(t, ctx, orderRepo, userID, "initial_status")

		// act: ord1 ve ord2'yi güncelle, ord3 dokunulmasın
		err := orderRepo.UpdateMultiState(ctx, []string{ord1.OrderCode, ord2.OrderCode}, "updated_multi")

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var status1, status2, status3 string
		_ = testDbInstance.QueryRow(ctx, "SELECT status FROM orders WHERE order_code = $1", ord1.OrderCode).Scan(&status1)
		_ = testDbInstance.QueryRow(ctx, "SELECT status FROM orders WHERE order_code = $1", ord2.OrderCode).Scan(&status2)
		_ = testDbInstance.QueryRow(ctx, "SELECT status FROM orders WHERE order_code = $1", ord3.OrderCode).Scan(&status3)

		assert.Equal(t, "updated_multi", status1)
		assert.Equal(t, "updated_multi", status2)
		assert.Equal(t, "initial_status", status3, "ord3 durumu değişmemeli")
	})

	t.Run("boş liste verildiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := orderRepo.UpdateMultiState(ctx, []string{}, "new_state")

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := orderRepo.UpdateMultiState(ctx, []string{"ORD123"}, "new_state")

		// assert
		assert.Error(t, err)
	})
}

func TestGetOrdersByStatus(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("duruma göre filtrelenmiş siparişleri sayfalanmış olarak getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		statusTarget := fmt.Sprintf("st_filt_%s", uuid.NewString()[:5])
		userID, userName, userSurname := createOrderTestUser(t, ctx, testDbInstance)

		ord1 := createOrderHelper(t, ctx, orderRepo, userID, statusTarget)
		ord2 := createOrderHelper(t, ctx, orderRepo, userID, statusTarget)
		_ = createOrderHelper(t, ctx, orderRepo, userID, "other_status")

		// act
		orders, err := orderRepo.GetOrdersByStatus(ctx, "10", "0", statusTarget)

		// assert
		assert.NoError(t, err)
		assert.Len(t, orders, 2)

		orderCodes := []string{orders[0].OrderId, orders[1].OrderId}
		assert.Contains(t, orderCodes, ord1.OrderCode)
		assert.Contains(t, orderCodes, ord2.OrderCode)
		assert.Equal(t, userName, orders[0].Name)
		assert.Equal(t, userSurname, orders[0].Surname)
	})

	t.Run("eşleşmeyen durum için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		orders, err := orderRepo.GetOrdersByStatus(ctx, "10", "0", "hic_olmayan_durum")

		// assert
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		orders, err := orderRepo.GetOrdersByStatus(ctx, "10", "0", "pending")

		// assert
		assert.Error(t, err)
		assert.Nil(t, orders)
	})
}

func TestGetOrderByOrderCode(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("order_code ile sipariş adres detaylarını getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order := createOrderHelper(t, ctx, orderRepo, userID, "pending")

		// act
		res, err := orderRepo.GetOrderByOrderCode(ctx, order.OrderCode)

		// assert
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, order.Order_id, res.OrderId)
		assert.Equal(t, order.ShippingAddress, res.ShippingAddress)
		assert.Equal(t, order.BillingAddress, res.BillingAddress)
	})

	t.Run("olmayan order_code sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentCode := "NOEXIST1"

		// act
		res, err := orderRepo.GetOrderByOrderCode(ctx, nonExistentCode)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		res, err := orderRepo.GetOrderByOrderCode(ctx, "CANCEL12")

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetOrderCountForWaiting(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("status='pending' olan sipariş sayısını doğru getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		var pendingBefore int
		err := testDbInstance.QueryRow(ctx, "SELECT count(*) FROM orders WHERE status='pending'").Scan(&pendingBefore)
		require.NoError(t, err)

		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		createOrderHelper(t, ctx, orderRepo, userID, "pending")
		createOrderHelper(t, ctx, orderRepo, userID, "pending")
		createOrderHelper(t, ctx, orderRepo, userID, "completed") // pending değil

		// act
		count, err := orderRepo.GetOrderCountForWaiting(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, pendingBefore+2, count)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := orderRepo.GetOrderCountForWaiting(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestSearchOrder(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("order_code ile sipariş aramasını başarıyla gerçekleştirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order := createOrderHelper(t, ctx, orderRepo, userID, "pending")

		// act
		results, err := orderRepo.SearchOrder(ctx, "10", order.OrderCode)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)

		found := false
		for _, r := range results {
			if r.OrderId == order.OrderCode {
				found = true
				assert.Equal(t, "pending", r.Status)
			}
		}
		assert.True(t, found, "Arama sonuçlarında orderCode bulunmalı")
	})

	t.Run("eşleşmeyen arama için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := orderRepo.SearchOrder(ctx, "10", "NOPE9999")

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := orderRepo.SearchOrder(ctx, "10", "ANY12345")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestSearchOrderCount(t *testing.T) {
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("arama sorgusuna uyan sipariş sayısını doğru dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID, _, _ := createOrderTestUser(t, ctx, testDbInstance)
		order := createOrderHelper(t, ctx, orderRepo, userID, "pending")

		// act
		count, err := orderRepo.SearchOrderCount(ctx, order.OrderCode)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("eşleşmeyen arama için 0 dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		count, err := orderRepo.SearchOrderCount(ctx, "ZERO9999")

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := orderRepo.SearchOrderCount(ctx, "ANY12345")

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}
