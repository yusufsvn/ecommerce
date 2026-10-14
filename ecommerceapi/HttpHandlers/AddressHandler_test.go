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

func setupAddressHandlerTest() (*AddressHandler, *mockAddressRepo) {
	repo := &mockAddressRepo{}
	svc := services.NewAddressService(repo)
	handler := NewAddressHandler(*svc)
	return handler, repo
}

func TestAddressHandler_GetAddress(t *testing.T) {
	handler, repo := setupAddressHandlerTest()

	t.Run("geçerli kullanıcı context ile adresleri başarıyla dönmeli", func(t *testing.T) {
		repo.getByUserIdFn = func(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error) {
			assert.Equal(t, "user-123", user_id)
			return &[]st.GetAddressDTO{
				{AddressId: "addr-1", City: "Istanbul", Country: "Turkey"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/address", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetAddress(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.GetAddressList
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Addresses)
		assert.Len(t, *resp.Addresses, 1)
		assert.Equal(t, "Istanbul", (*resp.Addresses)[0].City)
	})

	t.Run("context içerisinde user_id yoksa 500 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/address", nil)
		rec := httptest.NewRecorder()

		handler.GetAddress(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "adres verisi bulunamadı")
	})

	t.Run("servis/db hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.getByUserIdFn = func(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error) {
			return nil, errors.New("db error")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/address", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetAddress(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "adres verisi bulunamadı")
	})
}

func TestAddressHandler_UpdateAddress(t *testing.T) {
	handler, repo := setupAddressHandlerTest()

	t.Run("başarıyla adres güncellemeli", func(t *testing.T) {
		repo.updateFn = func(ctx context.Context, address st.Address, user_id string) error {
			assert.Equal(t, "user-123", user_id)
			assert.Equal(t, "Ankara", address.City)
			return nil
		}

		body := map[string]interface{}{
			"address_id":   "addr-1",
			"city":         "Ankara",
			"country":      "Turkey",
			"district":     "Cankaya",
			"name":         "Ev",
			"post_code":    "06100",
			"full_address": "Kizilay cad.",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/user/address", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.UpdateAddress(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.updateFn = func(ctx context.Context, address st.Address, user_id string) error {
			return errors.New("update failed")
		}

		body := map[string]interface{}{
			"address_id":   "addr-1",
			"city":         "Ankara",
			"country":      "Turkey",
			"district":     "Cankaya",
			"name":         "Ev",
			"post_code":    "06100",
			"full_address": "Kizilay cad.",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/user/address", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.UpdateAddress(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "adres güncellenemedi")
	})
}

func TestAddressHandler_CreateAddress(t *testing.T) {
	handler, repo := setupAddressHandlerTest()

	t.Run("başarıyla adres oluşturmalı", func(t *testing.T) {
		repo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			assert.Equal(t, "user-123", address.UserId)
			assert.Equal(t, "Izmir", address.City)
			return nil
		}

		body := map[string]interface{}{
			"city":         "Izmir",
			"country":      "Turkey",
			"district":     "Konak",
			"name":         "Is",
			"post_code":    "35000",
			"full_address": "Alsancak",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/address", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.CreateAddress(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			return errors.New("insert failed")
		}

		body := map[string]interface{}{
			"city":         "Izmir",
			"country":      "Turkey",
			"district":     "Konak",
			"name":         "Is",
			"post_code":    "35000",
			"full_address": "Alsancak",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/address", bytes.NewReader(bodyBytes))
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.CreateAddress(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Adres Eklenirken Hatayla Karşılaşıldı")
	})
}

func TestAddressHandler_DeleteAddress(t *testing.T) {
	handler, repo := setupAddressHandlerTest()

	t.Run("geçerli addressId ile başarıyla silmeli", func(t *testing.T) {
		deleted := false
		repo.deleteFn = func(ctx context.Context, addressId string) error {
			assert.Equal(t, "addr-123", addressId)
			deleted = true
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/user/address?addressId=addr-123", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.DeleteAddress(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, deleted)
	})

	t.Run("addressId query parametresi eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/address", nil)
		rec := httptest.NewRecorder()

		handler.DeleteAddress(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Adres silinemedi")
	})

	t.Run("servis silme sırasında hata verirse 500 dönmeli", func(t *testing.T) {
		repo.deleteFn = func(ctx context.Context, addressId string) error {
			return errors.New("delete failed")
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/user/address?addressId=addr-123", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.DeleteAddress(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmedik bir hatayla karşılaşıldı")
	})
}
