package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestUserForToken foreign key kısıtını sağlamak için izole bir kullanıcı oluşturur
func createTestUserForToken(t *testing.T, ctx context.Context, db *pgxpool.Pool) string {
	t.Helper()

	userID := uuid.NewString()
	email := fmt.Sprintf("token_user_%s@test.com", uuid.NewString())
	query := `INSERT INTO users (user_id, password, name, surname, phone, email) VALUES ($1, 'pass123', 'TokenUser', 'Test', '05550001122', $2)`

	_, err := db.Exec(ctx, query, userID, email)
	if err != nil {
		t.Fatalf("Test kullanıcısı oluşturulamadı: %v", err)
	}

	return userID
}

func TestAddToken(t *testing.T) {
	tokenRepo := NewTokenStore(testDbInstance)

	t.Run("başarıyla refresh token eklemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		tokenID := uuid.NewString()
		familyID := uuid.NewString()
		tokenHash := "hash_value_123456"
		expiresAt := time.Now().Add(7 * 24 * time.Hour).Truncate(time.Microsecond).UTC()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        tokenID,
			Family_Id:  familyID,
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: expiresAt,
			Token_hash: tokenHash,
		}

		// act
		err := tokenRepo.AddToken(ctx, tokenDTO)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var dbUserID, dbFamilyID, dbTokenHash string
		var dbIsUsed, dbIsRevoked bool
		var dbExpiresAt time.Time
		query := `SELECT user_id, family_id, token_hash, is_used, is_revoked, expires_at FROM refreshtokens WHERE token_id = $1`
		err = testDbInstance.QueryRow(ctx, query, tokenID).Scan(&dbUserID, &dbFamilyID, &dbTokenHash, &dbIsUsed, &dbIsRevoked, &dbExpiresAt)

		assert.NoError(t, err)
		assert.Equal(t, userID, dbUserID)
		assert.Equal(t, familyID, dbFamilyID)
		assert.Equal(t, tokenHash, dbTokenHash)
		assert.False(t, dbIsUsed)
		assert.False(t, dbIsRevoked)
		assert.WithinDuration(t, expiresAt, dbExpiresAt.UTC(), time.Second)
	})

	t.Run("olmayan user_id ile eklenmeye çalışıldığında foreign key hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentUserID := uuid.NewString()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        uuid.NewString(),
			Family_Id:  uuid.NewString(),
			User_id:    nonExistentUserID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "some_hash",
		}

		// act
		err := tokenRepo.AddToken(ctx, tokenDTO)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token db'ye eklenemedi")
	})

	t.Run("aynı token_id (jti) ile tekrar eklenmeye çalışıldığında primary key çakışma hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		tokenID := uuid.NewString()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        tokenID,
			Family_Id:  uuid.NewString(),
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "hash_first",
		}

		err := tokenRepo.AddToken(ctx, tokenDTO)
		require.NoError(t, err)

		// Aynı tokenID ile tekrar ekleme
		duplicateDTO := tokenDTO
		duplicateDTO.Token_hash = "hash_second"

		// act
		err = tokenRepo.AddToken(ctx, duplicateDTO)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token db'ye eklenemedi")
	})

	t.Run("bağlam iptal edildiğinde (context canceled) hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        uuid.NewString(),
			Family_Id:  uuid.NewString(),
			User_id:    uuid.NewString(),
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "cancel_hash",
		}

		// act
		err := tokenRepo.AddToken(ctx, tokenDTO)

		// assert
		assert.Error(t, err)
	})
}

func TestGetTokenById(t *testing.T) {
	tokenRepo := NewTokenStore(testDbInstance)

	t.Run("var olan token_id ile token detaylarını başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		tokenID := uuid.NewString()
		familyID := uuid.NewString()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        tokenID,
			Family_Id:  familyID,
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "query_hash",
		}
		err := tokenRepo.AddToken(ctx, tokenDTO)
		require.NoError(t, err)

		// act
		res, err := tokenRepo.GetTokenById(ctx, tokenID)

		// assert
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, userID, res.UserId)
		assert.Equal(t, familyID, res.FamilyId)
		assert.False(t, res.IsUsed)
		assert.False(t, res.Is_revoked)
	})

	t.Run("olmayan token_id sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentID := uuid.NewString()

		// act
		res, err := tokenRepo.GetTokenById(ctx, nonExistentID)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "struct a dönüştürmede hata")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		res, err := tokenRepo.GetTokenById(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestRemoveToken(t *testing.T) {
	tokenRepo := NewTokenStore(testDbInstance)

	t.Run("var olan token_id'yi başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		tokenID := uuid.NewString()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        tokenID,
			Family_Id:  uuid.NewString(),
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "delete_hash",
		}
		err := tokenRepo.AddToken(ctx, tokenDTO)
		require.NoError(t, err)

		// act
		err = tokenRepo.RemoveToken(ctx, tokenID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT count(*) FROM refreshtokens WHERE token_id = $1", tokenID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Token veritabanından silinmiş olmalı")
	})

	t.Run("olmayan token_id silinmek istendiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomTokenID := uuid.NewString()

		// act
		err := tokenRepo.RemoveToken(ctx, randomTokenID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := tokenRepo.RemoveToken(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token silinirken hatayla karsilasildi")
	})
}

func TestUpdateToken(t *testing.T) {
	tokenRepo := NewTokenStore(testDbInstance)

	t.Run("token'ın is_used durumunu başarıyla güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		tokenID := uuid.NewString()

		tokenDTO := st.RefreshTokenSaveDTO{
			Jti:        tokenID,
			Family_Id:  uuid.NewString(),
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "update_hash",
		}
		err := tokenRepo.AddToken(ctx, tokenDTO)
		require.NoError(t, err)

		// act
		err = tokenRepo.UpdateToken(ctx, tokenID, true)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var isUsed bool
		err = testDbInstance.QueryRow(ctx, "SELECT is_used FROM refreshtokens WHERE token_id = $1", tokenID).Scan(&isUsed)
		assert.NoError(t, err)
		assert.True(t, isUsed, "is_used alanı true olarak güncellenmiş olmalı")
	})

	t.Run("olmayan token_id ile çağrıldığında hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomTokenID := uuid.NewString()

		// act
		err := tokenRepo.UpdateToken(ctx, randomTokenID, true)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := tokenRepo.UpdateToken(ctx, uuid.NewString(), true)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update işleminde hata")
	})
}

func TestPanicUpdate(t *testing.T) {
	tokenRepo := NewTokenStore(testDbInstance)

	t.Run("aynı family_id'ye sahip birden fazla token'ı is_revoked olarak güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		userID := createTestUserForToken(t, ctx, testDbInstance)
		familyID := uuid.NewString()
		otherFamilyID := uuid.NewString()

		// Aynı family_id'ye ait 2 token ekle
		token1ID := uuid.NewString()
		token2ID := uuid.NewString()
		token3ID := uuid.NewString() // Farklı family_id

		err := tokenRepo.AddToken(ctx, st.RefreshTokenSaveDTO{
			Jti:        token1ID,
			Family_Id:  familyID,
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "panic_hash_1",
		})
		require.NoError(t, err)

		err = tokenRepo.AddToken(ctx, st.RefreshTokenSaveDTO{
			Jti:        token2ID,
			Family_Id:  familyID,
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "panic_hash_2",
		})
		require.NoError(t, err)

		// Farklı family_id'ye sahip token ekle
		err = tokenRepo.AddToken(ctx, st.RefreshTokenSaveDTO{
			Jti:        token3ID,
			Family_Id:  otherFamilyID,
			User_id:    userID,
			Is_used:    false,
			Is_revoked: false,
			Expires_at: time.Now().Add(24 * time.Hour),
			Token_hash: "other_family_hash",
		})
		require.NoError(t, err)

		// act
		err = tokenRepo.PanicUpdate(ctx, familyID, true)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: token1 ve token2 iptal edilmiş olmalı
		var isRevoked1, isRevoked2, isRevoked3 bool
		err = testDbInstance.QueryRow(ctx, "SELECT is_revoked FROM refreshtokens WHERE token_id = $1", token1ID).Scan(&isRevoked1)
		assert.NoError(t, err)
		assert.True(t, isRevoked1, "token1 is_revoked true olmalı")

		err = testDbInstance.QueryRow(ctx, "SELECT is_revoked FROM refreshtokens WHERE token_id = $1", token2ID).Scan(&isRevoked2)
		assert.NoError(t, err)
		assert.True(t, isRevoked2, "token2 is_revoked true olmalı")

		// token3 ise hala false olmalı
		err = testDbInstance.QueryRow(ctx, "SELECT is_revoked FROM refreshtokens WHERE token_id = $1", token3ID).Scan(&isRevoked3)
		assert.NoError(t, err)
		assert.False(t, isRevoked3, "farklı family_id'deki token3 iptal edilmemeli (false kalmalı)")
	})

	t.Run("olmayan family_id ile çağrıldığında hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomFamilyID := uuid.NewString()

		// act
		err := tokenRepo.PanicUpdate(ctx, randomFamilyID, true)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := tokenRepo.PanicUpdate(ctx, uuid.NewString(), true)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update işleminde hata")
	})
}
