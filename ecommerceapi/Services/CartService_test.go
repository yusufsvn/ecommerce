package services

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCartRepoForCartService struct {
	getCartByUserIDFn         func(ctx context.Context, user_id string) (string, error)
	createFn                  func(ctx context.Context, user_id string) (string, error)
	deleteFn                  func(ctx context.Context, cartId string) error
	getCartProductsByUserIDFn func(ctx context.Context, user_id string) ([]structs.CartProductDTO, error)
}

func (m *mockCartRepoForCartService) GetCartByUserID(ctx context.Context, user_id string) (string, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, user_id)
	}
	return "", nil
}

func (m *mockCartRepoForCartService) Create(ctx context.Context, user_id string) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user_id)
	}
	return "", nil
}

func (m *mockCartRepoForCartService) Delete(ctx context.Context, cartId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, cartId)
	}
	return nil
}

func (m *mockCartRepoForCartService) GetCartProductsByUserID(ctx context.Context, user_id string) ([]structs.CartProductDTO, error) {
	if m.getCartProductsByUserIDFn != nil {
		return m.getCartProductsByUserIDFn(ctx, user_id)
	}
	return nil, nil
}

type mockCartProductRepoForCartService struct {
	addFn                    func(ctx context.Context, cartProduct *st.CartProductDTO) error
	removeFn                 func(ctx context.Context, cartId string, productId string) error
	removeAllFn              func(ctx context.Context, cartId string) error
	updateQuantityFn         func(ctx context.Context, cartProduct *st.CartProductDTO) error
	getByCartIDFn            func(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error)
	getAllProductsFromCartFn func(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error)
}

func (m *mockCartProductRepoForCartService) Add(ctx context.Context, cartProduct *st.CartProductDTO) error {
	if m.addFn != nil {
		return m.addFn(ctx, cartProduct)
	}
	return nil
}

func (m *mockCartProductRepoForCartService) Remove(ctx context.Context, cartId string, productId string) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, cartId, productId)
	}
	return nil
}

func (m *mockCartProductRepoForCartService) RemoveAll(ctx context.Context, cartId string) error {
	if m.removeAllFn != nil {
		return m.removeAllFn(ctx, cartId)
	}
	return nil
}

func (m *mockCartProductRepoForCartService) UpdateQuantity(ctx context.Context, cartProduct *st.CartProductDTO) error {
	if m.updateQuantityFn != nil {
		return m.updateQuantityFn(ctx, cartProduct)
	}
	return nil
}

func (m *mockCartProductRepoForCartService) GetByCartID(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error) {
	if m.getByCartIDFn != nil {
		return m.getByCartIDFn(ctx, cartID)
	}
	return nil, nil
}

func (m *mockCartProductRepoForCartService) GetAllProductsFromCart(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error) {
	if m.getAllProductsFromCartFn != nil {
		return m.getAllProductsFromCartFn(ctx, cart_id)
	}
	return nil, nil
}

func TestCartService_GetUserCartId(t *testing.T) {
	cartRepo := &mockCartRepoForCartService{}
	cartProductRepo := &mockCartProductRepoForCartService{}
	svc := NewCartService(cartRepo, cartProductRepo)

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		res, err := svc.GetUserCartId(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Beklenmeyen bir hatayla karşılaşıldı")
	})

	t.Run("veritabanı GetCartByUserID hata verirse hata dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return "", errors.New("db cart error")
		}

		res, err := svc.GetUserCartId(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "sepet bulunamadı")
	})

	t.Run("başarıyla UserCartId dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			assert.Equal(t, userID, uID)
			return cartID, nil
		}

		res, err := svc.GetUserCartId(ctx)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, userID, res.User_id)
		assert.Equal(t, cartID, res.Cart_id)
	})
}

func TestCartService_AddToCart(t *testing.T) {
	cartRepo := &mockCartRepoForCartService{}
	cartProductRepo := &mockCartProductRepoForCartService{}
	svc := NewCartService(cartRepo, cartProductRepo)

	req := &st.AddToCartReq{
		ProductID: "prod-1",
		Quantity:  2,
	}

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		err := svc.AddToCart(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet bulunamadı")
	})

	t.Run("mevcut sepet ID yoksa yeni sepet oluşturup ürünü eklemeli", func(t *testing.T) {
		userID := uuid.NewString()
		newCartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return "", nil // sepet henüz yok
		}

		cartCreated := false
		cartRepo.createFn = func(ctx context.Context, uID string) (string, error) {
			cartCreated = true
			assert.Equal(t, userID, uID)
			return newCartID, nil
		}

		productAdded := false
		cartProductRepo.addFn = func(ctx context.Context, cartProduct *st.CartProductDTO) error {
			productAdded = true
			assert.Equal(t, newCartID, cartProduct.Cart_id)
			assert.Equal(t, req.ProductID, cartProduct.Product_id)
			assert.Equal(t, req.Quantity, cartProduct.Quantity)
			return nil
		}

		err := svc.AddToCart(ctx, req)
		assert.NoError(t, err)
		assert.True(t, cartCreated)
		assert.True(t, productAdded)
	})

	t.Run("sepet oluşturulurken hata çıkarsa hata dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return "", nil
		}
		cartRepo.createFn = func(ctx context.Context, uID string) (string, error) {
			return "", errors.New("db create cart error")
		}

		err := svc.AddToCart(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet oluşturulamadı")
	})

	t.Run("sepete ekleme işlemi hata verirse 'sepete eklenemedi' dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.addFn = func(ctx context.Context, cartProduct *st.CartProductDTO) error {
			return errors.New("db add error")
		}

		err := svc.AddToCart(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepete eklenemedi")
	})

	t.Run("mevcut sepete başarıyla eklemeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}

		productAdded := false
		cartProductRepo.addFn = func(ctx context.Context, cartProduct *st.CartProductDTO) error {
			productAdded = true
			assert.Equal(t, cartID, cartProduct.Cart_id)
			assert.Equal(t, "prod-1", cartProduct.Product_id)
			return nil
		}

		err := svc.AddToCart(ctx, req)
		assert.NoError(t, err)
		assert.True(t, productAdded)
	})
}

func TestCartService_RemoveFromCart(t *testing.T) {
	cartRepo := &mockCartRepoForCartService{}
	cartProductRepo := &mockCartProductRepoForCartService{}
	svc := NewCartService(cartRepo, cartProductRepo)

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		err := svc.RemoveFromCart(context.Background(), "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet bulunamadı")
	})

	t.Run("silme sorgusu hata verirse 'sepetten silinemedi' hatası dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.removeFn = func(ctx context.Context, cID string, pID string) error {
			return errors.New("db delete error")
		}

		err := svc.RemoveFromCart(ctx, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepetten silinemedi")
	})

	t.Run("başarıyla sepetten ürünü silmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}

		called := false
		cartProductRepo.removeFn = func(ctx context.Context, cID string, pID string) error {
			called = true
			assert.Equal(t, cartID, cID)
			assert.Equal(t, "prod-1", pID)
			return nil
		}

		err := svc.RemoveFromCart(ctx, "prod-1")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestCartService_UpdateQuantity(t *testing.T) {
	cartRepo := &mockCartRepoForCartService{}
	cartProductRepo := &mockCartProductRepoForCartService{}
	svc := NewCartService(cartRepo, cartProductRepo)

	req := &st.UpdateCartReq{
		Product_id: "prod-1",
		Quantity:   5,
	}

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		err := svc.UpdateQuantity(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet bulunamadı")
	})

	t.Run("adet güncelleme DB hatası verirse 'sepet güncellenemedi' dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.updateQuantityFn = func(ctx context.Context, dto *st.CartProductDTO) error {
			return errors.New("db update error")
		}

		err := svc.UpdateQuantity(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet güncellenemedi")
	})

	t.Run("başarıyla adet güncellemeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}

		called := false
		cartProductRepo.updateQuantityFn = func(ctx context.Context, dto *st.CartProductDTO) error {
			called = true
			assert.Equal(t, cartID, dto.Cart_id)
			assert.Equal(t, req.Product_id, dto.Product_id)
			assert.Equal(t, req.Quantity, dto.Quantity)
			return nil
		}

		err := svc.UpdateQuantity(ctx, req)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestCartService_GetCart(t *testing.T) {
	cartRepo := &mockCartRepoForCartService{}
	cartProductRepo := &mockCartProductRepoForCartService{}
	svc := NewCartService(cartRepo, cartProductRepo)

	t.Run("claims yoksa hata dönmeli", func(t *testing.T) {
		res, err := svc.GetCart(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("db hatası Customerrors.ZeroListResult ise doğrudan o hatayı dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, uID string) ([]structs.CartProductDTO, error) {
			return nil, Customerrors.ZeroListResult
		}

		res, err := svc.GetCart(ctx)
		assert.ErrorIs(t, err, Customerrors.ZeroListResult)
		assert.Nil(t, res)
	})

	t.Run("db sorgusu genel hata verirse hata dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, uID string) ([]structs.CartProductDTO, error) {
			return nil, errors.New("db error")
		}

		res, err := svc.GetCart(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "sepet bulunamadı")
	})

	t.Run("sepet boşsa (len == 0) Customerrors.ZeroListResult dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, uID string) ([]structs.CartProductDTO, error) {
			return []structs.CartProductDTO{}, nil
		}

		res, err := svc.GetCart(ctx)
		assert.ErrorIs(t, err, Customerrors.ZeroListResult)
		assert.Nil(t, res)
	})

	t.Run("başarıyla sepet ürünlerini dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := context.WithValue(context.Background(), Types.UserContextKey, jwt.MapClaims{"user_id": userID})

		cartRepo.getCartProductsByUserIDFn = func(ctx context.Context, uID string) ([]structs.CartProductDTO, error) {
			return []structs.CartProductDTO{
				{Cart_id: cartID, Product_id: "prod-1", Quantity: 2},
				{Cart_id: cartID, Product_id: "prod-2", Quantity: 1},
			}, nil
		}

		res, err := svc.GetCart(ctx)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, cartID, res.Cart_id)
		assert.Len(t, res.Products, 2)
		assert.Equal(t, "prod-1", res.Products[0].Product_id)
		assert.Equal(t, 2, res.Products[0].Quantity)
		assert.Equal(t, "prod-2", res.Products[1].Product_id)
		assert.Equal(t, 1, res.Products[1].Quantity)
	})
}
