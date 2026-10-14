package services

import (
	"context"
	e "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"ecommerceapi/Utils"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepoForUserService struct {
	getAllFn               func(limit int, offset int, ctx context.Context) ([]Entities.User, error)
	createFn               func(ctx context.Context, user *Entities.User) (string, error)
	deleteFn               func(ctx context.Context, userId *uuid.UUID) error
	updateFn               func(ctx context.Context, user *Entities.User) error
	getUserByEmailFn       func(ctx context.Context, email string) (*st.UserLoginRes, error)
	getUserByIdFn          func(ctx context.Context, userId string) (string, error)
	getUserFn              func(ctx context.Context, user_id string) (*st.UserWithoutAddress, error)
	getUserCountFn         func(ctx context.Context) (int, error)
	addNewUserWithGoogleFn func(ctx context.Context, user *st.GoogleUser) (string, error)
	getUserByOauthIdFn     func(ctx context.Context, oauthId string) (*st.UserLoginRes, error)
}

func (m *mockUserRepoForUserService) GetAll(limit int, offset int, ctx context.Context) ([]Entities.User, error) {
	if m.getAllFn != nil {
		return m.getAllFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockUserRepoForUserService) Create(ctx context.Context, user *Entities.User) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return "", nil
}
func (m *mockUserRepoForUserService) Delete(ctx context.Context, userId *uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userId)
	}
	return nil
}
func (m *mockUserRepoForUserService) Update(ctx context.Context, user *Entities.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}
func (m *mockUserRepoForUserService) GetUserByEmail(ctx context.Context, email string) (*st.UserLoginRes, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}
func (m *mockUserRepoForUserService) GetUserById(ctx context.Context, userId string) (string, error) {
	if m.getUserByIdFn != nil {
		return m.getUserByIdFn(ctx, userId)
	}
	return "", nil
}
func (m *mockUserRepoForUserService) GetUser(ctx context.Context, user_id string) (*st.UserWithoutAddress, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, user_id)
	}
	return nil, nil
}
func (m *mockUserRepoForUserService) GetUserCount(ctx context.Context) (int, error) {
	if m.getUserCountFn != nil {
		return m.getUserCountFn(ctx)
	}
	return 0, nil
}
func (m *mockUserRepoForUserService) AddNewUserWithGoogle(ctx context.Context, user *st.GoogleUser) (string, error) {
	if m.addNewUserWithGoogleFn != nil {
		return m.addNewUserWithGoogleFn(ctx, user)
	}
	return "", nil
}
func (m *mockUserRepoForUserService) GetUserByOauthId(ctx context.Context, oauthId string) (*st.UserLoginRes, error) {
	if m.getUserByOauthIdFn != nil {
		return m.getUserByOauthIdFn(ctx, oauthId)
	}
	return nil, nil
}

func setupUserServiceTest() (*UserService, *mockUserRepoForUserService, *mockCartRepoForCartService, *mockAddressRepoForAddressService, *mockTokenRepoForTokenService) {
	userRepo := &mockUserRepoForUserService{}
	cartRepo := &mockCartRepoForCartService{}
	addressRepo := &mockAddressRepoForAddressService{}
	tokenRepo := &mockTokenRepoForTokenService{}

	tokenService := NewTokenService(tokenRepo, "my-super-secret-jwt-key-32-chars!")
	userService := NewUserService(userRepo, *tokenService, cartRepo, addressRepo, "google-client-id-123")

	return userService, userRepo, cartRepo, addressRepo, tokenRepo
}

func TestUserService_RegisterUser(t *testing.T) {
	svc, userRepo, cartRepo, _, tokenRepo := setupUserServiceTest()
	ctx := context.Background()

	t.Run("parolalar eşleşmiyorsa ErrPasswordNotMatched dönmeli", func(t *testing.T) {
		req := &st.UserRegisterReq{
			Mail:           "test@example.com",
			Password:       "pass123",
			Password_check: "pass456",
		}
		res, err := svc.RegisterUser(ctx, req)
		assert.ErrorIs(t, err, e.ErrPasswordNotMatched)
		assert.Nil(t, res)
	})

	t.Run("kullanıcı oluşturma db hatası verirse hata dönmeli", func(t *testing.T) {
		req := &st.UserRegisterReq{
			Mail:           "TEST@EXAMPLE.COM ",
			Password:       "pass123",
			Password_check: "pass123",
			Name:           "Ali",
			SecondName:     "Veli",
		}

		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			assert.Equal(t, "test@example.com", user.Email, "Email normalize edilmeli (lowercase ve trim)")
			return "", errors.New("db create user error")
		}

		res, err := svc.RegisterUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "user oluşturulamadi")
	})

	t.Run("sepet oluşturma db hatası verirse hata dönmeli", func(t *testing.T) {
		req := &st.UserRegisterReq{
			Mail:           "ali@test.com",
			Password:       "pass123",
			Password_check: "pass123",
			Name:           "Ali",
			SecondName:     "Veli",
		}

		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			return "user-123", nil
		}
		cartRepo.createFn = func(ctx context.Context, userId string) (string, error) {
			return "", errors.New("db create cart error")
		}

		res, err := svc.RegisterUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "user oluşturulamadi")
	})

	t.Run("refresh token db kaydı başarısız olursa hata dönmeli", func(t *testing.T) {
		req := &st.UserRegisterReq{
			Mail:           "ali@test.com",
			Password:       "pass123",
			Password_check: "pass123",
			Name:           "Ali",
			SecondName:     "Veli",
		}

		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			return "user-123", nil
		}
		cartRepo.createFn = func(ctx context.Context, userId string) (string, error) {
			return "cart-123", nil
		}
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return errors.New("db add token error")
		}

		res, err := svc.RegisterUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "token veritabanına kaydedilemedi")
	})

	t.Run("başarıyla kayıt olup token dönmeli", func(t *testing.T) {
		req := &st.UserRegisterReq{
			Mail:           "  Ali@Test.COM  ",
			Password:       "pass123",
			Password_check: "pass123",
			Name:           "Ali",
			SecondName:     "Veli",
		}

		userCreated := false
		userRepo.createFn = func(ctx context.Context, user *Entities.User) (string, error) {
			userCreated = true
			assert.Equal(t, "ali@test.com", user.Email)
			assert.Equal(t, "Ali", user.Name)
			return "user-123", nil
		}

		cartCreated := false
		cartRepo.createFn = func(ctx context.Context, userId string) (string, error) {
			cartCreated = true
			assert.Equal(t, "user-123", userId)
			return "cart-123", nil
		}

		tokenSaved := false
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			tokenSaved = true
			assert.Equal(t, "user-123", token.User_id)
			return nil
		}

		res, err := svc.RegisterUser(ctx, req)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Access_token)
		assert.NotEmpty(t, res.Refresh_token)
		assert.Equal(t, "Ali", res.Name)
		assert.Equal(t, "Veli", res.Surname)
		assert.True(t, userCreated)
		assert.True(t, cartCreated)
		assert.True(t, tokenSaved)
	})
}

func TestUserService_Login(t *testing.T) {
	svc, userRepo, _, _, tokenRepo := setupUserServiceTest()
	ctx := context.Background()

	hashedPassword, err := Utils.HashPassword("correct-pass")
	require.NoError(t, err)

	t.Run("kullanıcı bulunamazsa hata dönmeli", func(t *testing.T) {
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return nil, errors.New("user not found")
		}

		res, err := svc.Login(ctx, st.LoginReq{Mail: "notfound@test.com", Password: "any"}, "user")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "giriş sırasında hata oluştu")
	})

	t.Run("kullanıcı rolü eşleşmiyorsa hata dönmeli", func(t *testing.T) {
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return &st.UserLoginRes{
				UserID:   "u-1",
				Role:     "user",
				Password: hashedPassword,
			}, nil
		}

		res, err := svc.Login(ctx, st.LoginReq{Mail: "ali@test.com", Password: "correct-pass"}, "admin")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "giriş sırasında hata oluştu")
	})

	t.Run("parola yanlışsa 'Parola yanlış' hatası dönmeli", func(t *testing.T) {
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return &st.UserLoginRes{
				UserID:   "u-1",
				Role:     "user",
				Password: hashedPassword,
			}, nil
		}

		res, err := svc.Login(ctx, st.LoginReq{Mail: "ali@test.com", Password: "wrong-password"}, "user")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Parola yanlış")
	})

	t.Run("başarıyla giriş yapıp token dönmeli", func(t *testing.T) {
		userRepo.getUserByEmailFn = func(ctx context.Context, email string) (*st.UserLoginRes, error) {
			return &st.UserLoginRes{
				UserID:   "u-1",
				Name:     "Ali",
				Surname:  "Veli",
				Role:     "user",
				Password: hashedPassword,
			}, nil
		}
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			assert.Equal(t, "u-1", token.User_id)
			return nil
		}

		res, err := svc.Login(ctx, st.LoginReq{Mail: "ali@test.com", Password: "correct-pass"}, "user")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.NotEmpty(t, res.Access_token)
		assert.NotEmpty(t, res.Refresh_token)
		assert.Equal(t, "Ali", res.Name)
		assert.Equal(t, "Veli", res.Surname)
	})
}

func TestUserService_IsValidUser(t *testing.T) {
	svc, userRepo, _, _, _ := setupUserServiceTest()

	t.Run("claims yoksa veya user_id boşsa hata dönmeli", func(t *testing.T) {
		valid, err := svc.IsValidUser(context.Background())
		assert.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("veritabanı GetUserById hata verirse false ve hata dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"user_id": "u-1",
			"role":    "user",
		})

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			return "", errors.New("db error")
		}

		valid, err := svc.IsValidUser(ctx)
		assert.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("rol eşleşmiyorsa false ve nil hata dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"user_id": "u-1",
			"role":    "admin", // claim admin diyor
		})

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			return "user", nil // db user diyor
		}

		valid, err := svc.IsValidUser(ctx)
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("rol eşleşiyorsa true ve nil hata dönmeli", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{
			"user_id": "u-1",
			"role":    "admin",
		})

		userRepo.getUserByIdFn = func(ctx context.Context, userId string) (string, error) {
			return "admin", nil
		}

		valid, err := svc.IsValidUser(ctx)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
}

func TestUserService_GetUserCount(t *testing.T) {
	svc, userRepo, _, _, _ := setupUserServiceTest()
	ctx := context.Background()

	t.Run("hata durumunda -1 ve hatayı dönmeli", func(t *testing.T) {
		userRepo.getUserCountFn = func(ctx context.Context) (int, error) {
			return -1, errors.New("count error")
		}

		count, err := svc.GetUserCount(ctx)
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})

	t.Run("başarıyla kullanıcı sayısını dönmeli", func(t *testing.T) {
		userRepo.getUserCountFn = func(ctx context.Context) (int, error) {
			return 42, nil
		}

		count, err := svc.GetUserCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 42, count)
	})
}

func TestUserService_VerifyGoogleIdToken(t *testing.T) {
	svc, _, _, _, _ := setupUserServiceTest()
	ctx := context.Background()

	t.Run("token boşsa ErrEmptyParameter dönmeli", func(t *testing.T) {
		user, err := svc.VerifyGoogleIdToken(ctx, "")
		assert.ErrorIs(t, err, e.ErrEmptyParameter)
		assert.Nil(t, user)
	})
}

func TestUserService_LoginWithGoogle(t *testing.T) {
	svc, _, _, _, _ := setupUserServiceTest()
	ctx := context.Background()

	t.Run("token boşsa ErrEmptyParameter dönmeli", func(t *testing.T) {
		res, err := svc.LoginWithGoogle(ctx, "")
		assert.ErrorIs(t, err, e.ErrEmptyParameter)
		assert.Nil(t, res)
	})
}

func TestUserService_InitToken(t *testing.T) {
	svc, _, _, _, tokenRepo := setupUserServiceTest()
	ctx := context.Background()

	t.Run("başarıyla token üretip db'ye kaydetmeli", func(t *testing.T) {
		tokenSaved := false
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			tokenSaved = true
			assert.Equal(t, "user-999", token.User_id)
			return nil
		}

		tokens, err := svc.InitToken(ctx, "user-999", "Mehmet", "Kaya")
		assert.NoError(t, err)
		require.NotNil(t, tokens)
		assert.NotEmpty(t, tokens.Access_token)
		assert.NotEmpty(t, tokens.Refresh_token)
		assert.Equal(t, "Mehmet", tokens.Name)
		assert.Equal(t, "Kaya", tokens.Surname)
		assert.True(t, tokenSaved)
	})

	t.Run("refresh token db kaydı hata verirse hata dönmeli", func(t *testing.T) {
		tokenRepo.addTokenFn = func(ctx context.Context, token st.RefreshTokenSaveDTO) error {
			return errors.New("db error")
		}

		tokens, err := svc.InitToken(ctx, "user-999", "Mehmet", "Kaya")
		assert.Error(t, err)
		assert.Nil(t, tokens)
		assert.Contains(t, err.Error(), "token veritabanına kaydedilemedi")
	})
}
