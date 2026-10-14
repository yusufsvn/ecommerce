package httphandlers

import (
	"bytes"
	"context"
	e "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
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

func setupUserHandlerTest() (*UserHandler, *mockUserRepo, *mockTokenRepo, *mockCartRepo, *mockAddressRepo) {
	userRepo := &mockUserRepo{}
	tokenRepo := &mockTokenRepo{}
	cartRepo := &mockCartRepo{}
	addrRepo := &mockAddressRepo{}

	tokenSvc := services.NewTokenService(tokenRepo, "super-secret-jwt-key-minimum-32-chars!")
	userSvc := services.NewUserService(userRepo, *tokenSvc, cartRepo, addrRepo, "test-client-id")
	handler := NewUserHandler(*userSvc)

	return handler, userRepo, tokenRepo, cartRepo, addrRepo
}

func TestUserHandler_Login(t *testing.T) {
	t.Run("başarılı girişte çerezleri ve kullanıcı adını dönmeli", func(t *testing.T) {
		handler, userRepo, tokenRepo, _, _ := setupUserHandlerTest()

		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			assert.Equal(t, "user@test.com", email)
			return &st.UserLoginRes{
				UserID:   "user-1",
				Name:     "Ahmet",
				Surname:  "Yilmaz",
				Password: string(hashedPw),
				Email:    "user@test.com",
				Role:     "user",
			}, nil
		}
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return nil
		}

		body := map[string]interface{}{
			"mail":     "user@test.com",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp st.TokenWithName
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Ahmet", resp.Name)
		assert.Equal(t, "Yilmaz", resp.Surname)

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
		assert.True(t, hasAccess)
		assert.True(t, hasRefresh)
	})

	t.Run("şifre yanlış olduğunda 500 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return &st.UserLoginRes{
				UserID:   "user-1",
				Password: string(hashedPw),
				Email:    "user@test.com",
				Role:     "user",
			}, nil
		}

		body := map[string]interface{}{
			"mail":     "user@test.com",
			"password": "wrong-password",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "giriş sırasında hata oluştu")
	})
}

func TestUserHandler_Register(t *testing.T) {
	t.Run("geçerli bilgilerle kayıt başarılı olmalı ve 201 dönmeli", func(t *testing.T) {
		handler, userRepo, tokenRepo, cartRepo, addrRepo := setupUserHandlerTest()

		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return nil, errors.New("user not found")
		}
		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			assert.Equal(t, "Mehmet", user.Name)
			return "user-new-id", nil
		}
		cartRepo.createFn = func(ctx context.Context, user_id string) (string, error) {
			return "cart-new-id", nil
		}
		addrRepo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			return nil
		}
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return nil
		}

		body := st.UserRegisterReq{
			Name:           "Mehmet",
			SecondName:     "Demir",
			Mail:           "mehmet@test.com",
			Password:       "GucluSifre123!",
			Password_check: "GucluSifre123!",
			Phone_number:   "05321234567",
			City:           "Bursa",
			Country:        "Turkiye",
			District:       "Nilufer",
			ZipCode:        "16000",
			Address:        "Ataevler Mah.",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
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
		assert.True(t, hasAccess)
		assert.True(t, hasRefresh)
	})

	t.Run("şifreler eşleşmediğinde 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _ := setupUserHandlerTest()

		body := st.UserRegisterReq{
			Name:           "Mehmet",
			SecondName:     "Demir",
			Mail:           "m@test.com",
			Password:       "Sifre123!",
			Password_check: "FarkliSifre!",
			Phone_number:   "05321234567",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Sunucu içi hata")
	})

	t.Run("kullanıcı oluşturulurken ErrUserExists hatasında 409 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			return "", e.ErrUserExists
		}

		body := st.UserRegisterReq{
			Name:           "Mehmet",
			SecondName:     "Demir",
			Mail:           "existing@test.com",
			Password:       "GucluSifre123!",
			Password_check: "GucluSifre123!",
			Phone_number:   "05321234567",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), e.ErrUserExists.Error())
	})
}

func TestUserHandler_GetUserRole(t *testing.T) {
	t.Run("geçerli kullanıcı için 200 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			assert.Equal(t, "user-123", userId)
			return "user", nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetUserRole(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("kullanıcı db'de bulunamazsa (isValid == false) 401 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			return "other-role", nil // rol uyuşmazlığı
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
		req = withUserClaims(req, "user-unknown", "user")
		rec := httptest.NewRecorder()

		handler.GetUserRole(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Kullanıcı bulunamadı")
	})

	t.Run("servis/db hatasında 500 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			return "", errors.New("db query failure")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
		req = withUserClaims(req, "user-123", "user")
		rec := httptest.NewRecorder()

		handler.GetUserRole(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Beklenmedik bir hata oluştu")
	})
}

func TestUserHandler_GetUserCount(t *testing.T) {
	t.Run("başarıyla kullanıcı sayısını dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.getUserCountFn = func(ctx context.Context) (int, error) {
			return 42, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/users/count", nil)
		rec := httptest.NewRecorder()

		handler.GetUserCount(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var count int
		err := json.Unmarshal(rec.Body.Bytes(), &count)
		assert.NoError(t, err)
		assert.Equal(t, 42, count)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		handler, userRepo, _, _, _ := setupUserHandlerTest()

		userRepo.getUserCountFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("count failed")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/users/count", nil)
		rec := httptest.NewRecorder()

		handler.GetUserCount(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Kullanıcı sayısı Getirilirken hata oluştu")
	})
}

func TestUserHandler_LoginWithGoogle(t *testing.T) {
	t.Run("Google ID token doğrulaması başarısız olduğunda 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _ := setupUserHandlerTest()

		body := st.LoginGoogle{
			ReturnIdToken: "invalid-google-token",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/user/google-login", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.LoginWithGoogle(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Google ile giriş yapılırken hata oluştu")
	})
}
