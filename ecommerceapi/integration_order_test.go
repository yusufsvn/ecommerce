package main

import (
	"context"
	constants "ecommerceapi/Constants"
	models "ecommerceapi/Database/Models"
	structs "ecommerceapi/Structs"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Order_CreateAndLifecycle(t *testing.T) {
	app := setupIntegrationApp(t)
	ctx := context.Background()

	// 1. Kullanıcı ve Admin oluştur
	userEmail := fmt.Sprintf("order-user-%s@example.com", uuid.NewString()[:8])
	adminEmail := fmt.Sprintf("order-admin-%s@example.com", uuid.NewString()[:8])
	userID, userAccessCookie, _ := app.createTestUser(t, userEmail, "Pass123!", "user")
	_, adminAccessCookie, _ := app.createTestUser(t, adminEmail, "Pass123!", "admin")

	// 2. Kategori ve Ürün oluştur
	cat, err := app.CategoryStore.Create(ctx, fmt.Sprintf("OrderCat-%s", uuid.NewString()[:6]))
	assert.NoError(t, err)

	prodID := uuid.NewString()
	now := time.Now()
	pCode := uuid.NewString()[:8]
	prod := models.Product{
		ProductID:           prodID,
		ProductName:         fmt.Sprintf("Sipariş Test Ürünü %s", uuid.NewString()[:4]),
		Stock:               25,
		Price:               2500,
		ThumbnailImage:      uuid.NewString(),
		Description:         "Sipariş entegrasyon ürünü",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		ItemType:            "PHYSICAL",
		ProductCode:         pCode,
		ProductPublishState: "published",
		TechnicalDetails:    []byte(`{}`),
	}
	err = app.ProductStore.Create(ctx, &prod)
	assert.NoError(t, err)

	err = app.PPCStore.AddMany(ctx, []string{cat.Category_id}, prodID)
	assert.NoError(t, err)

	// 3. Ürünü sepete ekle
	addPayload := structs.AddToCartReq{
		ProductID: prodID,
		Quantity:  2,
	}
	resp := app.doRequest(t, http.MethodPost, "/api/user/addtocart", addPayload, userAccessCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var createdOrderID string

	t.Run("kullanıcı başarıyla sipariş oluşturabilmeli ve sepet temizlenmeli", func(t *testing.T) {
		paymentReq := structs.CreatePaymentDto{
			Tckn:        "12345678901",
			PhoneNumber: "5551234567",
			CardData: structs.Card{
				CardName:    "Test User",
				CardNumber:  "5528790000000001",
				ExpireMonth: "12",
				ExpireYear:  "2030",
				Cvc:         "123",
			},
			BillingAddress: structs.Address{
				AddressName: "Ev",
				AddressStr:  "Barbaros Bulvarı No: 12",
				City:        "İstanbul",
				District:    "Beşiktaş",
				Country:     "Türkiye",
				ZipCode:     "34349",
			},
			ShippingAddress: structs.Address{
				AddressName: "Ev",
				AddressStr:  "Barbaros Bulvarı No: 12",
				City:        "İstanbul",
				District:    "Beşiktaş",
				Country:     "Türkiye",
				ZipCode:     "34349",
			},
			SecurePay:   false,
			SaveAddress: false,
		}

		orderResp := app.doRequest(t, http.MethodPost, "/api/user/createorder", paymentReq, userAccessCookie)
		assert.Equal(t, http.StatusOK, orderResp.StatusCode)

		// Sepetin temizlendiğini doğrula
		cartResp := app.doRequest(t, http.MethodGet, "/api/user/getcart", nil, userAccessCookie)
		assert.Equal(t, http.StatusOK, cartResp.StatusCode)
		var cartRes *structs.CartItems
		_ = json.NewDecoder(cartResp.Body).Decode(&cartRes)
		if cartRes != nil {
			assert.Empty(t, cartRes.Products, "Sipariş verildikten sonra sepet boş olmalı")
		}
	})

	t.Run("kullanıcı sipariş geçmişini listeleyebilmeli", func(t *testing.T) {
		historyResp := app.doRequest(t, http.MethodGet, "/api/user/getorderbyuserid", nil, userAccessCookie)
		assert.Equal(t, http.StatusOK, historyResp.StatusCode)

		var orders []structs.GetOrder
		err := json.NewDecoder(historyResp.Body).Decode(&orders)
		assert.NoError(t, err)
		assert.NotEmpty(t, orders, "Kullanıcının en az 1 siparişi olmalı")
		if len(orders) > 0 {
			if orders[0].OrderCode != "" {
				createdOrderID = orders[0].OrderCode
			} else {
				createdOrderID = orders[0].OrderId
			}
		}
		assert.NotEmpty(t, createdOrderID)
		if len(orders) > 0 {
			assert.Equal(t, constants.OrderStatePending, orders[0].Status)
		}
	})

	t.Run("kullanıcı sipariş detayını görebilmeli", func(t *testing.T) {
		assert.NotEmpty(t, createdOrderID)
		detailResp := app.doRequest(t, http.MethodGet, fmt.Sprintf("/api/user/getorderdetails?orderId=%s", createdOrderID), nil, userAccessCookie)
		assert.Equal(t, http.StatusOK, detailResp.StatusCode)

		var details structs.OrderResponse
		err := json.NewDecoder(detailResp.Body).Decode(&details)
		assert.NoError(t, err)
		assert.NotEmpty(t, details.Products, "Sipariş detayında ürünler bulunmalı")
		if len(details.Products) > 0 {
			assert.Equal(t, prodID, details.Products[0].ProductID)
		}
	})

	t.Run("admin bekleyen sipariş sayısını görebilmeli", func(t *testing.T) {
		countResp := app.doRequest(t, http.MethodGet, "/api/admin/getordercountforwait", nil, adminAccessCookie)
		assert.Equal(t, http.StatusOK, countResp.StatusCode)

		var count int
		err := json.NewDecoder(countResp.Body).Decode(&count)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, 1)
	})

	t.Run("admin sipariş durumunu güncelleyebilmeli", func(t *testing.T) {
		assert.NotEmpty(t, createdOrderID)
		updateStatePayload := structs.UpdateOrderStateReq{
			Order_id: createdOrderID,
			NewState: constants.OrderStateShipped,
		}

		updateResp := app.doRequest(t, http.MethodPut, "/api/admin/updateorderstate", updateStatePayload, adminAccessCookie)
		assert.Equal(t, http.StatusOK, updateResp.StatusCode)

		// Durumun güncellendiğini kullanıcı sipariş geçmişinden doğrula
		historyResp := app.doRequest(t, http.MethodGet, "/api/user/getorderbyuserid", nil, userAccessCookie)
		assert.Equal(t, http.StatusOK, historyResp.StatusCode)
		var orders []structs.GetOrder
		_ = json.NewDecoder(historyResp.Body).Decode(&orders)
		assert.NotEmpty(t, orders)
		assert.Equal(t, constants.OrderStateShipped, orders[0].Status)
	})

	t.Run("admin tüm siparişleri listeleyebilmeli", func(t *testing.T) {
		allResp := app.doRequest(t, http.MethodGet, "/api/admin/getallorders?limit=10&offset=0", nil, adminAccessCookie)
		assert.Equal(t, http.StatusOK, allResp.StatusCode)

		var orderResult structs.OrderResult
		err := json.NewDecoder(allResp.Body).Decode(&orderResult)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, orderResult.Count, 1)
	})

	_ = userID
}

func TestIntegration_Order_ValidationErrors(t *testing.T) {
	app := setupIntegrationApp(t)

	userEmail := fmt.Sprintf("invalid-order-%s@example.com", uuid.NewString()[:8])
	_, userAccessCookie, _ := app.createTestUser(t, userEmail, "Pass123!", "user")

	t.Run("geçersiz telefon numarası ile sipariş oluşturulamamalı", func(t *testing.T) {
		paymentReq := structs.CreatePaymentDto{
			PhoneNumber: "02121234567", // 5 ile başlamıyor!
			Tckn:        "12345678901",
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/createorder", paymentReq, userAccessCookie)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("boş sepet ile sipariş oluşturulamamalı", func(t *testing.T) {
		paymentReq := structs.CreatePaymentDto{
			PhoneNumber: "5551234567",
			Tckn:        "12345678901",
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/createorder", paymentReq, userAccessCookie)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
