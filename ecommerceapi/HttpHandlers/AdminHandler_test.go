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
	"golang.org/x/crypto/bcrypt"
)

func setupAdminHandlerTest() (*AdminHandler, *mockUserRepo, *mockTokenRepo) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	cartRepo := &mockCartRepo{}
	addrRepo := &mockAddressRepo{}

	tokenSvc := services.NewTokenService(tokenRepo, "test-jwt-secret-key-32-chars-long!")
	userSvc := services.NewUserService(userRepo, *tokenSvc, cartRepo, addrRepo, "google-client-id")
	handler := NewAdminHandler(*userSvc)
	return handler, userRepo, tokenRepo
}

func TestAdminHandler_Login(t *testing.T) {
	handler, userRepo, tokenRepo := setupAdminHandlerTest()

	t.Run("başarılı admin girişinde çerezleri ve 200 dönmeli", func(t *testing.T) {
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.MinCost)
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return &st.UserLoginRes{
				UserID:   "admin-uuid-1",
				Name:     "Admin",
				Surname:  "User",
				Password: string(hashedPw),
				Email:    "admin@test.com",
				Role:     "admin",
			}, nil
		}
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return nil
		}

		body := map[string]interface{}{
			"mail":     "admin@test.com",
			"password": "admin123",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		cookies := rec.Result().Cookies()
		var hasAccess, hasRefresh bool
		for _, c := range cookies {
			if c.Name == "access_token" && c.Value != "" {
				hasAccess = true
			}
			if c.Name == "refresh_token" && c.Value != "" {
				hasRefresh = true
			}
		}
		assert.True(t, hasAccess, "access_token çerezi bulunmalı")
		assert.True(t, hasRefresh, "refresh_token çerezi bulunmalı")
	})

	t.Run("geçersiz şifre veya rol uyuşmazlığında 500 dönmeli", func(t *testing.T) {
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.MinCost)
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			// Normal user admin login'e girmeye çalışıyor
			return &st.UserLoginRes{
				UserID:   "user-uuid-1",
				Password: string(hashedPw),
				Email:    "user@test.com",
				Role:     "user",
			}, nil
		}

		body := map[string]interface{}{
			"mail":     "user@test.com",
			"password": "admin123",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "giriş sırasında hata oluştu")
	})

	t.Run("veritabanı hatasında 500 dönmeli", func(t *testing.T) {
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return nil, errors.New("db connection failure")
		}

		body := map[string]interface{}{
			"mail":     "admin@test.com",
			"password": "admin123",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "giriş sırasında hata oluştu")
	})
}
