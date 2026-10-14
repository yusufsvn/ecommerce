package services

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTokenRepoForTokenService struct {
	addTokenFn     func(ctx context.Context, token st.RefreshTokenSaveDTO) error
	removeTokenFn  func(ctx context.Context, token_id string) error
	getTokenByIdFn func(ctx context.Context, token_id string) (*st.TokenQueryRes, error)
	updateTokenFn  func(ctx context.Context, token_id string, is_used bool) error
	panicUpdateFn  func(ctx context.Context, family_id string, is_revoked bool) error
}

func (m *mockTokenRepoForTokenService) AddToken(ctx context.Context, token st.RefreshTokenSaveDTO) error {
	if m.addTokenFn != nil {
		return m.addTokenFn(ctx, token)
	}
	return nil
}

func (m *mockTokenRepoForTokenService) RemoveToken(ctx context.Context, token_id string) error {
	if m.removeTokenFn != nil {
		return m.removeTokenFn(ctx, token_id)
	}
	return nil
}

func (m *mockTokenRepoForTokenService) GetTokenById(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
	if m.getTokenByIdFn != nil {
		return m.getTokenByIdFn(ctx, token_id)
	}
	return nil, nil
}

func (m *mockTokenRepoForTokenService) UpdateToken(ctx context.Context, token_id string, is_used bool) error {
	if m.updateTokenFn != nil {
		return m.updateTokenFn(ctx, token_id, is_used)
	}
	return nil
}

func (m *mockTokenRepoForTokenService) PanicUpdate(ctx context.Context, family_id string, is_revoked bool) error {
	if m.panicUpdateFn != nil {
		return m.panicUpdateFn(ctx, family_id, is_revoked)
	}
	return nil
}

func TestTokenService_GenToken(t *testing.T) {
	secretKey := "my-super-secret-test-key-32bytes!"
	mockRepo := &mockTokenRepoForTokenService{}
	svc := NewTokenService(mockRepo, secretKey)

	role := "admin"
	userID := uuid.NewString()
	familyID := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour).Truncate(time.Second)

	tokenDTO, err := svc.GenToken(role, expiresAt, familyID, userID)
	assert.NoError(t, err)
	require.NotNil(t, tokenDTO)
	assert.NotEmpty(t, tokenDTO.Token)
	assert.NotEmpty(t, tokenDTO.Jti)
	assert.Equal(t, expiresAt, tokenDTO.Expires_at)

	// Token'ı parse edip claim'leri kontrol et
	parsedToken, err := jwt.Parse(tokenDTO.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, role, claims["role"])
	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, familyID, claims["family_id"])
	assert.Equal(t, tokenDTO.Jti, claims["jti"])
	assert.Equal(t, "my-ecommerce", claims["iss"])
}

func TestTokenService_SaveRefreshToken(t *testing.T) {
	secretKey := "test-secret"
	mockRepo := &mockTokenRepoForTokenService{}
	svc := NewTokenService(mockRepo, secretKey)
	ctx := context.Background()

	t.Run("veritabanı AddToken hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return errors.New("db insert error")
		}

		err := svc.SaveRefreshToken(ctx, "sample-token", "user-1", time.Now().Add(time.Hour), "jti-1", "fam-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refresh token db ye kaydedilemedi")
	})

	t.Run("başarıyla refresh token'ı hashleyip veritabanına kaydetmeli", func(t *testing.T) {
		called := false
		expiresAt := time.Now().Add(time.Hour)
		mockRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			called = true
			assert.Equal(t, "jti-1", token.Jti)
			assert.Equal(t, "user-1", token.User_id)
			assert.Equal(t, "fam-1", token.Family_Id)
			assert.False(t, token.Is_used)
			assert.False(t, token.Is_revoked)
			assert.Equal(t, expiresAt, token.Expires_at)
			assert.NotEmpty(t, token.Token_hash)
			return nil
		}

		err := svc.SaveRefreshToken(ctx, "sample-token", "user-1", expiresAt, "jti-1", "fam-1")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestTokenService_RefreshingToken(t *testing.T) {
	secretKey := "my-refresh-test-key-32-chars-long"
	mockRepo := &mockTokenRepoForTokenService{}
	svc := NewTokenService(mockRepo, secretKey)

	t.Run("context içinde claims yoksa hata dönmeli", func(t *testing.T) {
		res, err := svc.RefreshingToken(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "cannot fetch token claims")
	})

	t.Run("claim alanları eksikse hata dönmeli", func(t *testing.T) {
		// jti eksik
		ctx1 := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "f-1",
		})
		_, err1 := svc.RefreshingToken(ctx1)
		assert.Error(t, err1)
		assert.Contains(t, err1.Error(), "jti is not a string")

		// user_id eksik
		ctx2 := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"role":      "user",
			"family_id": "f-1",
		})
		_, err2 := svc.RefreshingToken(ctx2)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "user_id is not a string")

		// role eksik
		ctx3 := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"family_id": "f-1",
		})
		_, err3 := svc.RefreshingToken(ctx3)
		assert.Error(t, err3)
		assert.Contains(t, err3.Error(), "role is not a string")

		// family_id eksik
		ctx4 := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":     "jti-1",
			"user_id": "u-1",
			"role":    "user",
		})
		_, err4 := svc.RefreshingToken(ctx4)
		assert.Error(t, err4)
		assert.Contains(t, err4.Error(), "family id is not a string")
	})

	t.Run("veritabanından token getirilemezse hata dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return nil, errors.New("token not found")
		}

		res, err := svc.RefreshingToken(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "token db den getirelemedi")
	})

	t.Run("token kullanılmışsa (IsUsed == true) PanicToken çağrılmalı ve ErrUnauthorized dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "u-1",
				FamilyId:   "fam-1",
				IsUsed:     true,
				Is_revoked: false,
			}, nil
		}

		panicCalled := false
		mockRepo.panicUpdateFn = func(ctx context.Context, family_id string, is_revoked bool) error {
			panicCalled = true
			assert.Equal(t, "fam-1", family_id)
			return nil
		}

		res, err := svc.RefreshingToken(ctx)
		assert.ErrorIs(t, err, Customerrors.ErrUnauthorized)
		assert.Nil(t, res)
		assert.True(t, panicCalled, "Kullanılmış token tespit edildiğinde PanicToken çağrılmalı")
	})

	t.Run("token iptal edilmişse (Is_revoked == true) PanicToken çağrılmalı ve ErrUnauthorized dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "u-1",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: true,
			}, nil
		}

		panicCalled := false
		mockRepo.panicUpdateFn = func(ctx context.Context, family_id string, is_revoked bool) error {
			panicCalled = true
			assert.Equal(t, "fam-1", family_id)
			return nil
		}

		res, err := svc.RefreshingToken(ctx)
		assert.ErrorIs(t, err, Customerrors.ErrUnauthorized)
		assert.Nil(t, res)
		assert.True(t, panicCalled)
	})

	t.Run("veritabanındaki token UserId ile talepteki user_id uyuşmuyorsa ErrUnauthorized dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "different-user",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: false,
			}, nil
		}

		res, err := svc.RefreshingToken(ctx)
		assert.ErrorIs(t, err, Customerrors.ErrUnauthorized)
		assert.Nil(t, res)
	})

	t.Run("token güncellemesi hata verirse hata dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "u-1",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: false,
			}, nil
		}
		mockRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return nil
		}
		mockRepo.updateTokenFn = func(ctx context.Context, token_id string, is_used bool) error {
			return errors.New("db update error")
		}

		res, err := svc.RefreshingToken(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "update sorgusunda hatayla karşılaşıldı")
	})

	t.Run("başarıyla yeni access ve refresh token dönmeli ve eski token'ı used yapmalı", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"jti":       "jti-1",
			"user_id":   "u-1",
			"role":      "user",
			"family_id": "fam-1",
			"exp":       float64(time.Now().Add(time.Hour).Unix()),
		})

		mockRepo.getTokenByIdFn = func(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
			return &st.TokenQueryRes{
				UserId:     "u-1",
				FamilyId:   "fam-1",
				IsUsed:     false,
				Is_revoked: false,
			}, nil
		}

		newRefreshSaved := false
		mockRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			newRefreshSaved = true
			assert.Equal(t, "u-1", token.User_id)
			assert.Equal(t, "fam-1", token.Family_Id)
			assert.False(t, token.Is_used)
			assert.False(t, token.Is_revoked)
			return nil
		}

		oldTokenUpdated := false
		mockRepo.updateTokenFn = func(ctx context.Context, token_id string, is_used bool) error {
			oldTokenUpdated = true
			assert.Equal(t, "jti-1", token_id)
			assert.True(t, is_used)
			return nil
		}

		res, err := svc.RefreshingToken(ctx)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Access_token)
		assert.NotEmpty(t, res.Refresh_token)
		assert.True(t, newRefreshSaved)
		assert.True(t, oldTokenUpdated)
	})
}

func TestTokenService_PanicToken(t *testing.T) {
	secretKey := "test-secret"
	mockRepo := &mockTokenRepoForTokenService{}
	svc := NewTokenService(mockRepo, secretKey)
	ctx := context.Background()

	t.Run("veritabanı PanicUpdate hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.panicUpdateFn = func(ctx context.Context, family_id string, is_revoked bool) error {
			return errors.New("db panic error")
		}

		err := svc.PanicToken(ctx, "fam-1", true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Panic token update sorgusunda hata")
	})

	t.Run("başarıyla PanicUpdate çağırmalı", func(t *testing.T) {
		called := false
		mockRepo.panicUpdateFn = func(ctx context.Context, family_id string, is_revoked bool) error {
			called = true
			assert.Equal(t, "fam-1", family_id)
			assert.True(t, is_revoked)
			return nil
		}

		err := svc.PanicToken(ctx, "fam-1", true)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}
