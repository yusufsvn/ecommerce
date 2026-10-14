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

func setupOrderHandlerTest() (
	*OrderHandler,
	*mockOrderRepo,
	*mockProductOrderRepo,
	*mockCartRepo,
	*mockCartProductRepo,
	*mockUserRepo,
	*mockProductRepo,
	*mockProductProductCategoryRepo,
	*mockAddressRepo,
	*mockPaymentService,
) {
	orderRepo := &mockOrderRepo{}
	prodOrderRepo := &mockProductOrderRepo{}
	cartRepo := &mockCartRepo{}
	cartProductRepo := &mockCartProductRepo{}
	userRepo := &mockUserRepo{}
	prodRepo := &mockProductRepo{}
	ppcRepo := &mockProductProductCategoryRepo{}
	addrRepo := &mockAddressRepo{}
	paymentSvc := &mockPaymentService{}

	cartSvc := services.NewCartService(cartRepo, cartProductRepo)
	addrSvc := services.NewAddressService(addrRepo)

	orderSvc := services.NewOrderService(
		cartSvc,
		cartProductRepo,
		userRepo,
		paymentSvc,
		orderRepo,
		prodOrderRepo,
		prodRepo,
		ppcRepo,
		"http://storage.api",
		"http://public.storage",
		"test-bucket",
		*addrSvc,
	)

	handler := NewOrderHandler(*orderSvc)

	return handler, orderRepo, prodOrderRepo, cartRepo, cartProductRepo, userRepo, prodRepo, ppcRepo, addrRepo, paymentSvc
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	t.Run("geçerli istek ve RemoteAddr ile başarıyla sipariş oluşturmalı", func(t *testing.T) {
		handler, orderRepo, prodOrderRepo, cartRepo, cartProductRepo, userRepo, _, ppcRepo, _, paymentSvc := setupOrderHandlerTest()

		cartRepo.getCartByUserIDFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-1", nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: "p-1", ProductName: "Urun", Price: 100, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*st.GetAllCategory, error) {
			return []*st.GetAllCategory{
				{ProductId: "p-1", Categories: []st.CategoryDbJson{{CategoryName: "Genel"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, user_id string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Yusuf", Surname: "A", Email: "y@t.com", Phone: "5301234567"}, nil
		}
		paymentSvc.createPaymentNonTDSFn = func(req st.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			return nil
		}
		prodOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			return nil
		}
		cartProductRepo.removeAllFn = func(ctx context.Context, cart_id string) error {
			return nil
		}

		body := st.CreatePaymentDto{
			PhoneNumber: "5301234567",
			ShippingAddress: st.Address{
				AddressStr: "Kadikoy",
			},
			BillingAddress: st.Address{
				AddressStr: "Kadikoy",
			},
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/order", bytes.NewReader(bodyBytes))
		req.RemoteAddr = "127.0.0.1:8080"
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.CreateOrder(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("geçersiz RemoteAddr formatında 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()
		body := st.CreatePaymentDto{
			PhoneNumber: "5301234567",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/order", bytes.NewReader(bodyBytes))
		req.RemoteAddr = "invalid-no-port"
		rec := httptest.NewRecorder()

		handler.CreateOrder(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Error splitting host port")
	})

	t.Run("telefon numarası geçersiz olduğunda veya servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()
		body := st.CreatePaymentDto{
			PhoneNumber: "05301234567", // 53 ile başlamıyor
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/order", bytes.NewReader(bodyBytes))
		req.RemoteAddr = "127.0.0.1:8080"
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.CreateOrder(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sipariş oluşturulurken hata oluştu")
	})
}

func TestOrderHandler_GetOrdersByUserId(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("başarıyla kullanıcının siparişlerini dönmeli", func(t *testing.T) {
		orderRepo.getOrdersByUserIdFn = func(ctx context.Context, user_id string) ([]st.GetOrder, error) {
			assert.Equal(t, "user-123", user_id)
			return []st.GetOrder{
				{OrderId: "order-1", TotalPrice: 250, Status: "pending"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetOrdersByUserId(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp []st.GetOrder
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "order-1", resp[0].OrderId)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.getOrdersByUserIdFn = func(ctx context.Context, user_id string) ([]st.GetOrder, error) {
			return nil, errors.New("db error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetOrdersByUserId(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "siparişler getirilemedi")
	})
}

func TestOrderHandler_GetOrderDetails(t *testing.T) {
	handler, orderRepo, prodOrderRepo, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("geçerli orderId ile sipariş detaylarını başarıyla dönmeli", func(t *testing.T) {
		orderRepo.getOrderByOrderCodeFn = func(ctx context.Context, order_code string) (*st.GetOrderAddress, error) {
			return &st.GetOrderAddress{OrderId: "o-1", ShippingAddress: "Adres 1"}, nil
		}
		prodOrderRepo.getByOrderIDFn = func(ctx context.Context, orderID string) ([]st.OrderFullDetail, error) {
			return []st.OrderFullDetail{
				{ProductID: "p-1", ProductName: "Urun 1", Quantity: 1},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/order/detail?orderId=ORD123", nil)
		rec := httptest.NewRecorder()

		handler.GetOrderDetails(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.OrderResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Adres 1", resp.ShippingAddress)
		assert.Len(t, resp.Products, 1)
	})

	t.Run("orderId eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/order/detail", nil)
		rec := httptest.NewRecorder()

		handler.GetOrderDetails(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "order_id eksik")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.getOrderByOrderCodeFn = func(ctx context.Context, order_code string) (*st.GetOrderAddress, error) {
			return nil, errors.New("order not found")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/order/detail?orderId=ORD123", nil)
		rec := httptest.NewRecorder()

		handler.GetOrderDetails(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "beklenmeyen hatayla karşılaşıldı")
	})
}

func TestOrderHandler_GetAllOrders(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("başarıyla tüm siparişleri dönmeli", func(t *testing.T) {
		orderRepo.getAllFn = func(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {
			return []st.OrderForAdmin{
				{OrderId: "ord-1", Name: "Ali", TotalAmount: 500},
			}, nil
		}
		orderRepo.getOrderCountFn = func(ctx context.Context) (int, error) {
			return 1, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAllOrders(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.OrderResult
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
		assert.Len(t, resp.Data, 1)
	})

	t.Run("limit veya offset eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders?limit=10", nil)
		rec := httptest.NewRecorder()

		handler.GetAllOrders(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "query stringsler boş olamaz")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("db error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetAllOrders(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmeyen Hatayla karşılaşıldı")
	})
}

func TestOrderHandler_UpdateOrderState(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("başarıyla sipariş durumunu güncellemeli", func(t *testing.T) {
		updated := false
		orderRepo.updateOrderStateFn = func(ctx context.Context, newState string, order_id string) error {
			assert.Equal(t, "shipped", newState)
			assert.Equal(t, "ord-1", order_id)
			updated = true
			return nil
		}

		body := st.UpdateOrderStateReq{
			Order_id: "ord-1",
			NewState: "shipped",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/order/state", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdateOrderState(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, updated)
	})

	t.Run("geçersiz durum veya servis hatasında 500 dönmeli", func(t *testing.T) {
		body := st.UpdateOrderStateReq{
			Order_id: "ord-1",
			NewState: "invalid-state",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/order/state", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdateOrderState(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmedik Hatayla karşılaşıldı")
	})
}

func TestOrderHandler_UpdateMultiState(t *testing.T) {
	t.Run("başarıyla birden fazla siparişin durumunu güncellemeli", func(t *testing.T) {
		handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()
		updated := false
		orderRepo.updateMultiStateFn = func(ctx context.Context, idlist []string, new_state string) error {
			assert.Equal(t, []string{"o-1", "o-2"}, idlist)
			assert.Equal(t, "delivered", new_state)
			updated = true
			return nil
		}

		body := st.UpdateMultiOrderReq{
			OrderIds: []string{"o-1", "o-2"},
			NewState: "delivered",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/order/multi-state", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdateMultiState(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, updated)
	})

	t.Run("servis/db hatasında 500 dönmeli", func(t *testing.T) {
		handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()
		orderRepo.updateMultiStateFn = func(ctx context.Context, idlist []string, new_state string) error {
			return errors.New("db error")
		}
		body := st.UpdateMultiOrderReq{
			OrderIds: []string{"o-1"},
			NewState: "delivered",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/order/multi-state", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdateMultiState(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün State i değiştirilemedi")
	})
}

func TestOrderHandler_GetOrdersByStatus(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("başarıyla duruma göre siparişleri dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountByStatusFn = func(ctx context.Context, status string) (int, error) {
			return 1, nil
		}
		orderRepo.getOrdersByStatusFn = func(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {
			return []st.OrderForAdmin{
				{OrderId: "o-1", Status: "pending"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/status?limit=10&offset=0&status=pending", nil)
		rec := httptest.NewRecorder()

		handler.GetOrdersByStatus(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.OrderResult
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
	})

	t.Run("parametrelerden biri eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/status?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetOrdersByStatus(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "query stringsler boş olamaz")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountByStatusFn = func(ctx context.Context, status string) (int, error) {
			return 0, errors.New("db error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/status?limit=10&offset=0&status=pending", nil)
		rec := httptest.NewRecorder()

		handler.GetOrdersByStatus(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmeyen Hatayla karşılaşıldı")
	})
}

func TestOrderHandler_GetOrderCountForWaiting(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("başarıyla bekleyen sipariş sayısını dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountForWaitingFn = func(ctx context.Context) (int, error) {
			return 7, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/waiting-count", nil)
		rec := httptest.NewRecorder()

		handler.GetOrderCountForWaiting(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var count int
		err := json.Unmarshal(rec.Body.Bytes(), &count)
		assert.NoError(t, err)
		assert.Equal(t, 7, count)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountForWaitingFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("count error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/waiting-count", nil)
		rec := httptest.NewRecorder()

		handler.GetOrderCountForWaiting(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sipariş sayısı Getirilirken hata oluştu")
	})
}

func TestOrderHandler_SearchOrder(t *testing.T) {
	handler, orderRepo, _, _, _, _, _, _, _, _ := setupOrderHandlerTest()

	t.Run("arama parametreleri eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/search?limit=10", nil)
		rec := httptest.NewRecorder()

		handler.SearchOrder(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "query stringsler boş olamaz")
	})

	t.Run("sonuç varsa başarıyla listeyi dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, search_query string) (int, error) {
			return 1, nil
		}
		orderRepo.searchOrderFn = func(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error) {
			return []st.OrderForAdmin{
				{OrderId: "ord-1", Name: "Ali"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/search?limit=10&search_query=Ali", nil)
		rec := httptest.NewRecorder()

		handler.SearchOrder(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp st.OrderResult
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
		assert.Len(t, resp.Data, 1)
	})

	t.Run("arama sonucu 0 ise (ZeroListResult) boş liste dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, search_query string) (int, error) {
			return 0, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/search?limit=10&search_query=Bilinmeyen", nil)
		rec := httptest.NewRecorder()

		handler.SearchOrder(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp st.OrderResult
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Count)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, search_query string) (int, error) {
			return 0, errors.New("search error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/orders/search?limit=10&search_query=Hata", nil)
		rec := httptest.NewRecorder()

		handler.SearchOrder(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmedik Hatayla karşılaşıldı")
	})
}
