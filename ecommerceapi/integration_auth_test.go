package main

import (
	"context"
	structs "ecommerceapi/Structs"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Auth_Register(t *testing.T) {
	app := setupIntegrationApp(t)

	uniqueMail := fmt.Sprintf("test-%s@example.com", uuid.NewString()[:8])
	regPayload := structs.UserRegisterReq{
		Name:           "Ahmet",
		SecondName:     "Yılmaz",
		Mail:           uniqueMail,
		Password:       "Password123!",
		Password_check: "Password123!",
		Phone_number:   "05551234567",
		Address:        "Bağdat Cad. No: 1",
		City:           "İstanbul",
		District:       "Kadıköy",
		Country:        "Türkiye",
		ZipCode:        "34710",
	}

	// 1. Başarılı Kayıt
	resp := app.doRequest(t, http.MethodPost, "/api/user/registeruser", regPayload)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Çerezlerin döndüğünü doğrula
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			accessTokenCookie = c
		}
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
		}
	}
	assert.NotNil(t, accessTokenCookie, "access_token çerezi dönmeli")
	assert.NotNil(t, refreshTokenCookie, "refresh_token çerezi dönmeli")

	// DB'de kullanıcının ve sepetinin oluştuğunu doğrula
	ctx := context.Background()
	user, err := app.UserStore.GetUserByEmail(ctx, uniqueMail)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.UserID)
	assert.Equal(t, "Ahmet", user.Name)
	assert.Equal(t, "user", user.Role)

	cartID, err := app.CartStore.GetCartByUserID(ctx, user.UserID)
	assert.NoError(t, err)
	assert.NotEmpty(t, cartID, "Kullanıcıya ait sepet oluşturulmuş olmalı")

	// 2. Aynı e-posta ile tekrar kayıt olmaya çalışıldığında 409 dönmeli
	dupResp := app.doRequest(t, http.MethodPost, "/api/user/registeruser", regPayload)
	assert.Equal(t, http.StatusConflict, dupResp.StatusCode)
}

func TestIntegration_Auth_Login(t *testing.T) {
	app := setupIntegrationApp(t)

	email := fmt.Sprintf("login-%s@example.com", uuid.NewString()[:8])
	password := "SecretPass123!"
	app.createTestUser(t, email, password, "user")

	t.Run("geçerli bilgilerle giriş başarılı olmalı ve çerezleri dönmeli", func(t *testing.T) {
		loginPayload := structs.LoginReq{
			Mail:     email,
			Password: password,
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/login", loginPayload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var hasAccessToken, hasRefreshToken bool
		for _, c := range resp.Cookies() {
			if c.Name == "access_token" && c.Value != "" {
				hasAccessToken = true
			}
			if c.Name == "refresh_token" && c.Value != "" {
				hasRefreshToken = true
			}
		}
		assert.True(t, hasAccessToken, "access_token çerezi dönmeli")
		assert.True(t, hasRefreshToken, "refresh_token çerezi dönmeli")
	})

	t.Run("yanlış şifre ile giriş başarısız olmalı", func(t *testing.T) {
		loginPayload := structs.LoginReq{
			Mail:     email,
			Password: "WrongPassword!",
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/login", loginPayload)
		assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_Auth_GetUserRole(t *testing.T) {
	app := setupIntegrationApp(t)

	email := fmt.Sprintf("me-%s@example.com", uuid.NewString()[:8])
	_, accessCookie, _ := app.createTestUser(t, email, "Pass1234!", "user")

	t.Run("geçerli token ile kullanıcı rolü doğrulanmalı ve 200 dönmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/user/me", nil, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("token olmadan yetkisiz 401 dönmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/user/me", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestIntegration_Auth_RefreshToken(t *testing.T) {
	app := setupIntegrationApp(t)

	email := fmt.Sprintf("refresh-%s@example.com", uuid.NewString()[:8])
	_, _, refreshCookie := app.createTestUser(t, email, "Pass1234!", "user")

	t.Run("geçerli refresh token ile yeni tokenlar alınabilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodPost, "/api/refreshtoken", nil, refreshCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var newAccessCookie, newRefreshCookie *http.Cookie
		for _, c := range resp.Cookies() {
			if c.Name == "access_token" {
				newAccessCookie = c
			}
			if c.Name == "refresh_token" {
				newRefreshCookie = c
			}
		}

		assert.NotNil(t, newAccessCookie)
		assert.NotNil(t, newRefreshCookie)
		assert.NotEmpty(t, newAccessCookie.Value)
		assert.NotEmpty(t, newRefreshCookie.Value)
	})

	t.Run("geçersiz refresh token ile 401 dönmeli", func(t *testing.T) {
		fakeCookie := &http.Cookie{
			Name:  "refresh_token",
			Value: "invalid-token-string",
		}
		resp := app.doRequest(t, http.MethodPost, "/api/refreshtoken", nil, fakeCookie)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestIntegration_Auth_AdminAuthorization(t *testing.T) {
	app := setupIntegrationApp(t)

	adminEmail := fmt.Sprintf("admin-%s@example.com", uuid.NewString()[:8])
	userEmail := fmt.Sprintf("user-%s@example.com", uuid.NewString()[:8])

	_, adminAccessCookie, _ := app.createTestUser(t, adminEmail, "AdminPass123!", "admin")
	_, userAccessCookie, _ := app.createTestUser(t, userEmail, "UserPass123!", "user")

	t.Run("admin token'ı ile admin endpoint'ine erişilebilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/admin/getusercount", nil, adminAccessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var count int
		err := json.NewDecoder(resp.Body).Decode(&count)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, 1)
	})

	t.Run("normal kullanıcı admin endpoint'ine eriştiğinde 403 Forbidden dönmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/admin/getusercount", nil, userAccessCookie)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("kimlik doğrulaması olmadan erişildiğinde 401 Unauthorized dönmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/admin/getusercount", nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestIntegration_Auth_RateLimiter(t *testing.T) {
	app := setupIntegrationApp(t)

	// Burst = 10 olarak tanımlanmış login endpoint'ine aynı IP'den istek atıyoruz
	payload := structs.LoginReq{
		Mail:     "ratelimit-test@example.com",
		Password: "wrongpassword",
	}

	t.Run("burst limitine kadar istekler 429 almamalı", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := app.doRequest(t, http.MethodPost, "/api/user/login", payload)
			assert.NotEqual(t, http.StatusTooManyRequests, resp.StatusCode, "İlk 10 istek rate limit'e takılmamalı")
		}
	})

	t.Run("kova boşaldığında 11. istek 429 Too Many Requests dönmeli", func(t *testing.T) {
		rateLimitedResp := app.doRequest(t, http.MethodPost, "/api/user/login", payload)
		assert.Equal(t, http.StatusTooManyRequests, rateLimitedResp.StatusCode)
		assert.Equal(t, "60", rateLimitedResp.Header.Get("Retry-After"))
	})
}

