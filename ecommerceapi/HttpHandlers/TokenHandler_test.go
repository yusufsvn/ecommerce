package httphandlers

import (
	"context"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupTokenHandlerTest() (*TokenHandler, *mockTokenRepo) {
	repo := &mockTokenRepo{}
	svc := services.NewTokenService(repo, "secret-key-that-is-at-least-32-bytes-long!")
	handler := NewTokenHandler(*svc)
	return handler, repo
}

func TestTokenHandler_GetNewAccessToken(t *testing.T) {
	handler, repo := setupTokenHandlerTest()

	t.Run("geçerli refresh token claim ile yeni tokenlar ve çerezler dönmeli", func(t *testing.T) {
		repo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			assert.Equal(t, "jti-123", token_id)
			return &st.TokenQueryRes{
				UserId:     "user-1",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: false,
			}, nil
		}
		repo.updateTokenFn = func(ctx context.Context, token_id string, is_used bool) error {
			assert.Equal(t, "jti-123", token_id)
			assert.True(t, is_used)
			return nil
		}
		repo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			assert.Equal(t, "user-1", token.User_id)
			assert.Equal(t, "fam-1", token.Family_Id)
			return nil
		}

		claims := jwt.MapClaims{
			"jti":       "jti-123",
			"user_id":   "user-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		}

		req := httptest.NewRequest(http.MethodPost, "/api/token/refresh", nil)
		ctx := context.WithValue(req.Context(), Types.UserContextKey, claims)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.GetNewAccessToken(rec, req)

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

	t.Run("context içerisinde claim yoksa 401 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/token/refresh", nil)
		rec := httptest.NewRecorder()

		handler.GetNewAccessToken(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "geçersiz refresh token")
	})

	t.Run("token iptal edilmiş (revoked) ise 401 dönmeli", func(t *testing.T) {
		repo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "user-1",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: true, // iptal edilmiş
			}, nil
		}
		repo.panicUpdateFn = func(ctx context.Context, family_id string, is_revoked bool) error {
			return nil
		}

		claims := jwt.MapClaims{
			"jti":       "jti-revoked",
			"user_id":   "user-1",
			"role":      "user",
			"family_id": "fam-1",
		}

		req := httptest.NewRequest(http.MethodPost, "/api/token/refresh", nil)
		ctx := context.WithValue(req.Context(), Types.UserContextKey, claims)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.GetNewAccessToken(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "geçersiz refresh token")
	})

	t.Run("db hatası oluştuğunda 401 dönmeli", func(t *testing.T) {
		repo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return nil, errors.New("db error")
		}

		claims := jwt.MapClaims{
			"jti":       "jti-err",
			"user_id":   "user-1",
			"role":      "user",
			"family_id": "fam-1",
		}

		req := httptest.NewRequest(http.MethodPost, "/api/token/refresh", nil)
		ctx := context.WithValue(req.Context(), Types.UserContextKey, claims)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		handler.GetNewAccessToken(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "geçersiz refresh token")
	})
}
