package services

import (
	"context"
	constants "ecommerceapi/Constants"
	e "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockOrderRepo struct {
	getAllFn                func(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error)
	createFn                func(ctx context.Context, order *st.OrderDTO) error
	deleteFn                func(ctx context.Context, orderId uuid.UUID) error
	getOrdersByUserIdFn     func(ctx context.Context, user_id string) ([]st.GetOrder, error)
	updateOrderStateFn      func(ctx context.Context, newState string, order_id string) error
	getOrderCountFn         func(ctx context.Context) (int, error)
	updateMultiStateFn      func(ctx context.Context, idlist []string, new_state string) error
	getOrderCountByStatusFn func(ctx context.Context, status string) (int, error)
	getOrdersByStatusFn     func(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error)
	getOrderByOrderCodeFn   func(ctx context.Context, order_code string) (*st.GetOrderAddress, error)
	getOrderCountWaitingFn  func(ctx context.Context) (int, error)
	searchOrderFn           func(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error)
	searchOrderCountFn      func(ctx context.Context, search_query string) (int, error)
}

func (m *mockOrderRepo) GetAll(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, limit, offset)
	}
	return nil, nil
}
func (m *mockOrderRepo) Create(ctx context.Context, order *st.OrderDTO) error {
	if m.createFn != nil {
		return m.createFn(ctx, order)
	}
	return nil
}
func (m *mockOrderRepo) Delete(ctx context.Context, orderId uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, orderId)
	}
	return nil
}
func (m *mockOrderRepo) GetOrdersByUserId(ctx context.Context, user_id string) ([]st.GetOrder, error) {
	if m.getOrdersByUserIdFn != nil {
		return m.getOrdersByUserIdFn(ctx, user_id)
	}
	return nil, nil
}
func (m *mockOrderRepo) UpdateOrderState(ctx context.Context, newState string, order_id string) error {
	if m.updateOrderStateFn != nil {
		return m.updateOrderStateFn(ctx, newState, order_id)
	}
	return nil
}
func (m *mockOrderRepo) GetOrderCount(ctx context.Context) (int, error) {
	if m.getOrderCountFn != nil {
		return m.getOrderCountFn(ctx)
	}
	return 0, nil
}
func (m *mockOrderRepo) UpdateMultiState(ctx context.Context, idlist []string, new_state string) error {
	if m.updateMultiStateFn != nil {
		return m.updateMultiStateFn(ctx, idlist, new_state)
	}
	return nil
}
func (m *mockOrderRepo) GetOrderCountByStatus(ctx context.Context, status string) (int, error) {
	if m.getOrderCountByStatusFn != nil {
		return m.getOrderCountByStatusFn(ctx, status)
	}
	return 0, nil
}
func (m *mockOrderRepo) GetOrdersByStatus(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {
	if m.getOrdersByStatusFn != nil {
		return m.getOrdersByStatusFn(ctx, limit, offset, status)
	}
	return nil, nil
}
func (m *mockOrderRepo) GetOrderByOrderCode(ctx context.Context, order_code string) (*st.GetOrderAddress, error) {
	if m.getOrderByOrderCodeFn != nil {
		return m.getOrderByOrderCodeFn(ctx, order_code)
	}
	return nil, nil
}
func (m *mockOrderRepo) GetOrderCountForWaiting(ctx context.Context) (int, error) {
	if m.getOrderCountWaitingFn != nil {
		return m.getOrderCountWaitingFn(ctx)
	}
	return 0, nil
}
func (m *mockOrderRepo) SearchOrder(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error) {
	if m.searchOrderFn != nil {
		return m.searchOrderFn(ctx, limit, search_query)
	}
	return nil, nil
}
func (m *mockOrderRepo) SearchOrderCount(ctx context.Context, search_query string) (int, error) {
	if m.searchOrderCountFn != nil {
		return m.searchOrderCountFn(ctx, search_query)
	}
	return 0, nil
}

type mockCartProductRepo struct {
	getAllProductsFromCartFn func(ctx context.Context, cartId string) ([]st.ProductsFromCartDb, error)
	removeAllFn              func(ctx context.Context, cartId string) error
	getByCartIDFn            func(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error)
	addFn                    func(ctx context.Context, cartProduct *st.CartProductDTO) error
	removeFn                 func(ctx context.Context, cartId string, productId string) error
	updateQuantityFn         func(ctx context.Context, cartProduct *st.CartProductDTO) error
}

func (m *mockCartProductRepo) GetAllProductsFromCart(ctx context.Context, cartId string) ([]st.ProductsFromCartDb, error) {
	if m.getAllProductsFromCartFn != nil {
		return m.getAllProductsFromCartFn(ctx, cartId)
	}
	return nil, nil
}
func (m *mockCartProductRepo) RemoveAll(ctx context.Context, cartId string) error {
	if m.removeAllFn != nil {
		return m.removeAllFn(ctx, cartId)
	}
	return nil
}
func (m *mockCartProductRepo) GetByCartID(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error) {
	if m.getByCartIDFn != nil {
		return m.getByCartIDFn(ctx, cartID)
	}
	return nil, nil
}
func (m *mockCartProductRepo) Add(ctx context.Context, cartProduct *st.CartProductDTO) error {
	if m.addFn != nil {
		return m.addFn(ctx, cartProduct)
	}
	return nil
}
func (m *mockCartProductRepo) Remove(ctx context.Context, cartId string, productId string) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, cartId, productId)
	}
	return nil
}
func (m *mockCartProductRepo) UpdateQuantity(ctx context.Context, cartProduct *st.CartProductDTO) error {
	if m.updateQuantityFn != nil {
		return m.updateQuantityFn(ctx, cartProduct)
	}
	return nil
}

type mockUserRepo struct {
	getUserFn func(ctx context.Context, user_id string) (*st.UserWithoutAddress, error)
}

func (m *mockUserRepo) GetAll(limit int, offset int, ctx context.Context) ([]Entities.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(ctx context.Context, user *Entities.User) (string, error) {
	return "", nil
}
func (m *mockUserRepo) Delete(ctx context.Context, userId *uuid.UUID) error { return nil }
func (m *mockUserRepo) Update(ctx context.Context, user *Entities.User) error { return nil }
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*st.UserLoginRes, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserById(ctx context.Context, userId string) (string, error) {
	return "", nil
}
func (m *mockUserRepo) GetUser(ctx context.Context, user_id string) (*st.UserWithoutAddress, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, user_id)
	}
	return nil, nil
}
func (m *mockUserRepo) GetUserCount(ctx context.Context) (int, error) { return 0, nil }
func (m *mockUserRepo) AddNewUserWithGoogle(ctx context.Context, user *st.GoogleUser) (string, error) {
	return "", nil
}
func (m *mockUserRepo) GetUserByOauthId(ctx context.Context, oauthId string) (*st.UserLoginRes, error) {
	return nil, nil
}

type mockProductOrderRepo struct {
	getByOrderIDFn func(ctx context.Context, orderID string) ([]st.OrderFullDetail, error)
	addFn          func(ctx context.Context, item Entities.ProductOrder) error
	addManyFn      func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error
}

func (m *mockProductOrderRepo) GetByOrderID(ctx context.Context, orderID string) ([]st.OrderFullDetail, error) {
	if m.getByOrderIDFn != nil {
		return m.getByOrderIDFn(ctx, orderID)
	}
	return nil, nil
}
func (m *mockProductOrderRepo) Add(ctx context.Context, item Entities.ProductOrder) error {
	if m.addFn != nil {
		return m.addFn(ctx, item)
	}
	return nil
}
func (m *mockProductOrderRepo) AddMany(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
	if m.addManyFn != nil {
		return m.addManyFn(ctx, items, order_id)
	}
	return nil
}

type mockProductProductCategoryRepo struct {
	getAllProductCategoriesFn func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error)
}

func (m *mockProductProductCategoryRepo) GetByProductID(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	return nil, nil
}
func (m *mockProductProductCategoryRepo) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	return nil, nil
}
func (m *mockProductProductCategoryRepo) GetAllProductCategories(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
	if m.getAllProductCategoriesFn != nil {
		return m.getAllProductCategoriesFn(ctx, ids)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepo) Add(ctx context.Context, item *Entities.ProductProductCategory) error {
	return nil
}
func (m *mockProductProductCategoryRepo) Remove(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
	return nil
}
func (m *mockProductProductCategoryRepo) GetByProductIDOneRow(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}
func (m *mockProductProductCategoryRepo) AddMany(ctx context.Context, items []string, product_id string) error {
	return nil
}
func (m *mockProductProductCategoryRepo) RemoveMany(ctx context.Context, items []string, product_id string) error {
	return nil
}

type mockCartRepo struct {
	getCartByUserIDFn func(ctx context.Context, userId string) (string, error)
}

func (m *mockCartRepo) GetCartByUserID(ctx context.Context, userId string) (string, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, userId)
	}
	return "", nil
}
func (m *mockCartRepo) GetCartProductsByUserID(ctx context.Context, user_id string) ([]structs.CartProductDTO, error) {
	return nil, nil
}
func (m *mockCartRepo) Create(ctx context.Context, userId string) (string, error) {
	return "", nil
}
func (m *mockCartRepo) Delete(ctx context.Context, cartId string) error { return nil }

type mockAddressRepo struct {
	createFn func(ctx context.Context, address *st.CreateAddressDTO) error
}

func (m *mockAddressRepo) Create(ctx context.Context, address *st.CreateAddressDTO) error {
	if m.createFn != nil {
		return m.createFn(ctx, address)
	}
	return nil
}
func (m *mockAddressRepo) Delete(ctx context.Context, addressId string) error { return nil }
func (m *mockAddressRepo) GetByUserId(ctx context.Context, userId string) (*[]st.GetAddressDTO, error) {
	return nil, nil
}
func (m *mockAddressRepo) Update(ctx context.Context, address st.Address, addressId string) error {
	return nil
}
func (m *mockAddressRepo) GetByAddressId(ctx context.Context, address_id string) (*st.GetAddressDTO, error) {
	return nil, nil
}

type mockPaymentFunctions struct {
	createPaymentNonTDSFn  func(req structs.CreateOrderDTO) bool
	createPaymentWith3DSFn func(data structs.CreateOrderDTO) (string, error)
}

func (m *mockPaymentFunctions) CreatePaymentNonTDS(req structs.CreateOrderDTO) bool {
	if m.createPaymentNonTDSFn != nil {
		return m.createPaymentNonTDSFn(req)
	}
	return false
}
func (m *mockPaymentFunctions) CreatePaymentWith3DS(data structs.CreateOrderDTO) (string, error) {
	if m.createPaymentWith3DSFn != nil {
		return m.createPaymentWith3DSFn(data)
	}
	return "", nil
}

// contextWithUserClaim helper'ı context içine user_id claim'i ekler
func contextWithUserClaim(ctx context.Context, userID string) context.Context {
	claims := jwt.MapClaims{
		"user_id": userID,
	}
	return context.WithValue(ctx, Types.UserContextKey, claims)
}

// setupOrderServiceHelper tüm mock bağımlılıklarıyla birlikte bir OrderService üretir
func setupOrderServiceHelper() (
	*OrderService,
	*mockOrderRepo,
	*mockCartProductRepo,
	*mockUserRepo,
	*mockProductOrderRepo,
	*mockProductProductCategoryRepo,
	*mockCartRepo,
	*mockAddressRepo,
	*mockPaymentFunctions,
) {
	orderRepo := &mockOrderRepo{}
	cartProductRepo := &mockCartProductRepo{}
	userRepo := &mockUserRepo{}
	productOrderRepo := &mockProductOrderRepo{}
	ppcRepo := &mockProductProductCategoryRepo{}
	cartRepo := &mockCartRepo{}
	addressRepo := &mockAddressRepo{}
	paymentFuncs := &mockPaymentFunctions{}

	cartService := NewCartService(cartRepo, cartProductRepo)
	addressService := NewAddressService(addressRepo)

	svc := NewOrderService(
		cartService,
		cartProductRepo,
		userRepo,
		paymentFuncs,
		orderRepo,
		productOrderRepo,
		nil,
		ppcRepo,
		"http://storage.api",
		"http://public.storage",
		"test-bucket",
		*addressService,
	)

	return svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, addressRepo, paymentFuncs
}

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("telefon numarası geçersizse (53 ile başlamıyorsa veya boşsa) ErrInvalidPhoneNumber dönmeli", func(t *testing.T) {
		svc, _, _, _, _, _, _, _, _ := setupOrderServiceHelper()
		ctx := context.Background()

		// Boş telefon
		err1 := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: ""})
		assert.ErrorIs(t, err1, e.ErrInvalidPhoneNumber)

		// 5 ile başlamayan telefon (örn. 0 ile başlayan)
		err2 := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "05421234567"})
		assert.ErrorIs(t, err2, e.ErrInvalidPhoneNumber)
	})

	t.Run("context içinde claims yoksa sepet ID alınamaz ve hata dönmeli", func(t *testing.T) {
		svc, _, _, _, _, _, _, _, _ := setupOrderServiceHelper()
		ctx := context.Background() // no claims

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "id error")
	})

	t.Run("sepet ürünleri getirilemezse hata dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, _, cartRepo, _, _ := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return nil, errors.New("db cart error")
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com", Phone: "0532"}, nil
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet verileri getirelemedi")
	})

	t.Run("kategori bilgisi getirilemezse hata dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, _ := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: uuid.NewString(), ProductName: "Ürün 1", Price: 100, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return nil, errors.New("category fetch error")
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category fetch error")
	})

	t.Run("kullanıcı bilgisi alınamazsa hata dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, _ := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 100, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return nil, errors.New("user not found")
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user cannot fetched")
	})

	t.Run("toplam tutar 0 veya eksi ise hata dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, _ := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 0, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total price 0")
	})

	t.Run("ödeme başarısız olduğunda hata dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 500, Quantity: 2, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return false // ödeme başarısız
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233", SecurePay: false})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ödeme başarısız")
	})

	t.Run("3DS ödeme talep edildiğinde desteklenmediği hatasını dönmeli", func(t *testing.T) {
		svc, _, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, _ := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 500, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233", SecurePay: true})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not supported 3ds payment")
	})

	t.Run("sipariş oluşturma (orderFunc.Create) başarısız olduğunda hata dönmeli", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, _, ppcRepo, cartRepo, _, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 500, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			return errors.New("db create order error")
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order save error")
	})

	t.Run("sipariş ürünleri ekleme (productorders.AddMany) başarısız olduğunda hata dönmeli", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, _, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 500, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			return nil
		}
		productOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			return errors.New("addMany error")
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "beklenmeyen hatayla karşılaşıldı")
	})

	t.Run("sepet temizleme (cartProducts.RemoveAll) başarısız olduğunda hata dönmeli", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, _, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Ürün 1", Price: 500, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Elektronik"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ali", Surname: "Veli", Email: "ali@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error { return nil }
		productOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			return nil
		}
		cartProductRepo.removeAllFn = func(ctx context.Context, cID string) error {
			return errors.New("removeAll error")
		}

		err := svc.CreateOrder(ctx, st.CreatePaymentDto{PhoneNumber: "5321112233"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sepet ürünleri silinemedi")
	})

	t.Run("tüm adımlar başarılı olduğunda ve SaveAddress true iken adres kaydedilip nil dönmeli", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, addressRepo, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Laptop", Price: 15000, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Bilgisayar"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ahmet", Surname: "Yılmaz", Email: "ahmet@test.com", Phone: "05321112233"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderCreated := false
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			orderCreated = true
			assert.Equal(t, userID, order.User_id)
			assert.Equal(t, constants.OrderStatePending, order.Status)
			assert.Equal(t, "15000", order.TotalAmount)
			return nil
		}
		productOrdersAdded := false
		productOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			productOrdersAdded = true
			assert.Len(t, items, 1)
			assert.Equal(t, "Bilgisayar", items[0].Category)
			return nil
		}
		cartCleaned := false
		cartProductRepo.removeAllFn = func(ctx context.Context, cID string) error {
			cartCleaned = true
			return nil
		}
		addressSaved := false
		addressRepo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			addressSaved = true
			return nil
		}

		req := st.CreatePaymentDto{
			PhoneNumber: "5321112233",
			SaveAddress: true,
			ShippingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
			BillingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
		}

		err := svc.CreateOrder(ctx, req)

		assert.NoError(t, err)
		assert.True(t, orderCreated, "Sipariş veritabanına eklenmeli")
		assert.True(t, productOrdersAdded, "Sipariş ürünleri eklenmeli")
		assert.True(t, cartCleaned, "Sepet temizlenmeli")
		assert.True(t, addressSaved, "Adres kaydedilmeli")
	})

	t.Run("SaveAddress false iken adres kaydedilmemeli ve işlem başarılı olmalı", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, addressRepo, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Mouse", Price: 300, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Aksesuar"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ayşe", Surname: "Demir", Email: "ayse@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			return nil
		}
		productOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			return nil
		}
		cartProductRepo.removeAllFn = func(ctx context.Context, cID string) error {
			return nil
		}
		addressSaved := false
		addressRepo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			addressSaved = true
			return nil
		}

		req := st.CreatePaymentDto{
			PhoneNumber: "5321112233",
			SaveAddress: false,
			ShippingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
			BillingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
		}

		err := svc.CreateOrder(ctx, req)
		assert.NoError(t, err)
		assert.False(t, addressSaved, "SaveAddress false iken adres kaydedilmemeli")
	})

	t.Run("SaveAddress true iken adres kaydı başarısız olursa hata dönmeli", func(t *testing.T) {
		svc, orderRepo, cartProductRepo, userRepo, productOrderRepo, ppcRepo, cartRepo, addressRepo, paymentFuncs := setupOrderServiceHelper()
		userID := uuid.NewString()
		cartID := uuid.NewString()
		prodID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)

		cartRepo.getCartByUserIDFn = func(ctx context.Context, uID string) (string, error) {
			return cartID, nil
		}
		cartProductRepo.getAllProductsFromCartFn = func(ctx context.Context, cID string) ([]st.ProductsFromCartDb, error) {
			return []st.ProductsFromCartDb{
				{ProductID: prodID, ProductName: "Klavye", Price: 800, Quantity: 1, ItemType: "PHYSICAL"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: prodID, Categories: []structs.CategoryDbJson{{CategoryName: "Aksesuar"}}},
			}, nil
		}
		userRepo.getUserFn = func(ctx context.Context, uID string) (*st.UserWithoutAddress, error) {
			return &st.UserWithoutAddress{Name: "Ayşe", Surname: "Demir", Email: "ayse@test.com"}, nil
		}
		paymentFuncs.createPaymentNonTDSFn = func(req structs.CreateOrderDTO) bool {
			return true
		}
		orderRepo.createFn = func(ctx context.Context, order *st.OrderDTO) error {
			return nil
		}
		productOrderRepo.addManyFn = func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
			return nil
		}
		cartProductRepo.removeAllFn = func(ctx context.Context, cID string) error {
			return nil
		}
		addressRepo.createFn = func(ctx context.Context, address *st.CreateAddressDTO) error {
			return errors.New("address save error")
		}

		req := st.CreatePaymentDto{
			PhoneNumber: "5321112233",
			SaveAddress: true,
			ShippingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
			BillingAddress: st.Address{
				AddressName: "Ev",
				AddressStr:  "Kadıköy",
				City:        "İstanbul",
				Country:     "Türkiye",
			},
		}

		err := svc.CreateOrder(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address save error")
	})
}

func TestOrderService_GetOrdersByUserId(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()

	t.Run("claims içinde user_id yoksa hata dönmeli", func(t *testing.T) {
		ctx := context.Background()
		orders, err := svc.GetOrdersByUserId(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Contains(t, err.Error(), "beklenmeyen hatayla karşılaşıldı")
	})

	t.Run("veritabanı sorgusu hata verirse hata dönmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)
		orderRepo.getOrdersByUserIdFn = func(ctx context.Context, uID string) ([]st.GetOrder, error) {
			return nil, errors.New("db error")
		}

		orders, err := svc.GetOrdersByUserId(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("başarıyla kullanıcının siparişlerini getirmeli", func(t *testing.T) {
		userID := uuid.NewString()
		ctx := contextWithUserClaim(context.Background(), userID)
		expectedOrders := []st.GetOrder{
			{
				OrderId:       "ORD12345",
				UserId:        userID,
				TotalPrice:    3500,
				Status:        constants.OrderStatePending,
				CreatedAt:     time.Now(),
				TotalQuantity: 3,
				OrderCode:     "ORD12345",
			},
		}
		orderRepo.getOrdersByUserIdFn = func(ctx context.Context, uID string) ([]st.GetOrder, error) {
			assert.Equal(t, userID, uID)
			return expectedOrders, nil
		}

		orders, err := svc.GetOrdersByUserId(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
	})
}

func TestOrderService_GetOrderDetails(t *testing.T) {
	svc, orderRepo, _, _, productOrderRepo, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("order_code ile sipariş bulunamazsa hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderByOrderCodeFn = func(ctx context.Context, orderCode string) (*st.GetOrderAddress, error) {
			return nil, errors.New("order not found")
		}

		res, err := svc.GetOrderDetails(ctx, "NONEXIST")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("sipariş ürün detayları getirilirken hata oluşursa hata dönmeli", func(t *testing.T) {
		orderID := uuid.NewString()
		orderRepo.getOrderByOrderCodeFn = func(ctx context.Context, orderCode string) (*st.GetOrderAddress, error) {
			return &st.GetOrderAddress{
				OrderId:         orderID,
				ShippingAddress: "Adres 1",
				BillingAddress:  "Adres 2",
			}, nil
		}
		productOrderRepo.getByOrderIDFn = func(ctx context.Context, oID string) ([]st.OrderFullDetail, error) {
			return nil, errors.New("productorder error")
		}

		res, err := svc.GetOrderDetails(ctx, "ORD123")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "beklenmeyen hatayla karşılaşıldı")
	})

	t.Run("başarıyla sipariş detaylarını ve görsel URL'lerini biçimlendirerek dönmeli", func(t *testing.T) {
		orderID := uuid.NewString()
		orderRepo.getOrderByOrderCodeFn = func(ctx context.Context, orderCode string) (*st.GetOrderAddress, error) {
			return &st.GetOrderAddress{
				OrderId:         orderID,
				ShippingAddress: "Kargo Adresi",
				BillingAddress:  "Fatura Adresi",
			}, nil
		}
		productOrderRepo.getByOrderIDFn = func(ctx context.Context, oID string) ([]st.OrderFullDetail, error) {
			return []st.OrderFullDetail{
				{
					ProductID:       uuid.NewString(),
					ProductName:     "Telefon",
					PriceAtPurchase: 10000,
					Quantity:        1,
					ThumbnailImage:  "image.jpg",
				},
			}, nil
		}

		res, err := svc.GetOrderDetails(ctx, "ORD123")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "Kargo Adresi", res.ShippingAddress)
		assert.Equal(t, "Fatura Adresi", res.BillingAddress)
		assert.Len(t, res.Products, 1)
		assert.Contains(t, res.Products[0].ThumbnailImage, "test-bucket")
	})
}

func TestOrderService_GetAllOrders(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("GetOrderCount hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountFn = func(ctx context.Context) (int, error) {
			return -1, errors.New("count error")
		}

		res, err := svc.GetAllOrders(ctx, "10", "0")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Beklenmeyen Hatayla Karşılaşıldı")
	})

	t.Run("GetAll sorgusu hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountFn = func(ctx context.Context) (int, error) {
			return 10, nil
		}
		orderRepo.getAllFn = func(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {
			return nil, errors.New("get all query error")
		}

		res, err := svc.GetAllOrders(ctx, "10", "0")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "get all query error")
	})

	t.Run("başarıyla tüm siparişleri ve toplam sayıyı dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountFn = func(ctx context.Context) (int, error) {
			return 5, nil
		}
		expectedList := []st.OrderForAdmin{
			{OrderId: "ORD1", Status: constants.OrderStatePending, TotalAmount: 1000, Name: "Ali", Surname: "Veli"},
		}
		orderRepo.getAllFn = func(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {
			return expectedList, nil
		}

		res, err := svc.GetAllOrders(ctx, "10", "0")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 5, res.Count)
		assert.Equal(t, expectedList, res.Data)
	})
}

func TestOrderService_UpdateOrderState(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("geçersiz sipariş durumu verilirse hata dönmeli", func(t *testing.T) {
		err := svc.UpdateOrderState(ctx, "invalid_state_xyz", "ORD123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Hatalı State Girdiniz")
	})

	t.Run("veritabanı güncellemesi hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.updateOrderStateFn = func(ctx context.Context, newState string, order_id string) error {
			return errors.New("db error")
		}

		err := svc.UpdateOrderState(ctx, constants.OrderStateShipped, "ORD123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Beklenmeyen Hatayla Karşılaşıldı")
	})

	t.Run("geçerli durum ile başarıyla güncellemeli", func(t *testing.T) {
		called := false
		orderRepo.updateOrderStateFn = func(ctx context.Context, newState string, order_id string) error {
			called = true
			assert.Equal(t, constants.OrderStateDelivered, newState)
			assert.Equal(t, "ORD123", order_id)
			return nil
		}

		err := svc.UpdateOrderState(ctx, constants.OrderStateDelivered, "ORD123")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestOrderService_UpdateMultiState(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("hata durumunda hatayı döndürmeli", func(t *testing.T) {
		orderRepo.updateMultiStateFn = func(ctx context.Context, idlist []string, new_state string) error {
			return errors.New("multi update error")
		}

		err := svc.UpdateMultiState(ctx, &st.UpdateMultiOrderReq{
			OrderIds: []string{"ORD1", "ORD2"},
			NewState: "cancelled",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multi update error")
	})

	t.Run("başarıyla birden fazla sipariş durumunu güncellemeli", func(t *testing.T) {
		called := false
		orderRepo.updateMultiStateFn = func(ctx context.Context, idlist []string, new_state string) error {
			called = true
			assert.Equal(t, []string{"ORD1", "ORD2"}, idlist)
			assert.Equal(t, "cancelled", new_state)
			return nil
		}

		err := svc.UpdateMultiState(ctx, &st.UpdateMultiOrderReq{
			OrderIds: []string{"ORD1", "ORD2"},
			NewState: "cancelled",
		})
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestOrderService_GetOrdersByStatus(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("GetOrderCountByStatus hata verirse veya -1 dönerse hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountByStatusFn = func(ctx context.Context, status string) (int, error) {
			return -1, errors.New("status count error")
		}

		res, err := svc.GetOrdersByStatus(ctx, "10", "0", "pending")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetOrdersByStatus sorgusu hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountByStatusFn = func(ctx context.Context, status string) (int, error) {
			return 3, nil
		}
		orderRepo.getOrdersByStatusFn = func(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {
			return nil, errors.New("query error")
		}

		res, err := svc.GetOrdersByStatus(ctx, "10", "0", "pending")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("başarıyla filtrelenmiş siparişleri getirmeli", func(t *testing.T) {
		orderRepo.getOrderCountByStatusFn = func(ctx context.Context, status string) (int, error) {
			return 2, nil
		}
		expectedOrders := []st.OrderForAdmin{
			{OrderId: "ORD1", Status: "pending"},
			{OrderId: "ORD2", Status: "pending"},
		}
		orderRepo.getOrdersByStatusFn = func(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {
			return expectedOrders, nil
		}

		res, err := svc.GetOrdersByStatus(ctx, "10", "0", "pending")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 2, res.Count)
		assert.Equal(t, expectedOrders, res.Data)
	})
}

func TestOrderService_GetOrderCountForWaiting(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("hata durumunda -1 ve hata dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountWaitingFn = func(ctx context.Context) (int, error) {
			return -1, errors.New("waiting count error")
		}

		count, err := svc.GetOrderCountForWaiting(ctx)
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})

	t.Run("başarıyla beklemedeki sipariş sayısını dönmeli", func(t *testing.T) {
		orderRepo.getOrderCountWaitingFn = func(ctx context.Context) (int, error) {
			return 7, nil
		}

		count, err := svc.GetOrderCountForWaiting(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 7, count)
	})
}

func TestOrderService_SearchOrder(t *testing.T) {
	svc, orderRepo, _, _, _, _, _, _, _ := setupOrderServiceHelper()
	ctx := context.Background()

	t.Run("SearchOrderCount hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, searchQuery string) (int, error) {
			return -1, errors.New("search count error")
		}

		res, err := svc.SearchOrder(ctx, "10", "ORD")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "sipariş arama sayısı hatası")
	})

	t.Run("arama sonucunda hiç kayıt yoksa e.ZeroListResult dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, searchQuery string) (int, error) {
			return 0, nil
		}

		res, err := svc.SearchOrder(ctx, "10", "NONEXIST")
		assert.ErrorIs(t, err, e.ZeroListResult)
		assert.Nil(t, res)
	})

	t.Run("SearchOrder sorgusu hata verirse hata dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, searchQuery string) (int, error) {
			return 1, nil
		}
		orderRepo.searchOrderFn = func(ctx context.Context, limit string, searchQuery string) ([]st.OrderForAdmin, error) {
			return nil, errors.New("search query error")
		}

		res, err := svc.SearchOrder(ctx, "10", "ORD")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "sipariş arama hatası")
	})

	t.Run("başarıyla arama sonuçlarını ve sayısını dönmeli", func(t *testing.T) {
		orderRepo.searchOrderCountFn = func(ctx context.Context, searchQuery string) (int, error) {
			return 1, nil
		}
		expectedList := []st.OrderForAdmin{
			{OrderId: "ORD12345", Status: "pending", Name: "Ali", Surname: "Veli"},
		}
		orderRepo.searchOrderFn = func(ctx context.Context, limit string, searchQuery string) ([]st.OrderForAdmin, error) {
			return expectedList, nil
		}

		res, err := svc.SearchOrder(ctx, "10", "ORD")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Count)
		assert.Equal(t, expectedList, res.Data)
	})
}
