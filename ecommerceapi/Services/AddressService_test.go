package services

import (
	"context"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAddressRepoForAddressService struct {
	createFn      func(ctx context.Context, address *st.CreateAddressDTO) error
	deleteFn      func(ctx context.Context, addressId string) error
	updateFn      func(ctx context.Context, address st.Address, user_id string) error
	getByUserIdFn func(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error)
}

func (m *mockAddressRepoForAddressService) Create(ctx context.Context, address *st.CreateAddressDTO) error {
	if m.createFn != nil {
		return m.createFn(ctx, address)
	}
	return nil
}

func (m *mockAddressRepoForAddressService) Delete(ctx context.Context, addressId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, addressId)
	}
	return nil
}

func (m *mockAddressRepoForAddressService) Update(ctx context.Context, address st.Address, user_id string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, address, user_id)
	}
	return nil
}

func (m *mockAddressRepoForAddressService) GetByUserId(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error) {
	if m.getByUserIdFn != nil {
		return m.getByUserIdFn(ctx, user_id)
	}
	return nil, nil
}

func TestAddressService_GetAddress(t *testing.T) {
	mockRepo := &mockAddressRepoForAddressService{}
	svc := NewAddressService(mockRepo)

	t.Run("claims yoksa veya user_id alınamazsa hata dönmeli", func(t *testing.T) {
		res, err := svc.GetAddress(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("veritabanı GetByUserId hata verirse 'adres alınamadı' hatası dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		mockRepo.getByUserIdFn = func(ctx context.Context, uID string) (*[]st.GetAddressDTO, error) {
			return nil, errors.New("db error")
		}

		res, err := svc.GetAddress(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "adres alınamadı")
	})

	t.Run("başarıyla kullanıcının adres listesini dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		expectedAddresses := []st.GetAddressDTO{
			{AddressId: "addr-1", Name: "Ev", City: "İstanbul", Country: "Türkiye", FullAddress: "Kadıköy"},
		}

		mockRepo.getByUserIdFn = func(ctx context.Context, uID string) (*[]st.GetAddressDTO, error) {
			assert.Equal(t, userID, uID)
			return &expectedAddresses, nil
		}

		res, err := svc.GetAddress(ctx)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, &expectedAddresses, res.Addresses)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	mockRepo := &mockAddressRepoForAddressService{}
	svc := NewAddressService(mockRepo)

	addr := st.Address{
		AddressName: "İş",
		AddressStr:  "Levent",
		City:        "İstanbul",
		Country:     "Türkiye",
		ZipCode:     "34330",
	}

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		err := svc.UpdateAddress(context.Background(), addr)
		assert.Error(t, err)
	})

	t.Run("veritabanı Update hata verirse 'db error' dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		mockRepo.updateFn = func(ctx context.Context, a st.Address, uID string) error {
			return errors.New("sql syntax error")
		}

		err := svc.UpdateAddress(ctx, addr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("başarıyla adresi güncellemeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		called := false
		mockRepo.updateFn = func(ctx context.Context, a st.Address, uID string) error {
			called = true
			assert.Equal(t, userID, uID)
			assert.Equal(t, addr, a)
			return nil
		}

		err := svc.UpdateAddress(ctx, addr)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestAddressService_CreateAddress(t *testing.T) {
	mockRepo := &mockAddressRepoForAddressService{}
	svc := NewAddressService(mockRepo)

	addr := st.Address{
		AddressName: "Yazlık",
		AddressStr:  "Bodrum",
		City:        "Muğla",
		Country:     "Türkiye",
		District:    "Bodrum",
		ZipCode:     "48400",
	}

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		err := svc.CreateAddress(context.Background(), addr)
		assert.Error(t, err)
	})

	t.Run("veritabanı Create hata verirse hatayı dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		mockRepo.createFn = func(ctx context.Context, dto *st.CreateAddressDTO) error {
			return errors.New("insert address error")
		}

		err := svc.CreateAddress(ctx, addr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insert address error")
	})

	t.Run("başarıyla adres DTO'sunu oluşturup kaydetmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		called := false
		mockRepo.createFn = func(ctx context.Context, dto *st.CreateAddressDTO) error {
			called = true
			assert.Equal(t, userID, dto.UserId)
			assert.Equal(t, addr.AddressName, dto.Name)
			assert.Equal(t, addr.AddressStr, dto.FullAddress)
			assert.Equal(t, addr.City, dto.City)
			assert.Equal(t, addr.Country, dto.Country)
			assert.Equal(t, addr.District, dto.District)
			assert.Equal(t, addr.ZipCode, dto.PostCode)
			return nil
		}

		err := svc.CreateAddress(ctx, addr)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestAddressService_DeleteAddress(t *testing.T) {
	mockRepo := &mockAddressRepoForAddressService{}
	svc := NewAddressService(mockRepo)
	ctx := context.Background()

	t.Run("veritabanı Delete hata verirse hatayı dönmeli", func(t *testing.T) {
		mockRepo.deleteFn = func(ctx context.Context, addressId string) error {
			return errors.New("delete error")
		}

		err := svc.DeleteAddress(ctx, "addr-123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})

	t.Run("başarıyla adresi silmeli", func(t *testing.T) {
		called := false
		mockRepo.deleteFn = func(ctx context.Context, addressId string) error {
			called = true
			assert.Equal(t, "addr-123", addressId)
			return nil
		}

		err := svc.DeleteAddress(ctx, "addr-123")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}
