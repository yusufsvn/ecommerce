package main

import (
	"context"
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

func TestIntegration_Cart_Lifecycle(t *testing.T) {
	app := setupIntegrationApp(t)
	ctx := context.Background()

	// 1. Kullanıcı ve sepet oluştur
	userEmail := fmt.Sprintf("cart-user-%s@example.com", uuid.NewString()[:8])
	_, accessCookie, _ := app.createTestUser(t, userEmail, "Pass123!", "user")

	// 2. Ürün oluştur
	prodID := uuid.NewString()
	now := time.Now()
	pCode := uuid.NewString()[:8]
	prod := models.Product{
		ProductID:           prodID,
		ProductName:         fmt.Sprintf("Sepet Test Ürünü %s", uuid.NewString()[:4]),
		Stock:               20,
		Price:               1200,
		ThumbnailImage:      uuid.NewString(),
		Description:         "Sepet entegrasyon testi ürünü",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		ItemType:            "PHYSICAL",
		ProductCode:         pCode,
		ProductPublishState: "published",
		TechnicalDetails:    []byte(`{}`),
	}
	err := app.ProductStore.Create(ctx, &prod)
	assert.NoError(t, err)

	t.Run("sepete ürün başarıyla eklenmeli", func(t *testing.T) {
		addPayload := structs.AddToCartReq{
			ProductID: prodID,
			Quantity:  2,
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/addtocart", addPayload, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("sepet içeriği başarıyla çekilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/user/getcart", nil, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cartRes structs.CartItems
		err := json.NewDecoder(resp.Body).Decode(&cartRes)
		assert.NoError(t, err)
		assert.NotEmpty(t, cartRes.Products)

		var found bool
		for _, item := range cartRes.Products {
			if item.Product_id == prodID {
				found = true
				assert.Equal(t, 2, item.Quantity)
				break
			}
		}
		assert.True(t, found, "Eklenen ürün sepette bulunmalı")
	})

	t.Run("sepetteki ürün adedi güncellenebilmeli", func(t *testing.T) {
		updatePayload := structs.UpdateCartReq{
			Product_id: prodID,
			Quantity:   5,
		}

		resp := app.doRequest(t, http.MethodPut, "/api/user/updatecartitem", updatePayload, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Güncellemenin yansıdığını doğrula
		getResp := app.doRequest(t, http.MethodGet, "/api/user/getcart", nil, accessCookie)
		var cartRes structs.CartItems
		_ = json.NewDecoder(getResp.Body).Decode(&cartRes)

		for _, item := range cartRes.Products {
			if item.Product_id == prodID {
				assert.Equal(t, 5, item.Quantity)
				break
			}
		}
	})

	t.Run("ürün sepetten çıkarılabilmeli", func(t *testing.T) {
		removePayload := structs.RemoveFromCartReq{
			Product_id: prodID,
		}

		resp := app.doRequest(t, http.MethodDelete, "/api/user/removefromcart", removePayload, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Sepetin boşaldığını doğrula (ZeroListResult nil döner)
		getResp := app.doRequest(t, http.MethodGet, "/api/user/getcart", nil, accessCookie)
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var cartRes *structs.CartItems
		_ = json.NewDecoder(getResp.Body).Decode(&cartRes)
		if cartRes != nil {
			assert.Empty(t, cartRes.Products)
		}
	})
}
