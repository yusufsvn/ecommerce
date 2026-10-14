package repository

import (
	"context"
	e "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createUniqueTestUser testler için benzersiz bir kullanıcı oluşturur
func createUniqueTestUser(t *testing.T, ctx context.Context, repo *UserTable) (*Entities.User, string) {
	t.Helper()

	email := fmt.Sprintf("user_%s@test.com", uuid.NewString())
	user := &Entities.User{
		Name:     "TestName",
		Surname:  "TestSurname",
		Phone:    "05551112233",
		Email:    email,
		Password: "Password123!",
	}

	userID, err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, userID)

	uID, err := uuid.Parse(userID)
	require.NoError(t, err)
	user.UserID = uID

	return user, userID
}

func TestCreateUser(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("başarıyla kullanıcı oluşturmalı ve geçerli UUID dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		rawPassword := "SecurePass123!"
		user := &Entities.User{
			Name:     "Ahmet",
			Surname:  "Yılmaz",
			Phone:    "05559876543",
			Email:    fmt.Sprintf("ahmet_%s@test.com", uuid.NewString()),
			Password: rawPassword,
		}

		// act
		userID, err := userRepo.Create(ctx, user)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, userID)

		parsedUUID, err := uuid.Parse(userID)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, parsedUUID)

		// DB doğrulama
		var dbName, dbSurname, dbPhone, dbEmail, dbPassword, dbRole string
		query := `SELECT name, surname, phone, email, password, user_role FROM users WHERE user_id = $1`
		err = testDbInstance.QueryRow(ctx, query, parsedUUID).Scan(&dbName, &dbSurname, &dbPhone, &dbEmail, &dbPassword, &dbRole)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, dbName)
		assert.Equal(t, user.Surname, dbSurname)
		assert.Equal(t, user.Phone, dbPhone)
		assert.Equal(t, user.Email, dbEmail)
		assert.Equal(t, "user", dbRole)
		assert.True(t, Utils.CheckPasswordHash(rawPassword, dbPassword), "Veritabanındaki parola hash doğrulaması başarılı olmalı")
	})

	t.Run("aynı email ile kullanıcı oluşturulmaya çalışıldığında ErrUserExists dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		existingUser, _ := createUniqueTestUser(t, ctx, userRepo)

		duplicateUser := &Entities.User{
			Name:     "Duplicate",
			Surname:  "User",
			Phone:    "05550000000",
			Email:    existingUser.Email, // aynı email
			Password: "AnotherPassword123!",
		}

		// act
		userID, err := userRepo.Create(ctx, duplicateUser)

		// assert
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.True(t, errors.Is(err, e.ErrUserExists), "ErrUserExists hatası bekleniyordu")
	})

	t.Run("bağlam iptal edildiğinde (context canceled) hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		user := &Entities.User{
			Name:     "Cancel",
			Surname:  "Test",
			Phone:    "05551110000",
			Email:    fmt.Sprintf("cancel_%s@test.com", uuid.NewString()),
			Password: "Password123!",
		}

		// act
		userID, err := userRepo.Create(ctx, user)

		// assert
		assert.Error(t, err)
		assert.Empty(t, userID)
	})
}

func TestGetAllUsers(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("tüm kullanıcıları başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user1, _ := createUniqueTestUser(t, ctx, userRepo)
		user2, _ := createUniqueTestUser(t, ctx, userRepo)

		// act
		users, err := userRepo.GetAll(0, 0, ctx)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, users)

		found1, found2 := false, false
		for _, u := range users {
			if u.UserID == user1.UserID {
				found1 = true
				assert.Equal(t, user1.Email, u.Email)
				assert.Equal(t, user1.Name, u.Name)
				assert.Equal(t, user1.Surname, u.Surname)
			}
			if u.UserID == user2.UserID {
				found2 = true
				assert.Equal(t, user2.Email, u.Email)
			}
		}
		assert.True(t, found1, "user1 listede bulunmalı")
		assert.True(t, found2, "user2 listede bulunmalı")
	})

	t.Run("limit ve offset ile sayfalanmış kullanıcıları getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		createUniqueTestUser(t, ctx, userRepo)
		createUniqueTestUser(t, ctx, userRepo)

		// act
		page1, err1 := userRepo.GetAll(1, 0, ctx)
		page2, err2 := userRepo.GetAll(1, 1, ctx)

		// assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Len(t, page1, 1)
		assert.Len(t, page2, 1)
		assert.NotEqual(t, page1[0].UserID, page2[0].UserID, "Farklı sayfalar farklı kullanıcıları dönmeli")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		users, err := userRepo.GetAll(10, 0, ctx)

		// assert
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "sorgu sırasında hata")
	})
}

func TestDeleteUser(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("kullanıcıyı başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user, _ := createUniqueTestUser(t, ctx, userRepo)

		// act
		err := userRepo.Delete(ctx, &user.UserID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT count(*) FROM users WHERE user_id = $1", user.UserID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Kullanıcı veritabanından silinmiş olmalı")
	})

	t.Run("olmayan kullanıcı silinmeye çalışıldığında hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomUUID := uuid.New()

		// act
		err := userRepo.Delete(ctx, &randomUUID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		randomUUID := uuid.New()

		// act
		err := userRepo.Delete(ctx, &randomUUID)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kullanıcı silme hatası")
	})
}

func TestUpdateUser(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("kullanıcı bilgilerini ve parolasını başarıyla güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user, _ := createUniqueTestUser(t, ctx, userRepo)

		newPassword := "BrandNewPass999!"
		user.Name = "Güncellenmiş İsim"
		user.Surname = "Güncellenmiş Soyisim"
		user.Phone = "05321112233"
		user.Password = newPassword

		// act
		err := userRepo.Update(ctx, user)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, user.UpdatedAt)

		// DB doğrulama
		var dbName, dbSurname, dbPhone, dbPassword string
		query := `SELECT name, surname, phone, password FROM users WHERE user_id = $1`
		err = testDbInstance.QueryRow(ctx, query, user.UserID).Scan(&dbName, &dbSurname, &dbPhone, &dbPassword)
		assert.NoError(t, err)
		assert.Equal(t, "Güncellenmiş İsim", dbName)
		assert.Equal(t, "Güncellenmiş Soyisim", dbSurname)
		assert.Equal(t, "05321112233", dbPhone)
		assert.True(t, Utils.CheckPasswordHash(newPassword, dbPassword), "Yeni parolanın hash'i eşleşmeli")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		user := &Entities.User{
			UserID:   uuid.New(),
			Name:     "Test",
			Surname:  "Cancel",
			Phone:    "05550000000",
			Email:    "cancel@test.com",
			Password: "Password123!",
		}

		// act
		err := userRepo.Update(ctx, user)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "güncelleme sorgusu hatası")
	})
}

func TestGetUserByEmail(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("email ile kullanıcıyı başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user, _ := createUniqueTestUser(t, ctx, userRepo)

		// act
		res, err := userRepo.GetUserByEmail(ctx, user.Email)

		// assert
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, user.UserID.String(), res.UserID)
		assert.Equal(t, user.Name, res.Name)
		assert.Equal(t, user.Surname, res.Surname)
		assert.Equal(t, "user", res.Role)
		assert.True(t, Utils.CheckPasswordHash("Password123!", res.Password), "Dönen hash şifre ile eşleşmeli")
	})

	t.Run("olmayan email sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentEmail := fmt.Sprintf("nonexistent_%s@test.com", uuid.NewString())

		// act
		res, err := userRepo.GetUserByEmail(ctx, nonExistentEmail)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		res, err := userRepo.GetUserByEmail(ctx, "any@test.com")

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetUserById(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("kullanıcı id ile rolünü başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user, userID := createUniqueTestUser(t, ctx, userRepo)

		// act
		role, err := userRepo.GetUserById(ctx, userID)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, "user", role)
		_ = user
	})

	t.Run("olmayan kullanıcı id sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomID := uuid.NewString()

		// act
		role, err := userRepo.GetUserById(ctx, randomID)

		// assert
		assert.Error(t, err)
		assert.Empty(t, role)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		role, err := userRepo.GetUserById(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Empty(t, role)
	})
}

func TestGetUser(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("kullanıcı id ile UserWithoutAddress detaylarını başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		user, userID := createUniqueTestUser(t, ctx, userRepo)

		// act
		res, err := userRepo.GetUser(ctx, userID)

		// assert
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, user.Name, res.Name)
		assert.Equal(t, user.Surname, res.Surname)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, user.Phone, res.Phone)
	})

	t.Run("olmayan kullanıcı id sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomID := uuid.NewString()

		// act
		res, err := userRepo.GetUser(ctx, randomID)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "kullanıcı bulunamadı")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		res, err := userRepo.GetUser(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetUserCount(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("kullanıcı sayısını doğru dönmeli (admin hariç)", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// Mevcut toplam satır sayısını al
		var currentTotalRows int
		err := testDbInstance.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&currentTotalRows)
		require.NoError(t, err)

		// 1 kullanıcı ekle
		createUniqueTestUser(t, ctx, userRepo)

		// act
		count, err := userRepo.GetUserCount(ctx)

		// assert
		assert.NoError(t, err)
		// GetUserCount returns `count - 1` (excluding admin)
		assert.Equal(t, (currentTotalRows+1)-1, count)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := userRepo.GetUserCount(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestAddNewUserWithGoogle(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("Google OAuth kullanıcısını başarıyla eklemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		googleSub := "google_oauth_" + uuid.NewString()
		email := fmt.Sprintf("google_%s@gmail.com", uuid.NewString())
		googleUser := &st.GoogleUser{
			UserId:  googleSub,
			Email:   email,
			Name:    "GoogleAd",
			Surname: "GoogleSoyad",
		}

		// act
		userID, err := userRepo.AddNewUserWithGoogle(ctx, googleUser)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, userID)

		parsedUUID, err := uuid.Parse(userID)
		assert.NoError(t, err)

		// DB doğrulama
		var dbName, dbSurname, dbEmail, dbProvider, dbOauthId, dbRole string
		query := `SELECT name, surname, email, oauth_provider, oauth_id, user_role FROM users WHERE user_id = $1`
		err = testDbInstance.QueryRow(ctx, query, parsedUUID).Scan(&dbName, &dbSurname, &dbEmail, &dbProvider, &dbOauthId, &dbRole)
		assert.NoError(t, err)
		assert.Equal(t, googleUser.Name, dbName)
		assert.Equal(t, googleUser.Surname, dbSurname)
		assert.Equal(t, googleUser.Email, dbEmail)
		assert.Equal(t, "google", dbProvider)
		assert.Equal(t, googleSub, dbOauthId)
		assert.Equal(t, "user", dbRole)
	})

	t.Run("aynı email ile tekrar eklenmeye çalışıldığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		email := fmt.Sprintf("google_dup_%s@gmail.com", uuid.NewString())
		googleUser1 := &st.GoogleUser{
			UserId:  "sub_1_" + uuid.NewString(),
			Email:   email,
			Name:    "İlk",
			Surname: "Kullanıcı",
		}
		_, err := userRepo.AddNewUserWithGoogle(ctx, googleUser1)
		require.NoError(t, err)

		googleUser2 := &st.GoogleUser{
			UserId:  "sub_2_" + uuid.NewString(),
			Email:   email, // aynı email
			Name:    "İkinci",
			Surname: "Kullanıcı",
		}

		// act
		userID, err := userRepo.AddNewUserWithGoogle(ctx, googleUser2)

		// assert
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		googleUser := &st.GoogleUser{
			UserId:  "sub_" + uuid.NewString(),
			Email:   "cancelled@gmail.com",
			Name:    "Cancel",
			Surname: "Test",
		}

		// act
		userID, err := userRepo.AddNewUserWithGoogle(ctx, googleUser)

		// assert
		assert.Error(t, err)
		assert.Empty(t, userID)
	})
}

func TestGetUserByOauthId(t *testing.T) {
	userRepo := NewUserStore(testDbInstance)

	t.Run("Google OAuth ID ile kullanıcıyı başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		googleSub := "google_oauth_" + uuid.NewString()
		email := fmt.Sprintf("oauth_find_%s@gmail.com", uuid.NewString())
		googleUser := &st.GoogleUser{
			UserId:  googleSub,
			Email:   email,
			Name:    "OAuthAd",
			Surname: "OAuthSoyad",
		}

		createdID, err := userRepo.AddNewUserWithGoogle(ctx, googleUser)
		require.NoError(t, err)

		// act
		res, err := userRepo.GetUserByOauthId(ctx, googleSub)

		// assert
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, createdID, res.UserID)
		assert.Equal(t, googleUser.Name, res.Name)
		assert.Equal(t, googleUser.Surname, res.Surname)
		assert.Equal(t, googleUser.Email, res.Email)
		assert.Equal(t, "user", res.Role)
	})

	t.Run("olmayan OAuth ID sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomSub := "nonexistent_sub_" + uuid.NewString()

		// act
		res, err := userRepo.GetUserByOauthId(ctx, randomSub)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("farklı provider ile kayıtlı olan oauth_id google sorgusunda bulunamamalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		githubSub := "github_sub_" + uuid.NewString()
		email := fmt.Sprintf("github_%s@test.com", uuid.NewString())

		// provider'ı google olmayan bir kayıt ekle
		query := `INSERT INTO users (name, surname, email, oauth_provider, oauth_id) VALUES ($1, $2, $3, 'github', $4)`
		_, err := testDbInstance.Exec(ctx, query, "Git", "Hub", email, githubSub)
		require.NoError(t, err)

		// act (GetUserByOauthId sorgusu oauth_provider='google' kontrolü yapar)
		res, err := userRepo.GetUserByOauthId(ctx, githubSub)

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		res, err := userRepo.GetUserByOauthId(ctx, "any-oauth-id")

		// assert
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
