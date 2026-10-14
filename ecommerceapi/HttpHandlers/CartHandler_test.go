package httphandlers

import (
	"bytes"
	"context"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCartHandlerTest() (*CartHandler, *mockCartRepo, *mockCartProductRepo) {
	cartRepo := &mockCartRepo{}
	cartProductRepo := &mockCartProductRepo{}
	svc := services.NewCartService(cartRepo, cartProductRepo)
	handler := NewCartHandler(*svc)
	return handler, cartRepo, cartProductRepo
}

func TestCartHandler_AddToCart(t *testing.T) {
	handler, cartRepo, cartProductRepo := setupCartHandlerTest()

	t.Run("başarıyla sepete ürün eklemeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.addFn = func(ctx context.Context, cartProduct *st.CartProductDTO) error {
			assert.Equal(t, "cart-1", cartProduct.Cart_id)
			assert.Equal(t, "prod-1", cartProduct.Product_id)
			assert.Equal(t, 2, cartProduct.Quantity)
			return nil
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
			"quantity":   2,
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.AddToCart(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "", errors.New("cart error")
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
			"quantity":   2,
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.AddToCart(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sepete ürün eklenirken hata oluştu")
	})
}

func TestCartHandler_GetCart(t *testing.T) {
	handler, cartRepo, _ := setupCartHandlerTest()

	t.Run("sepet dolu olduğunda ürünleri başarıyla dönmeli", func(t *testing.T) {
		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, user_id string) ([]st.CartProductDTO, error) {
			assert.Equal(t, "user-123", user_id)
			return []st.CartProductDTO{
				{
					Cart_id:    "cart-1",
					Product_id: "prod-1",
					Quantity:   2,
				},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/cart", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetCart(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.CartItems
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "cart-1", resp.Cart_id)
		assert.Len(t, resp.Products, 1)
		assert.Equal(t, "prod-1", resp.Products[0].Product_id)
	})

	t.Run("sepet boş olduğunda (ZeroListResult) 200 dönmeli", func(t *testing.T) {
		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, user_id string) ([]st.CartProductDTO, error) {
			return []st.CartProductDTO{}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/cart", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetCart(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("servis/db hata verdiğinde 500 dönmeli", func(t *testing.T) {
		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, user_id string) ([]st.CartProductDTO, error) {
			return nil, errors.New("db failed")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/cart", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetCart(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sepette ürün bulunamadı")
	})
}

func TestCartHandler_RemoveFromCart(t *testing.T) {
	handler, cartRepo, cartProductRepo := setupCartHandlerTest()

	t.Run("başarıyla sepetten ürün silmeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.removeFn = func(ctx context.Context, cart_id string, product_id string) error {
			assert.Equal(t, "cart-1", cart_id)
			assert.Equal(t, "prod-1", product_id)
			return nil
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodDelete, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.RemoveFromCart(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.removeFn = func(ctx context.Context, cart_id string, product_id string) error {
			return errors.New("remove error")
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodDelete, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.RemoveFromCart(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sepetten ürün çıkarılırken hata oluştu")
	})
}

func TestCartHandler_UpdateCart(t *testing.T) {
	handler, cartRepo, cartProductRepo := setupCartHandlerTest()

	t.Run("başarıyla sepet miktarını güncellemeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.updateQuantityFn = func(ctx context.Context, data *st.CartProductDTO) error {
			assert.Equal(t, "cart-1", data.Cart_id)
			assert.Equal(t, "prod-1", data.Product_id)
			assert.Equal(t, 5, data.Quantity)
			return nil
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
			"quantity":   5,
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.UpdateCart(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.updateQuantityFn = func(ctx context.Context, data *st.CartProductDTO) error {
			return errors.New("update error")
		}

		body := map[string]interface{}{
			"product_id": "prod-1",
			"quantity":   5,
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/user/cart", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.UpdateCart(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sepet güncellenirken hata oluştu")
	})
}
